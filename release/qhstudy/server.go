package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lyyym/zinx-wsbase/global"
	"github.com/lyyym/zinx-wsbase/release/qhstudy/api"
	"github.com/lyyym/zinx-wsbase/release/qhstudy/core"
	"github.com/lyyym/zinx-wsbase/release/qhstudy/internal/server/router"
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

	//==============同步周边玩家上线信息，与现实周边玩家信息========
	//player.SyncSurrounding()

	fmt.Println("=====> Player pIDID = ", player.PID, " arrived ====")

	//fmt.Println(gameutils.GlobalScene)
}

// 当客户端断开连接的时候的hook函数
func OnConnectionLost(conn ziface.IConnection) {

	//获取当前连接的PID属性
	pID, _ := conn.GetProperty("pID")
	//fmt.Println("pID = " , pID)
	fmt.Println("有客户端断开了连接 , pid =  ", pID)
	//根据pID获取对应的玩家对象
	player := core.WorldMgrObj.GetPlayerByPID(pID.(int32))
	if player != nil {
		//fmt.Println(player)
		fmt.Println("[断开连接用户]pid= ", player.PID, ",UserName =", player.UserName)
		//触发玩家下线业务
		core.WorldMgrObj.LostConnection(player.PID, player.UserName, 0)
		player.LostConnection()
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
	s.AddRouter(1, &api.AccountApi{}) //登录路由
	s.AddRouter(2, &api.WorkApi{})    //作品
	s.AddRouter(3, &api.DeviceApi{})  //设备
	s.AddRouter(4, &api.RoomApi{})    //房间
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

func StunServer(ip string, port int) {

	//fmt.Println("Stun Server Ready To Start")
	//
	//publicIP := ip //flag.String("public-ip", "192.168.0.22", "IP Address that STUN can be contacted by.")
	////port           //flag.Int("port", 3478, "Listening port.")
	//flag.Parse()
	//
	//if len(publicIP) == 0 {
	//	fmt.Println("'public-ip' is required")
	//	log.Fatalf("'public-ip' is required")
	//}
	//
	//// Create a UDP listener to pass into pion/turn
	//// pion/turn itself doesn't allocate any UDP sockets, but lets the user pass them in
	//// this allows us to add logging, storage or modify inbound/outbound traffic
	//udpListener, err := net.ListenPacket("udp4", "0.0.0.0:"+strconv.Itoa(port))
	//if err != nil {
	//	fmt.Println("Failed to create STUN server listener:", err)
	//	log.Panicf("Failed to create STUN server listener: %s", err)
	//}
	//fmt.Println("Stun Server Ready To NewServer")
	//s, err := turn.NewServer(turn.ServerConfig{
	//	// PacketConnConfigs is a list of UDP Listeners and the configuration around them
	//	PacketConnConfigs: []turn.PacketConnConfig{
	//		{
	//			PacketConn: udpListener,
	//		},
	//	},
	//})
	//if err != nil {
	//	log.Panic(err)
	//}
	//fmt.Println("Stun Server Ready To ServerUped")
	//fmt.Println("Listener Ip = " + publicIP)
	//fmt.Println("Listener Port " + strconv.Itoa(port))
	//// Block until user sends SIGINT or SIGTERM
	//sigs := make(chan os.Signal, 1)
	//signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//<-sigs
	//
	//if err = s.Close(); err != nil {
	//	log.Panic(err)
	//}

}
