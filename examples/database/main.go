package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lyyym/zinx-wsbase/global"
	"github.com/lyyym/zinx-wsbase/ziface"
	"github.com/lyyym/zinx-wsbase/znet"
)

func main() {
	server := znet.NewServer()

	global.InitGormMysql()
	global.InitRedis()
	fmt.Println(global.Mysql)
	fmt.Println(global.Redis)
	server.SetOnConnStart(OnConnecionAdd)
	bindAddress := fmt.Sprintf("%s:%d", global.Object.Host, global.Object.TCPPort)
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.GET("/ws", server.Serve)
	router.Run(bindAddress)
}

func OnConnecionAdd(conn ziface.IConnection) {
	fmt.Println("OnConnecionAdd : ", conn.RemoteAddr().String())
}
