package znet

import (
	"context"
	"errors"
	"github.com/lyyym/zinx-wsbase/global"
	"github.com/lyyym/zinx-wsbase/ziface"
	"go.uber.org/zap"
	"io"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Connection 链接
type Connection struct {
	//当前Conn属于哪个Server
	TCPServer ziface.IServer
	//当前连接的socket TCP套接字
	Conn *websocket.Conn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID int64
	//消息管理MsgID和对应处理方法的消息管理模块
	MsgHandler ziface.IMsgHandle
	//告知该链接已经退出/停止的channel
	ctx    context.Context
	cancel context.CancelFunc
	//无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgChan chan *Message
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan *Message

	sync.RWMutex
	//链接属性
	property map[string]interface{}
	////保护当前property的锁
	propertyLock sync.Mutex
	//当前连接的关闭状态
	isClosed bool
	//最后一次心跳时间
	lastHeartBeatTime time.Time
}

// NewConnection 创建连接的方法
func NewConnection(server ziface.IServer, conn *websocket.Conn, connID int64, msgHandler ziface.IMsgHandle) *Connection {
	//初始化Conn属性
	c := &Connection{
		TCPServer:   server,
		Conn:        conn,
		ConnID:      connID,
		isClosed:    false,
		MsgHandler:  msgHandler,
		msgChan:     make(chan *Message),
		msgBuffChan: make(chan *Message, global.Object.MaxMsgChanLen),
		property:    nil,
	}

	//将新创建的Conn添加到链接管理中
	c.TCPServer.GetConnMgr().Add(c)
	return c
}

// StartWriter 写消息Goroutine， 用户将数据发送给客户端
func (c *Connection) StartWriter() {
	//fmt.Println("[Writer Goroutine is running]")
	defer global.Glog.Warn("[conn Writer exit!]")

	for {
		select {
		case msg := <-c.msgChan:
			//有数据要写给客户端
			if err := c.Conn.WriteMessage(msg.GetMsgType(), msg.GetData()); err != nil {
				global.Glog.Error("Send Data error:, ", zap.Error(err))
				return
			}
			c.KeepAlive()
		case msg, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if err := c.Conn.WriteMessage(msg.GetMsgType(), msg.GetData()); err != nil {
					global.Glog.Error("Send Data error:, ", zap.Error(err))
					return
				}
				c.KeepAlive()
			} else {
				global.Glog.Warn("msgBuffChan is Closed")
				break
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// StartReader 读消息Goroutine，用于从客户端中读取数据
func (c *Connection) StartReader() {
	//fmt.Println("[Reader Goroutine is running]")
	defer global.Glog.Warn("[conn Reader exit!]")
	defer c.Stop()

	// 创建拆包解包的对象
	for {
		select {
		case <-c.ctx.Done():
			global.Glog.Warn("conn reader Done")
			return
		default:
			msgType, ioReader, err := c.Conn.NextReader()
			if err != nil {
				global.Glog.Error("get read reader error ", zap.Error(err))
				return
			}
			//读取客户端的Msg head
			headData := make([]byte, c.TCPServer.Packet().GetHeadLen())
			if _, err := io.ReadFull(ioReader, headData); err != nil {
				global.Glog.Error("read msg head error ", zap.Error(err))
				return
			}
			//拆包，得到msgID 和 dataLen 放在msg中
			msg, err := c.TCPServer.Packet().Unpack(headData)
			if err != nil {
				global.Glog.Error("unpack error ", zap.Error(err))
				return
			}
			msg.SetMsgType(msgType)
			//fmt.Println("readMsg - msgType = ", msgType)
			//根据 dataLen 读取 data，放在msg.Data中
			var data []byte
			if msg.GetDataLen() > 0 {
				data = make([]byte, msg.GetDataLen())
				if _, err := io.ReadFull(ioReader, data); err != nil {
					global.Glog.Error("read msg data error ", zap.Error(err))
					return
				}
			}
			msg.SetData(data)
			req := Request{
				conn: c,
				msg:  msg,
			}
			c.KeepAlive()
			if global.Object.WorkerPoolSize > 0 {
				//已经启动工作池机制，将消息交给Worker处理
				c.MsgHandler.SendMsgToTaskQueue(&req)
			} else {
				//从绑定好的消息和对应的处理方法中执行对应的Handle方法
				go c.MsgHandler.DoMsgHandler(&req)
			}
		}
	}
}

// Start 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	//设置消息最大size
	c.Conn.SetReadLimit(int64(global.Object.MaxPacketSize))
	//1 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()
	//2 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()
	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TCPServer.CallOnConnStart(c)
	//开启心跳检测
	go c.heartBeatChecker()

}

// Stop 停止连接，结束当前连接状态M
func (c *Connection) Stop() {

	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TCPServer.CallOnConnStop(c)

	c.Lock()
	defer c.Unlock()

	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}

	global.Glog.Warn("Conn Stop()...", zap.Int64("ConnID =", c.ConnID))

	// 关闭socket链接
	err := c.Conn.Close()
	if err != nil {
		global.Glog.Error("关闭socket链接", zap.Error(err))
	}
	//关闭Writer
	c.cancel()

	//将链接从连接管理器中删除
	c.TCPServer.GetConnMgr().Remove(c)

	//关闭该链接全部管道
	close(c.msgBuffChan)
	//设置标志位
	c.isClosed = true

}

// GetTCPConnection 从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *websocket.Conn {
	return c.Conn
}

// GetConnID 获取当前连接ID
func (c *Connection) GetConnID() int64 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendBinaryMsg 直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendBinaryMsg(msgID uint16, data []byte) (err error) {
	return c.SendMsg(msgID, websocket.BinaryMessage, data)
}

// SendMsg 直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgID uint16, msgType int, data []byte) (err error) {
	c.RLock()
	defer c.RUnlock()
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	//将data封包，并且发送
	dp := c.TCPServer.Packet()
	msg := NewMsgPackage(msgID, msgType, data)
	pack, err := dp.Pack(msg)
	if err != nil {
		global.Glog.Error("Pack error ", zap.Uint16("msg ID = ", msgID))
		return errors.New("Pack error msg ")
	}
	msg.SetData(pack)
	select {
	//写回客户端
	case c.msgChan <- msg:
	default: // 写操作不会阻塞, 因为channel已经预留给websocket一定的缓冲空间
		err = errors.New("ERR_SEND_MESSAGE_FULL")
	}

	return err
}

