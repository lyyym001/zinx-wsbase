package core

import (
	"fmt"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/internal/helper"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/internal/models_sqlite"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/pb"
	"github.com/lyyym/zinx-wsbase/ziface"
	"google.golang.org/protobuf/proto"
	"net"
	"sync"
	"time"
)

type DeviceInfo struct {
	Status  int32  //0-故障，1-正常
	Ip      string //ip地址
	Battery int32  //电池电量
	Dir     int64  //目录是否存在
	Free    int32  //是否空闲
}

// 玩家对象
type Player struct {
	//下面是new
	BeatTime int    //心跳检测时间
	UserName string //用户名
	Status   byte   //0-未运行 1-登录服务器中 2-绑定中 3-数据获取中 4-服务器登录成功
	UName    string //真实姓名
	//Ip          string             //ip
	Conn        ziface.IConnection //当前玩家的连接
	AccountType uint32             //账号类型 0-学生 1-老师
	PID         int32              //玩家ID

	CDevice *DeviceInfo //设备信息

	X float32 //平面x坐标
	Y float32 //高度
	Z float32 //平面y坐标 (注意不是Y)
	V float32 //旋转0-360度

}

/*
Player ID 生成器
*/
var PIDGen int32 = 1  //用来生成玩家ID的计数器
var IDLock sync.Mutex //保护PIDGen的互斥机制

// 创建一个玩家对象
func NewPlayer(conn ziface.IConnection) *Player {
	//生成一个PID
	IDLock.Lock()
	ID := PIDGen

	PIDGen++
	IDLock.Unlock()

	tcpAddr := conn.RemoteAddr().(*net.TCPAddr)

	p := &Player{
		PID:     ID,
		Conn:    conn,
		CDevice: &DeviceInfo{Ip: tcpAddr.IP.String(), Status: 0},
		//CLFlag:false,
	}
	return p
}

func (p *Player) Bind(userToken string) {

	userClaim, err := helper.AnalyseToken(userToken)
	if err != nil {
		fmt.Println("tokenAnalyseError")
		return
	}
	p.Status = 2
	p.BeatTime = int(time.Now().Unix())
	p.AccountType = userClaim.AccountType
	p.UserName = userClaim.UserName
	p.UName = userClaim.NickName
	//p.CDevice.Ip = userClaim.Ip
	fmt.Println("user Bind : UserName = ", p.UserName)

	//read from db
	if userClaim.AccountType == 0 {
		data := new(models_sqlite.DeviceBasic)
		err = models_sqlite.DB.Where("username = ?", userClaim.UserName).First(&data).Error
		if err != nil {
			fmt.Println("user read data error : from db ", err.Error())
		} else {
			p.CDevice.Status = data.Status
		}
	}

	p.Status = 3
	//发送userInfo
	p.SendUserInfo()
	//检测心跳
	//go p.BroadcastPlayer()

}

func (p *Player) SendUserInfo() {

	userCall := &pb.Tcp_UserInfo{
		NickName:    p.UName,
		AccountType: p.AccountType,
		UserName:    p.UserName,
	}

	fmt.Println("Send userInfo = ", userCall)

	p.SendMsg(1, 10002, userCall)
	p.Status = 4

}

// 告知客户端被踢了
func (p *Player) Kicked() {

	//p.SendMsg(1, 10003, []byte("ok"))
}

// 告知客户端pID,同步已经生成的玩家ID给客户端
func (p *Player) SyncPID() {

	////发送数据给客户端
	fmt.Println("SendPID To Client", p.PID)
	p.SendMsg(1, 10001, &pb.SyncPID{PID: p.PID})
}

