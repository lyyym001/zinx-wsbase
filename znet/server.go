package znet

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lyyym/zinx-wsbase/global"
	"github.com/lyyym/zinx-wsbase/ziface"
	"github.com/lyyym/zinx-wsbase/zutil/zuid"
	"go.uber.org/zap"
	"net/http"
)

// Server 接口实现，定义一个Server服务类
type Server struct {
	//服务器的名称
	Name string
	//tcp4 or other
	IPVersion string
	//服务绑定的IP地址
	IP string
	//服务绑定的端口
	Port int
	//当前Server的消息管理模块，用来绑定MsgID和对应的处理方法
	msgHandler ziface.IMsgHandle
	//当前Server的链接管理器
	ConnMgr ziface.IConnManager
	//该Server的连接创建时Hook函数
	OnConnStart func(conn ziface.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn ziface.IConnection)

	packet ziface.Packet
}

// NewServer 创建一个服务器句柄
func NewServer(opts ...Option) ziface.IServer {
	//global.InitObject()
	//global.InitZap()
	//zuid.Init()

	printLogo()
	s := &Server{
		Name:       global.Object.Name,
		IPVersion:  "tcp4",
		IP:         global.Object.Host,
		Port:       global.Object.TCPPort,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
		packet:     NewDataPack(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

//============== 实现 ziface.IServer 里的全部接口方法 ========

// Start 开启网络服务
func (s *Server) Start(c *gin.Context) {

	//开启一个go去做服务端Linster业务
	go func() {
		var (
			err        error
			wsSocket   *websocket.Conn
			wsUpgrader = websocket.Upgrader{
				// 允许所有CORS跨域请求
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			}
		)
		if wsSocket, err = wsUpgrader.Upgrade(c.Writer, c.Request, nil); err != nil {
			global.Glog.Error("将HTTP服务器连接升级到WebSocket协议失败 ", zap.Error(err))
			return
		}

		//3 启动server网络连接业务
		//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
		//todo 是否可以关闭此新的连接?
		//if s.ConnMgr.Len() >= utils.Object.MaxConn {
		//	wsSocket.Close()
		//	return
		//}
		//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
		dealConn := NewConnection(s, wsSocket, zuid.Gen64(), s.msgHandler)
		//3.4 启动当前链接的处理业务
		go dealConn.Start()
	}()
}

// Stop 停止服务
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)

	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnMgr.ClearConn()
}

// Serve 运行服务
func (s *Server) Serve(c *gin.Context) {
	s.Start(c)

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出
	select {}
}

// AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(msgID uint16, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
}

// GetConnMgr 得到链接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// SetOnConnStart 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		//fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		//fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}

func (s *Server) Packet() ziface.Packet {
	return s.packet
}

var topLine = `┌───────────────────────────────────────────────────┐`
var borderLine = `│`
var bottomLine = `└───────────────────────────────────────────────────┘`

func printLogo() {
	//fmt.Println(zinxLogo)
	fmt.Println(topLine)
	fmt.Println(fmt.Sprintf("%s [Github] https://github.com/sun-fight/zinx-websocket                 %s", borderLine, borderLine))
	fmt.Println(fmt.Sprintf("%s [tutorial] https://github.com/sun-fight/zinx-websocket/blob/master/README.md     %s", borderLine, borderLine))
	fmt.Println(bottomLine)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		global.Object.Version,
		global.Object.MaxConn,
		global.Object.MaxPacketSize)
}