// SendBinaryBuffMsg  发生BuffMsg
func (c *Connection) SendBinaryBuffMsg(msgID uint16, data []byte) (err error) {
	return c.SendBuffMsg(msgID, websocket.BinaryMessage, data)
}
func (c *Connection) SendBuffMsg(msgID uint16, msgType int, data []byte) (err error) {
	c.RLock()
	defer c.RUnlock()
	if c.isClosed == true {
		return errors.New("connection closed when send buff msg")
	}

	//将data封包，并且发送
	dp := c.TCPServer.Packet()
	msg := NewMsgPackage(msgID, msgType, data)
	pack, err := dp.Pack(msg)
	if err != nil {
		global.Glog.Error("Pack error ", zap.Uint16("msg ID = ", msgID))
		return errors.New("Pack error msg ")
	}
	//fmt.Println("packLenght = ", len(pack), pack)
	msg.SetData(pack)
	select {
	//写回客户端
	case c.msgBuffChan <- msg:
	default: // 写操作不会阻塞, 因为channel已经预留给websocket一定的缓冲空间
		err = errors.New("ERR_SEND_Buff_MESSAGE_FULL")
	}
	return err
}

// SetProperty 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if c.property == nil {
		c.property = make(map[string]interface{})
	}

	c.property[key] = value
}

// GetProperty 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}

	return nil, errors.New("no property found")
}

// RemoveProperty 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

// Context 返回ctx，用于用户自定义的go程获取连接退出状态
func (c *Connection) Context() context.Context {
	return c.ctx
}