// 广播玩家位置移动
func (p *Player) UpdatePos(x float32, y float32, z float32, v float32) {

	//触发消失视野和添加视野业务
	//计算旧格子gID
	//oldGID := WorldMgrObj.AoiMgr.GetGIDByPos(p.X, p.Z)
	//计算新格子gID
	//newGID := WorldMgrObj.AoiMgr.GetGIDByPos(x, z)

	//更新玩家的位置信息
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	//if oldGID != newGID {
	//	//触发gird切换
	//	//把pID从就的aoi格子中删除
	//	WorldMgrObj.AoiMgr.RemovePIDFromGrID(int(p.PID), oldGID)
	//	//把pID添加到新的aoi格子中去
	//	WorldMgrObj.AoiMgr.AddPIDToGrID(int(p.PID), newGID)
	//
	//	_ = p.OnExchangeAoiGrID(oldGID, newGID)
	//}
	//
	////组装protobuf协议，发送位置给周围玩家
	//msg := &pb.BroadCast{
	//	PID: p.PID,
	//	Tp:  4, //4- 移动之后的坐标信息
	//	Data: &pb.BroadCast_P{
	//		P: &pb.Position{
	//			X: p.X,
	//			Y: p.Y,
	//			Z: p.Z,
	//			V: p.V,
	//		},
	//	},
	//}
	//
	////获取当前玩家周边全部玩家
	//players := p.GetSurroundingPlayers()
	////向周边的每个玩家发送MsgID:200消息，移动位置更新消息
	//for _, player := range players {
	//	player.SendMsg(200, msg)
	//}
}

////广播玩家聊天
//func (p *Player) Talk(content string) {
//
//	//拼接：
//	msg := "["+p.CID+"]：" + content
//
//	//1. 组建MsgID200 proto数据
//	data,_ := json.Marshal(&pb.Talk{Content:msg,})
//
//	//2. 得到当前世界所有的在线玩家
//	players := WorldMgrObj.GetAllPlayers()
//
//	fmt.Println("Talk -> players Length:" , len(players))
//
//	//3. 向所有的玩家发送MsgID:200消息
//	for _, player := range players {
//		player.SendMsg(2,10002 , data)
//	}
//}

// 玩家下线
func (p *Player) LostConnection() {

	//WorldMgrObj.LostConnection(p.PID,p.UserName)
	//pID, err := p.Conn.GetTCPConnection().GetProperty("pID")
	//if err != nil {
	//	fmt.Println("GetProperty pID error", err)
	//	request.GetConnection().Stop()
	//	return
	//}
	//_, err := p.Conn.GetTCPConnection(). .Write([]byte("After ping .....\n"))
	//if err != nil {
	//	fmt.Println("call back ping ping ping error")
	//	p.Conn.Stop()
	//}
	//if p.Conn.GetTCPConnection() {
	//
	//}
	p.Conn.Stop()
	p = nil
}

//广播玩家自己的出生地点
//func (p *Player) BroadCastStartPosition() {
//
//	//组建MsgID200 proto数据
//	msg := &pb.BroadCast{
//		PID: p.PID,
//		Tp:  2, //TP2 代表广播坐标
//		Data: &pb.BroadCast_P{
//			P: &pb.Position{
//				X: p.X,
//				Y: p.Y,
//				Z: p.Z,
//				V: p.V,
//			},
//		},
//	}
//
//	//发送数据给客户端
//	p.SendMsg(200, msg)
//}

//给当前玩家周边的(九宫格内)玩家广播自己的位置，让他们显示自己
//func (p *Player) SyncSurrounding() {
//	//1 根据自己的位置，获取周围九宫格内的玩家pID
//	pIDs := WorldMgrObj.AoiMgr.GetPIDsByPos(p.X, p.Z)
//	//2 根据pID得到所有玩家对象
//	players := make([]*Player, 0, len(pIDs))
//	//3 给这些玩家发送MsgID:200消息，让自己出现在对方视野中
//	for _, pID := range pIDs {
//		players = append(players, WorldMgrObj.GetPlayerByPID(int32(pID)))
//	}
//	//3.1 组建MsgID200 proto数据
//	msg := &pb.BroadCast{
//		PID: p.PID,
//		Tp:  2, //TP2 代表广播坐标
//		Data: &pb.BroadCast_P{
//			P: &pb.Position{
//				X: p.X,
//				Y: p.Y,
//				Z: p.Z,
//				V: p.V,
//			},
//		},
//	}
//	//3.2 每个玩家分别给对应的客户端发送200消息，显示人物
//	for _, player := range players {
//		player.SendMsg(200, msg)
//	}
//	//4 让周围九宫格内的玩家出现在自己的视野中
//	//4.1 制作Message SyncPlayers 数据
//	playersData := make([]*pb.Player, 0, len(players))
//	for _, player := range players {
//		p := &pb.Player{
//			PID: player.PID,
//			P: &pb.Position{
//				X: player.X,
//				Y: player.Y,
//				Z: player.Z,
//				V: player.V,
//			},
//		}
//		playersData = append(playersData, p)
//	}
//
//	//4.2 封装SyncPlayer protobuf数据
//	SyncPlayersMsg := &pb.SyncPlayers{
//		Ps: playersData[:],
//	}
//
//	//4.3 给当前玩家发送需要显示周围的全部玩家数据
//	p.SendMsg(202, SyncPlayersMsg)
//}

