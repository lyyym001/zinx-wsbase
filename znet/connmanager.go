package znet

import (
	"errors"
	"fmt"
	"github.com/lyyym/zinx-wsbase/global"
	"github.com/lyyym/zinx-wsbase/ziface"
	"go.uber.org/zap"
	"sync"
)

// ConnManager 连接管理模块
type ConnManager struct {
	connections map[int64]ziface.IConnection
	connLock    sync.RWMutex
}

// NewConnManager 创建一个链接管理
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[int64]ziface.IConnection),
	}
}

// Add 添加链接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	//将conn连接添加到ConnManager中
	connMgr.connections[conn.GetConnID()] = conn
	connMgr.connLock.Unlock()

	global.Glog.Info("connection add to ConnManager successfully:",
		zap.Int64("connID= ", conn.GetConnID()), zap.Int(" conn num = ", connMgr.Len()))
}

// Remove 删除连接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	//删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	connMgr.connLock.Unlock()
	global.Glog.Info("connection Remove ", zap.Int64("ConnID=", conn.GetConnID()), zap.Int("successfully: conn num = ", connMgr.Len()))
}

// Get 利用ConnID获取链接
func (connMgr *ConnManager) Get(connID int64) (ziface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	}

	return nil, errors.New("connection not found")

}

// Len 获取当前连接
func (connMgr *ConnManager) Len() int {
	connMgr.connLock.RLock()
	length := len(connMgr.connections)
	connMgr.connLock.RUnlock()
	return length
}

// ClearConn 清除并停止所有连接
func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()

	//停止并删除全部的连接信息
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
	}
	connMgr.connLock.Unlock()
	fmt.Println("Clear All Connections successfully: conn num = ", connMgr.Len())
}

// ClearOneConn  利用ConnID获取一个链接 并且删除
func (connMgr *ConnManager) ClearOneConn(connID int64) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	connections := connMgr.connections
	if conn, ok := connections[connID]; ok {
		//停止
		conn.Stop()
		//删除
		delete(connections, connID)
		fmt.Println("Clear Connections ID:  ", connID, "succeed")
		return
	}

	fmt.Println("Clear Connections ID:  ", connID, "err")
	return
}
