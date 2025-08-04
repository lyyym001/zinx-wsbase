package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lyyym/zinx-wsbase/global"
	"github.com/lyyym/zinx-wsbase/release/sis_server/api"
	"github.com/lyyym/zinx-wsbase/release/sis_server/core"
	"github.com/lyyym/zinx-wsbase/release/sis_server/internal/server/router"
	"github.com/lyyym/zinx-wsbase/ziface"
	"github.com/lyyym/zinx-wsbase/znet"
	"github.com/lyyym/zinx-wsbase/zutil/zuid"
	"go.uber.org/zap"
	"io"
	"log"
	"os"
)

//业务Api 这里定义跟客户都安通信的业务关联
//1	-	登录账号相关
//2 - 	房间业务

// 当客户端建立连接的时候的hook函数
func OnConnecionAdd(conn ziface.IConnection) {
	//创建一个玩家
	player := core.NewPlayer(conn)

	//同步当前玩家的初始化坐标信息给客户端，走MsgID:200消息
	//player.BroadCastStartPosition()

	//将当前新上线玩家添加到worldManager中
	core.WorldMgrObj.AddPlayer(player)

	//将该连接绑定属性PID
	conn.SetProperty("pID", player.PID)

	//同步周边玩家上线信息，与现实周边玩家信息
	//player.SyncSurrounding()

	//同步当前的PlayerID给客户端， 走MsgID:1 消息 这里需要客户端回执 登录信息
	player.SyncPID()

	fmt.Println("=====> Player pIDID = ", player.PID, " arrived ====")
}

// 当客户端断开连接的时候的hook函数
func OnConnectionLost(conn ziface.IConnection) {

	//获取当前连接的PID属性
	pID, _ := conn.GetProperty("pID")
	//fmt.Println("pID = " , pID)
	//根据pID获取对应的玩家对象
	player := core.WorldMgrObj.GetPlayerByPID(pID.(int32))
	if player != nil {
		fmt.Println(player)
		fmt.Println("Player Lost  player= ", player.CID, " Room = ", player.TID)
		//触发玩家下线业务
		if pID != nil {
			fmt.Println("Player Lost  pID= ", pID)
			player.LostConnection()
		}
	}

	//fmt.Println("====> Player ", pID, " left =====")
	//fmt.Println("123")
}

func main() {
	var AppID string = "mdbc"
	//1.初始化配置
	global.InitObject("/conf/app.yaml") ///release/mdbc_server/conf/app.yaml
	//2.初始化日志
	global.InitZap()
	//3.初始化uuid
	zuid.Init()
	//4. 连接sqlite
	//models_sqlite.NewDB()
	pwd, err := os.Getwd()
	if err != nil {
		pwd = ""
	}
	sqlitePath := pwd + global.Object.Sqlite.Dns
	global.SqliteInst = znet.NewSqliteHandle(sqlitePath)
	println("sqlitePath = ", sqlitePath)
	//
	//启动本地老师端
	//print("Start Up TClient\n")
	_, err1 := os.StartProcess("StartTClient.bat", nil, &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
	if err1 != nil {
		print("TClient Started Error\n")
	} else {
		print("TClient Started Succ\n")
	}

	//5.启动InteractionServices
	s := InteractionServices(AppID)

	//6.启动ginServer
	GinServices(s)

	global.Glog.Info("server runing ")

}

func InteractionServices(AppID string) ziface.IServer {

	//从世界启动app
	core.WorldMgrObj.StartApp(AppID)
	//创建服务器句柄
	s := znet.NewServer()
	//注册客户端连接建立和丢失函数
	s.SetOnConnStart(OnConnecionAdd)
	s.SetOnConnStop(OnConnectionLost)
	//注册路由
	//登录路由
	s.AddRouter(1, &api.AccountApi{})
	//聊天路由
	s.AddRouter(2, &api.RoomApi{})
	//课程路由
	s.AddRouter(3, &api.CourseApi{})
	//启动服务
	//global.Glog.Debug("6.InteractionServices run in %s , port = %d \n", utils.GlobalObject.Host, utils.GlobalObject.TCPPort)
	//s.Serve()
	return s
}

// http - services
func GinServices(server ziface.IServer) {

	gin.DisableConsoleColor()
	f, _ := os.Create("./log/gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	e := router.Router(server)
	//e.Use(TlsHandler(8080))
	ginHost := fmt.Sprintf("%s:%d", global.Object.Host, global.Object.GinPort)
	fmt.Println("ginHost = ", ginHost)
	global.Glog.Info(" GinServices run in", zap.String(" iddress = ", ginHost))
	err := e.Run(ginHost)
	//err := e.RunTLS(config.YamlConfig.App.Host, "./cert/server.pem", "./cert/server.key")
	if err != nil {
		log.Fatalln("run err.", err)
		return
	}
	//zlog.Debug("4.GinServices run in ", ginHost)

	fmt.Println("GinServices run in ", ginHost)
}
