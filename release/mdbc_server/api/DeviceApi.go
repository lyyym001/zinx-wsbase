package api

import (
	"fmt"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/core"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/internal/models_sqlite"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/pb"
	"github.com/lyyym/zinx-wsbase/ziface"
	"github.com/lyyym/zinx-wsbase/znet"
	"google.golang.org/protobuf/proto"
)

type DeviceApi struct {
	znet.BaseRouter
}

func (aa *DeviceApi) Handle(request ziface.IRequest) {

	//1. 得到消息的Sub，用来细化业务实现
	sub := request.GetMsgID()
	//fmt.Println("Account Api Do : msgID = " , request.GetMsgID() , " Sub = " , request.GetMsgSub() , " msgLength = " , len(request.GetData()))

	//2. 得知当前的消息是从哪个玩家传递来的,从连接属性pID中获取
	pID, err := request.GetConnection().GetProperty("pID")
	if err != nil {
		fmt.Println("GetProperty pID error", err)
		request.GetConnection().Stop()
		return
	}
	//3. 根据pID得到player对象
	player := core.WorldMgrObj.GetPlayerByPID(pID.(int32))
	if player == nil {
		return
	}
	//fmt.Println("[Receive Account Msg] : Player = " , player.PID )

	switch sub {

	case 10001: //设置故障
		aa.Handle_onRequest10001(player, request.GetData())
		break
	case 10002: //同步目录版本
		aa.Handle_onRequest10002(player, request.GetData())
		break
	case 10003: //推送学生端数据(电池电量...)
		aa.Handle_onRequest10003(player, request.GetData())
		break
	case 10004: //请求监控
		aa.Handle_onRequest10004(player, request.GetData())
		break
	case 10005: //监控流启动成功
		aa.Handle_onRequest10005(player, request.GetData())
		break
	case 10006: //关闭监控
		aa.Handle_onRequest10006(player, request.GetData())
		break
	case 10007: //监控质量调整
		aa.Handle_onRequest10007(player, request.GetData())
		break
	case 10008: //开启推流
		aa.Handle_onRequest10008(player, request.GetData())
		break
	}

}

// 设置故障
func (aa *DeviceApi) Handle_onRequest10001(p *core.Player, data []byte) {

	request_data := &pb.Tcp_Gz{}
	err := proto.Unmarshal(data, request_data)
	if err != nil {
		fmt.Println("proto.Unmarshal err", err)
		return
	}
	call := &pb.Tcp_Gz_Response{Info: &pb.Tcp_Info{Code: -1}}
	//1. 设置数据
	p.CDevice.Status = request_data.Status
	// update db
	err = models_sqlite.DB.Model(&models_sqlite.DeviceBasic{}).Where("username = ?", request_data.UserName).Update("status", request_data.Status).Error
	if err != nil {
		call.Info.Msg = "设置故障db异常"
	} else {
		call.Info.Code = 200
		call.Info.Msg = "操作成功"
		call.GzStatus = request_data
	}

	p.SendMsg(3, 10001, call)
}

// 推送目录版本
func (aa *DeviceApi) Handle_onRequest10002(p *core.Player, data []byte) {

	request_data := &pb.Tcp_Version{}
	err := proto.Unmarshal(data, request_data)
	if err != nil {
		fmt.Println("proto.Unmarshal err", err)
		return
	}

	//1. 设置数据
	p.CDevice.Dir = request_data.DirVersion
}

// 推送学生端数据(电池电量...)
func (aa *DeviceApi) Handle_onRequest10003(p *core.Player, data []byte) {

	request_data := &pb.Tcp_StudentStatus{}
	err := proto.Unmarshal(data, request_data)
	if err != nil {
		fmt.Println("proto.Unmarshal err", err)
		return
	}

	//1. 设置数据
	p.CDevice.Battery = request_data.Battery
	p.CDevice.Free = request_data.Free

	response_data := &pb.Response_StudentStatus{}
	response_data.Free = request_data.Free
	response_data.Number = p.UserName
	//fmt.Println(response_data)
	core.WorldMgrObj.ToTeacher(3, 10003, response_data)
}

// 请求监控
func (aa *DeviceApi) Handle_onRequest10004(p *core.Player, data []byte) {

	request_data := &pb.Tcp_RequestJk{}
	err := proto.Unmarshal(data, request_data)
	if err != nil {
		fmt.Println("proto.Unmarshal err", err)
		return
	}
	response_data := &pb.Tcp_ResponseJk{}
	//新用户
	player_new := core.WorldMgrObj.GetPlayerByUserName(request_data.UserName)
	if player_new == nil {
		response_data.Code = 0
		p.SendMsg(3, 10004, response_data)
		return
	}

	//// 1.停掉之前的流
	//if len(request_data.OldUserName) > 0 {
	//	player := core.WorldMgrObj.GetPlayerByUserName(request_data.OldUserName)
	//	if player != nil {
	//		response_data.Code = 1
	//		player.SendMsg(3, 10004, response_data)
	//	}
	//}

	//开启新的流
	if player_new != nil {
		response_data.Code = 2
		response_data.Progress = request_data.Progress
		player_new.SendMsg(3, 10004, response_data)
		p.SendMsg(3, 10004, response_data)
	}
}

// 监控流播放成功
func (aa *DeviceApi) Handle_onRequest10005(p *core.Player, data []byte) {

	core.WorldMgrObj.ToTeacher(3, 10005, nil)
}

// 关闭监控
func (aa *DeviceApi) Handle_onRequest10006(p *core.Player, data []byte) {

	request_data := &pb.Tcp_UInfo{}
	err := proto.Unmarshal(data, request_data)
	if err != nil {
		fmt.Println("proto.Unmarshal err", err)
		return
	}

	if len(request_data.UserName) == 0 {
		return
	}

	//response_data := &pb.Tcp_Info{}
	//关闭监控流
	player := core.WorldMgrObj.GetPlayerByUserName(request_data.UserName)
	if player != nil {
		//response_data.Code = 2
		player.SendMsg(3, 10006, nil)
		//p.SendMsg(3, 10004, response_data)
	}
	//else {
	//	response_data.Code = 0
	//	p.SendMsg(3, 10004, response_data)
	//}
}

// 监控质量调整
func (aa *DeviceApi) Handle_onRequest10007(p *core.Player, data []byte) {
	request_data := &pb.Tcp_RequestJk{}
	err := proto.Unmarshal(data, request_data)
	if err != nil {
		fmt.Println("proto.Unmarshal err", err)
		return
	}
	response_data := &pb.Tcp_ResponseJk{}
	player := core.WorldMgrObj.GetPlayerByUserName(request_data.UserName)
	if player != nil {
		response_data.Code = 1
		response_data.Progress = request_data.Progress
		player.SendMsg(3, 10007, response_data)
	}
}

// 开启结束推流
func (aa *DeviceApi) Handle_onRequest10008(p *core.Player, data []byte) {
	request_data := &pb.SyncPID{}
	err := proto.Unmarshal(data, request_data)
	if err != nil {
		fmt.Println("proto.Unmarshal err", err)
		return
	}

	//转发
	core.WorldMgrObj.Toa_NoGzNoTeacher(3, 10008, request_data)
}