////广播玩家位置移动
//func (p *Player) UpdatePos(x float32, y float32, z float32, v float32) {
//
//	//触发消失视野和添加视野业务
//	//计算旧格子gID
//	oldGID := WorldMgrObj.AoiMgr.GetGIDByPos(p.X, p.Z)
//	//计算新格子gID
//	newGID := WorldMgrObj.AoiMgr.GetGIDByPos(x, z)
//
//	//更新玩家的位置信息
//	p.X = x
//	p.Y = y
//	p.Z = z
//	p.V = v
//
//	if oldGID != newGID {
//		//触发gird切换
//		//把pID从就的aoi格子中删除
//		WorldMgrObj.AoiMgr.RemovePIDFromGrID(int(p.PID), oldGID)
//		//把pID添加到新的aoi格子中去
//		WorldMgrObj.AoiMgr.AddPIDToGrID(int(p.PID), newGID)
//
//		_ = p.OnExchangeAoiGrID(oldGID, newGID)
//	}
//
//	//组装protobuf协议，发送位置给周围玩家
//	msg := &pb.BroadCast{
//		PID: p.PID,
//		Tp:  4, //4- 移动之后的坐标信息
//		Data: &pb.BroadCast_P{
//			P: &pb.Position{
//				X: p.X,
//				Y: p.Y,
//				Z: p.Z,
//				V: p.V,
//			},
//		},
//	}
//
//	//获取当前玩家周边全部玩家
//	players := p.GetSurroundingPlayers()
//	//向周边的每个玩家发送MsgID:200消息，移动位置更新消息
//	for _, player := range players {
//		player.SendMsg(200, msg)
//	}
//}

//func (p *Player) OnExchangeAoiGrID(oldGID, newGID int) error {
//	//获取就的九宫格成员
//	oldGrIDs := WorldMgrObj.AoiMgr.GetSurroundGrIDsByGID(oldGID)
//
//	//为旧的九宫格成员建立哈希表,用来快速查找
//	oldGrIDsMap := make(map[int]bool, len(oldGrIDs))
//	for _, grID := range oldGrIDs {
//		oldGrIDsMap[grID.GID] = true
//	}
//
//	//获取新的九宫格成员
//	newGrIDs := WorldMgrObj.AoiMgr.GetSurroundGrIDsByGID(newGID)
//	//为新的九宫格成员建立哈希表,用来快速查找
//	newGrIDsMap := make(map[int]bool, len(newGrIDs))
//	for _, grID := range newGrIDs {
//		newGrIDsMap[grID.GID] = true
//	}
//
//	//------ > 处理视野消失 <-------
//	offlineMsg := &pb.SyncPID{
//		PID: p.PID,
//	}
//
//	//找到在旧的九宫格中出现,但是在新的九宫格中没有出现的格子
//	leavingGrIDs := make([]*GrID, 0)
//	for _, grID := range oldGrIDs {
//		if _, ok := newGrIDsMap[grID.GID]; !ok {
//			leavingGrIDs = append(leavingGrIDs, grID)
//		}
//	}
//
//	//获取需要消失的格子中的全部玩家
//	for _, grID := range leavingGrIDs {
//		players := WorldMgrObj.GetPlayersByGID(grID.GID)
//		for _, player := range players {
//			//让自己在其他玩家的客户端中消失
//			player.SendMsg(201, offlineMsg)
//
//			//将其他玩家信息 在自己的客户端中消失
//			anotherOfflineMsg := &pb.SyncPID{
//				PID: player.PID,
//			}
//			p.SendMsg(201, anotherOfflineMsg)
//			time.Sleep(200 * time.Millisecond)
//		}
//	}
//
//	//------ > 处理视野出现 <-------
//
//	//找到在新的九宫格内出现,但是没有在就的九宫格内出现的格子
//	enteringGrIDs := make([]*GrID, 0)
//	for _, grID := range newGrIDs {
//		if _, ok := oldGrIDsMap[grID.GID]; !ok {
//			enteringGrIDs = append(enteringGrIDs, grID)
//		}
//	}
//
//	onlineMsg := &pb.BroadCast{
//		PID: p.PID,
//		Tp:  2,
//		Data: &pb.BroadCast_P{
//			P: &pb.Position{
//				X: p.X,
//				Y: p.Y,
//				Z: p.Z,
//				V: p.V,
//			},
//		},
//	}
//
//	//获取需要显示格子的全部玩家
//	for _, grID := range enteringGrIDs {
//		players := WorldMgrObj.GetPlayersByGID(grID.GID)
//
//		for _, player := range players {
//			//让自己出现在其他人视野中
//			player.SendMsg(200, onlineMsg)
//
//			//让其他人出现在自己的视野中
//			anotherOnlineMsg := &pb.BroadCast{
//				PID: player.PID,
//				Tp:  2,
//				Data: &pb.BroadCast_P{
//					P: &pb.Position{
//						X: player.X,
//						Y: player.Y,
//						Z: player.Z,
//						V: player.V,
//					},
//				},
//			}
//
//			time.Sleep(200 * time.Millisecond)
//			p.SendMsg(200, anotherOnlineMsg)
//		}
//	}
//
//	return nil
//}

