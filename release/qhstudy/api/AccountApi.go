package api

import (
	"github.com/lyyym/zinx-wsbase/global"
	"github.com/lyyym/zinx-wsbase/release/qhstudy/core"
	"github.com/lyyym/zinx-wsbase/release/qhstudy/pb"
	"github.com/lyyym/zinx-wsbase/ziface"
	"github.com/lyyym/zinx-wsbase/znet"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"fmt"
	//"log"
)

type AccountApi struct {
	znet.BaseRouter
}

func (aa *AccountApi) Handle(request ziface.IRequest) {

	//1. 得到消息的Sub，用来细化业务实现
	sub := request.GetSubID()
	//global.Glog.Info("recv from client : ", zap.Any("sub", request.GetSubID()),
	//	zap.Any("data", string(request.GetData())))
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

	case 10002: //绑定
		aa.Handle_onRequest10002(player, request.GetData())
		break
	case 10003:
		//同步位置
		aa.Handle_onRequest10003(player, request.GetData())
		break
	case 10004:
		//同步聊天
		//fmt.Println("keep alive")
		aa.Handle_onRequest10004(player, request.GetData())
		break

	}

}

// 绑定
func (aa *AccountApi) Handle_onRequest10002(p *core.Player, data []byte) {

	request_data := &pb.Tcp_Bind{}
	err := proto.Unmarshal(data, request_data)
	if err != nil {
		fmt.Println("proto.Unmarshal err", err)
		return
	}

	//1. 先绑定数据给玩家
	p.CDevice.Dir = request_data.DirVersion
	if p.Bind(request_data.UserToken) {
		fmt.Println("bind ", p.UserName)
		global.Glog.Info("Bind UserName = %s , AccountType = %d , Pid = %d",
			zap.String("UserName", p.UserName),
			zap.Uint32("AccountType", p.AccountType),
			zap.Int32("PID", p.PID),
		)
		//2. 绑定玩家到世界
		core.WorldMgrObj.BindPlayer(p)

	} else {
		global.Glog.Info("Bind Error")
	}
}

func (aa *AccountApi) Handle_onRequest10003(p *core.Player, data []byte) {

	request_data := &pb.Sync_Pos{}
	err := proto.Unmarshal(data, request_data)
	if err != nil {
		fmt.Println("proto.Unmarshal err", err)
		return
	}
	//fmt.Println("位置同步：", request_data.X, request_data.Y, request_data.Z)
	//1. 先绑定数据给玩家
	core.WorldMgrObj.WorldToo(1, 10003, request_data, p.UserName)
}

func (aa *AccountApi) Handle_onRequest10004(p *core.Player, data []byte) {

	request_data := &pb.Sync_Chat{}
	err := proto.Unmarshal(data, request_data)
	if err != nil {
		fmt.Println("proto.Unmarshal err", err)
		return
	}
	fmt.Println("聊天：", request_data.UserName, request_data.Text)
	core.WorldMgrObj.WorldToo(1, 10004, request_data, p.UserName)
}
