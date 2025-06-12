package znet

import (
	"fmt"
	"github.com/lyyym/zinx-wsbase/global"
	"time"
)

// 定时检测心跳包
func (c *Connection) heartBeatChecker() {
	if global.Object.HeartbeatTime == 0 {
		return
	}
	var (
		timer *time.Timer
	)

	timer = time.NewTimer((global.Object.HeartbeatTime) * time.Second)

	for {
		select {
		case <-timer.C:
			if !c.IsAlive() {
				c.Stop()
				//心跳检测失败，结束连接
				if c.isClosed {
					global.Glog.Warn("连接已关闭")
				} else {
					global.Glog.Warn("心跳过期")
				}
				return
			}
			timer.Reset(global.Object.HeartbeatTime * time.Second)
		case <-c.ctx.Done():
			timer.Stop()
			global.Glog.Warn("连接关闭 by ctx.Done")
			return
		}
	}

}

// IsAlive 检测心跳
func (c *Connection) IsAlive() bool {
	var (
		now = time.Now()
	)
	c.Lock()
	defer c.Unlock()
	fmt.Println("now.Sub(c.lastHeartBeatTime) = ", now.Sub(c.lastHeartBeatTime), global.Object.HeartbeatTime*time.Second)
	if c.isClosed || now.Sub(c.lastHeartBeatTime) >
		global.Object.HeartbeatTime*time.Second {
		return false
	}
	return true

}

// KeepAlive 更新心跳
func (c *Connection) KeepAlive() {
	var (
		now = time.Now()
	)
	c.Lock()
	defer c.Unlock()

	c.lastHeartBeatTime = now
}
