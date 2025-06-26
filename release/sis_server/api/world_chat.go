package api

import (
	"encoding/json"
	"fmt"
	"github.com/lyyym/zinx-wsbase/release/sis_server/core"
	"github.com/lyyym/zinx-wsbase/release/sis_server/pb"
	"github.com/lyyym/zinx-wsbase/ziface"
	"github.com/lyyym/zinx-wsbase/znet"
)

// 世界聊天 路由业务
type WorldChatApi struct {
	znet.BaseRouter
}

func (*WorldChatApi) Handle(request ziface.IRequest) {
	//1. 将客户端传来的proto协议解码
	msg := &pb.Talk{}

	fmt.Println("WorldChat Msg Lenght = ", len(request.GetData()), " msg = ", request.GetData())

	err := json.Unmarshal(request.GetData(), msg)
	if err != nil {
		fmt.Println("Talk Unmarshal error ", err)
		return
	}

	//2. 得知当前的消息是从哪个玩家传递来的,从连接属性pID中获取
	pID, err := request.GetConnection().GetProperty("pID")
	if err != nil {
		fmt.Println("GetProperty pID error", err)
		request.GetConnection().Stop()
		return
	}
	//3. 根据pID得到player对象
	player := core.WorldMgrObj.GetPlayerByPID(pID.(int32))

	//4. 让player对象发起聊天广播请求
	player.Talk(msg.Content)
}
