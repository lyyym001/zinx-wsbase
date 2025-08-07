package api

import (
	"fmt"
	"github.com/lyyym/zinx-wsbase/global"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/core"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/pb"
	"github.com/lyyym/zinx-wsbase/ziface"
	"github.com/lyyym/zinx-wsbase/znet"
	"google.golang.org/protobuf/proto"
	"net"
)

type RoomApi struct {
	znet.BaseRouter
}

func (ra *RoomApi) Handle(request ziface.IRequest) {

	//1. 得到消息的Sub，用来细化业务实现
	sub := request.GetSubID()

	//fmt.Println("Room Api Do : msgID = " , request.GetMsgID() , " Sub = " , request.GetMsgSub() , " msgLength = " , len(request.GetData()) , " msg = " )

	//2. 得知当前的消息是从哪个玩家传递来的,从连接属性pID中获取
	pID, err := request.GetConnection().GetProperty("pID")
	if err != nil {
		fmt.Println("GetProperty pID error", err)
		request.GetConnection().Stop()
		return
	}
	//3. 根据pID得到player对象
	//fmt.Println("pID ---  = ", pID)
	player := core.WorldMgrObj.GetPlayerByPID(pID.(int32))
	if player == nil {
		return
	}
	//fmt.Println("[Receive Room Msg] : Player = " , player.CID )

	switch sub {

	case 10001: //转发控制解控消息
		ra.Handle_onRequest10001(player, request.GetData())
		break
	case 10002: //控制播放课件
		ra.Handle_onRequest10002(player, request.GetData())
		break
	case 10003: //课件资源加载完成
		ra.Handle_onRequest10003(player, request.GetData())
		break
	case 10004: //对齐进度
		ra.Handle_onRequest10004(player, request.GetData())
		break
	case 10005: //关闭课件
		ra.Handle_onRequest10005(player, request.GetData())
		break
	case 10006: //播放课件
		ra.Handle_onRequest10006(player, request.GetData())
		break
	case 10007: //暂停课件
		ra.Handle_onRequest10007(player, request.GetData())
		break
	case 10008: //快进课件
		ra.Handle_onRequest10008(player, request.GetData())
		break
	case 10009: //快退课件
		ra.Handle_onRequest10009(player, request.GetData())
		break
	case 10010: //播放到
		ra.Handle_onRequest10010(player, request.GetData())
		break
	case 10011: //转发课件加载进度
		ra.Handle_onRequest10011(player, request.GetData())
		break
	case 10012: //转发同步目录
		ra.Handle_onRequest10012(player, request.GetData())
		break
	case 10013:
		ra.Handle_onRequest10013(player, request.GetData())
		break
	}

}

// 转发控制解控消息
func (aa *RoomApi) Handle_onRequest10001(p *core.Player, data []byte) {

	request_data := &pb.SyncPID{}
	err := proto.Unmarshal(data, request_data)
	if err != nil {
		fmt.Println("proto.Unmarshal err", err)
		return
	}

	//赋值
	core.WorldMgrObj.AutoStatus = request_data.PID

	//转发消息
	core.WorldMgrObj.Toa_NoGzNoTeacher(4, 10001, request_data)

}

// 控制播放课件
func (aa *RoomApi) Handle_onRequest10002(p *core.Player, data []byte) {

	request_data := &pb.Tcp_CourseInfo{}
	err := proto.Unmarshal(data, request_data)
	if err != nil {
		fmt.Println("proto.Unmarshal err", err)
		return
	}

	//处理消息
	call := core.WorldMgrObj.SetupMScene(request_data.Cid, request_data.CType, request_data.CMode)
	//转发消息
	core.WorldMgrObj.Toa_NoGzNoTeacher_InScene(4, 10002, request_data)
	//给老师端发送成员列表
	p.SendMsg(4, 10002, call)
}

// 课件资源加载完成
func (aa *RoomApi) Handle_onRequest10003(p *core.Player, data []byte) {

	//处理消息
	readyed, response := core.WorldMgrObj.ReadyMScene(p.UserName)
	//转发消息
	if readyed {
		core.WorldMgrObj.Toa_NoGz_InScene(4, 10003, response)
	}
}

// 对齐进度
func (aa *RoomApi) Handle_onRequest10004(p *core.Player, data []byte) {

	request_data := &pb.Tcp_VideoStep{}
	err := proto.Unmarshal(data, request_data)
	if err != nil {
		fmt.Println("proto.Unmarshal err", err)
		return
	}

	//转发消息
	core.WorldMgrObj.Toa_NoGzNoTeacher_InScene(4, 10004, request_data)

}

