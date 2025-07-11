package znet

import "github.com/lyyym/zinx-wsbase/ziface"

// Request 请求
type Request struct {
	conn ziface.IConnection //已经和客户端建立好的 链接
	msg  ziface.IMessage    //客户端请求的数据
}

// GetConnection 获取请求连接信息
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgID 获取请求的消息的ID
func (r *Request) GetMsgID() uint16 {
	return r.msg.GetMsgID()
}

// GetSubID 获取请求的消息的ID
func (r *Request) GetSubID() uint16 {
	return (r.msg.GetMsgID()-10000)%1000 + 10000
}