//获得当前玩家的AOI周边玩家信息
//func (p *Player) GetSurroundingPlayers() []*Player {
//	//得到当前AOI区域的所有pID
//	pIDs := WorldMgrObj.AoiMgr.GetPIDsByPos(p.X, p.Z)
//
//	//将所有pID对应的Player放到Player切片中
//	players := make([]*Player, 0, len(pIDs))
//	for _, pID := range pIDs {
//		players = append(players, WorldMgrObj.GetPlayerByPID(int32(pID)))
//	}
//
//	return players
//}

/*
发送消息给客户端，
主要是将pb的protobuf数据序列化之后发送
*/
func (p *Player) SendMsg(msgID uint16, subID uint16, data proto.Message) {
	//fmt.Println("SendMsgToClient,user = ", p.UserName, " MsgId = ", msgID, " msgSub = ", msgSub)
	//将NetBody结构体序列化
	//创建一个存放bytes字节的缓冲
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err: ", err)
		return
	}
	//fmt.Printf("after Marshal data = %+v\n", msg)
	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}
	//调用Zinx框架的SendMsg发包
	//if err := p.Conn.SendMsg(msgID, msgType, msg); err != nil {
	//	fmt.Println("Player SendMsg error !")
	//	return
	//}
	mId := msgID*1000 + subID
	if err := p.Conn.SendBinaryBuffMsg(mId, msg); err != nil {
		fmt.Println("Player SendMsg error !")
		return
	}
	return
}

//func (p *Player) BroadcastPlayer()  {
//
//
//
//
//	t1 := time.NewTimer(time.Millisecond * 5000) //5s
//L:
//	for {
//		if p == nil || p.TID == p.CID {
//			break L
//		}
//		//fmt.Println("len(cRoom.ClientHandle), ", len(cRoom.ClientHandle), "cRoom.TeacherCli,", cRoom.TeacherCli)
//		/*if len(cRoom.ClientHandle) == 0 && cRoom.TeacherCli == nil{
//			goto ForEnd
//		}*/
//		select {
//		case <-t1.C:
//			t1.Reset(time.Millisecond * 5000)
//			//SendUdpBroadcastToStudent(rm.TID)
//			if p == nil{
//				break L
//			}
//			_BeatTime := int(time.Now().Unix())
//			if _BeatTime - p.BeatTime > 22 {
//				//fmt.Println("学生端掉线了")
//				fmt.Println("学生掉线了 ： " ,p.PID,p.CID , p.TID)
//				//p.Conn.Stop()
//				p.LostConnection()
//				//p.Conn.Stop()
//
//				break L
//			}
//		}
//	}
//
//}