// 关闭课件
func (aa *RoomApi) Handle_onRequest10005(p *core.Player, data []byte) {

	request_data := &pb.Tcp_CourseInfo{}
	err := proto.Unmarshal(data, request_data)
	if err != nil {
		fmt.Println("proto.Unmarshal err", err)
		return
	}

	//结束课件
	core.WorldMgrObj.CloseMScene()

	//转发消息
	core.WorldMgrObj.Toa_NoGzNoTeacher_InScene(4, 10005, request_data)

	//发送UDP - 关闭课程
	SendUdpBroadcastToStudent_CloseCourse()

}

func SendUdpBroadcastToStudent_CloseCourse() {

	raddrStu := net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: global.Object.UdpPortDir,
	}
	fmt.Println("SendUdpBroadcast To Student raddrStu = ", global.Object.UdpPortDir, " raddrStu = ", raddrStu)
	connStu, err := net.DialUDP("udp", nil, &raddrStu) //&laddrStu
	if err != nil {
		println(err.Error())
		return
	}

	connStu.Write([]byte("QuitCourse"))
	connStu.Close()
}

// 播放
func (aa *RoomApi) Handle_onRequest10006(p *core.Player, data []byte) {

	//转发消息
	core.WorldMgrObj.Toa_NoGzNoTeacher_InScene(4, 10006, nil)

}

// 暂停
func (aa *RoomApi) Handle_onRequest10007(p *core.Player, data []byte) {

	//转发消息
	core.WorldMgrObj.Toa_NoGzNoTeacher_InScene(4, 10007, nil)

}

// 快进
func (aa *RoomApi) Handle_onRequest10008(p *core.Player, data []byte) {

	//转发消息
	core.WorldMgrObj.Toa_NoGzNoTeacher_InScene(4, 10008, nil)

}

// 快退
func (aa *RoomApi) Handle_onRequest10009(p *core.Player, data []byte) {

	//转发消息
	core.WorldMgrObj.Toa_NoGzNoTeacher_InScene(4, 10009, nil)

}

// 播放到
func (aa *RoomApi) Handle_onRequest10010(p *core.Player, data []byte) {

	request_data := &pb.Tcp_VideoStep{}
	err := proto.Unmarshal(data, request_data)
	if err != nil {
		fmt.Println("proto.Unmarshal err", err)
		return
	}
	//转发消息
	core.WorldMgrObj.Toa_NoGzNoTeacher_InScene(4, 10010, request_data)

}

// 转发课件进度
func (aa *RoomApi) Handle_onRequest10011(p *core.Player, data []byte) {

	request_data := &pb.Tcp_Progress{}
	err := proto.Unmarshal(data, request_data)
	if err != nil {
		fmt.Println("proto.Unmarshal err", err)
		return
	}
	//转发消息给老师端
	fmt.Println("课件加载进度 ", request_data)
	core.WorldMgrObj.ToTeacher(4, 10011, request_data)

}

// 同步目录
func (aa *RoomApi) Handle_onRequest10012(p *core.Player, data []byte) {

	request_data := &pb.Tcp_UInfo{}
	err := proto.Unmarshal(data, request_data)
	if err != nil {
		fmt.Println("proto.Unmarshal err", err)
		return
	}
	//转发消息
	if request_data.UserName == "-99" {
		core.WorldMgrObj.Toa_NoGzNoTeacher(4, 10012, nil)
	} else {
		p := core.WorldMgrObj.GetPlayerByUserName(request_data.UserName)
		if p != nil {
			p.SendMsg(4, 10012, nil)
		}
	}
}

// 控制关机
func (aa *RoomApi) Handle_onRequest10013(p *core.Player, data []byte) {

	request_data := &pb.Tcp_UInfo{}
	err := proto.Unmarshal(data, request_data)
	if err != nil {
		fmt.Println("proto.Unmarshal err", err)
		return
	}
	//转发消息
	if request_data.UserName == "-99" {
		core.WorldMgrObj.ShotDown(4, 10013, nil)
	} else {
		p := core.WorldMgrObj.GetPlayerByUserName(request_data.UserName)
		if p != nil {
			p.SendMsg(4, 10013, nil)
			p.LostConnection()
		}
	}
}
