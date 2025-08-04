package core

import (
	"encoding/json"
	"fmt"
	"github.com/lyyym/zinx-wsbase/global"
	"github.com/lyyym/zinx-wsbase/release/qhstudy/pb"
	"google.golang.org/protobuf/proto"
	"net"
	"sync"
	"time"
)

type SceneInfo struct {
	Mode           byte //1-学练模式 2-考评模式 3-协同模式
	Running        bool //是否正在运行
	teacherReadyed bool //老师准备状态
	Players        map[string]int32
	TakeObjects    map[int]int32
	JHObjects      map[int]*JHObject
	Steps          map[int]*Step
	Questions      map[int]int32
	CI             *CourseInfo
	//GlobalStepId int
}

type CourseInfo struct {
	Cid          string //课程id
	CType        int32  //课件类型
	Mode         int32  //课程模式
	MainCtroller string //主控玩家
}

/*
当前游戏世界的总管理模块
*/
type WorldManager struct {
	AoiMgr     *AOIManager       //当前世界地图的AOI规划管理器
	Players    map[int32]*Player //当前在线的玩家集合
	PMs        map[string]int32  //记录登录的玩家，防止重复登录,map[UserName]Pid
	pLock      sync.RWMutex      //保护Players的互斥读写机制
	AppID      string            //用户登录标识，无此标识则不允许登录
	TUserName  string            //老师账号
	TUid       int32             //老师UID
	AutoStatus int32             //控制状态 0-未控制 1-控制
	DirVersion int64             //当前目录版本
	MScene     *SceneInfo        //当前场景状态
}

// 提供一个对外的世界管理模块句柄
var WorldMgrObj *WorldManager

// 提供WorldManager 初始化方法
func init() {
	WorldMgrObj = &WorldManager{
		Players: make(map[int32]*Player),
		AoiMgr:  NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
		PMs:     make(map[string]int32),
		MScene:  &SceneInfo{Running: false},
	}

	//// 2.读取系统状态
	//sData := &models_sqlite.SysBasic{}
	//err := models_sqlite.DB.Where("sid = ?", 1).First(sData).Error
	//if err == nil {
	//	//设置状态
	//	WorldMgrObj.DirVersion = sData.DirVersion
	//}

}

// 提供添加一个玩家的的功能，将玩家添加进玩家信息表Players
func (wm *WorldManager) StartApp(appid string) {
	//将player添加到 世界管理器中
	wm.pLock.Lock()
	wm.AppID = appid
	wm.pLock.Unlock()

	//将player 添加到AOI网络规划中
	//wm.AoiMgr.AddToGrIDByPos(int(player.PID), player.X, player.Z)

	//启动udp用来广播自己的ip及端口
	go wm.BroadcastNet()

}

func (wm *WorldManager) BroadcastNet() {

	t1 := time.NewTimer(time.Millisecond * 1000) //1s
	//L:
	for {

		select {
		case <-t1.C:
			t1.Reset(time.Millisecond * 1000)
			SendUdpBroadcastToAll()
		}
	}
}

func SendUdpBroadcastToAll() {
	var sData pb.Sync_Hello
	sData.Ip = global.Object.Host
	sData.Port = global.Object.TCPPort
	sData.GinPort = global.Object.GinPort
	//zlog.Debugf("Broadcast Ip=%s,TcpPort=%d,GinPort = %d", sData.Ip, sData.Port, sData.GinPort)
	//global.Glog.Info("udp broadcast ", zap.String("Ip", sData.Ip), zap.Int("TcpPort", sData.Port), zap.Int("GinPort", sData.GinPort))
	// 这里设置接收者的IP地址为广播地址
	raddrStu := net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: global.Object.UdpPort,
	}
	//fmt.Println("SendUdpBroadcast To Student raddrStu = ",laddrStu," raddrStu = ",raddrStu)
	connStu, err := net.DialUDP("udp", nil, &raddrStu) //&laddrStu
	if err != nil {
		println(err.Error())
		return
	}

	//fmt.Println("pBroadcast LanServer ", raddrStu, "data = ", sData)
	x, _ := json.Marshal(sData)
	connStu.Write(x)
	connStu.Close()
}

// 提供添加一个玩家的的功能，将玩家添加进玩家信息表Players
func (wm *WorldManager) BindPlayer(player *Player) {
	//将player添加到 世界管理器中
	wm.pLock.Lock()

	if _, ok := wm.Players[player.PID]; ok {

		//3. 记录真实玩家(已经发生登录的玩家)
		if _, has := wm.PMs[player.UserName]; has {
			fmt.Println("已经有相同账号记录在服务器")
			//1. 有可能被顶号
			//2. 有可能重复绑定
			// 暂时不处理
		}
		wm.PMs[player.UserName] = player.PID

		if player.AccountType == 1 {
			//记录老师状态
			wm.TUserName = player.UserName
			wm.TUid = player.PID
			//fmt.Println("[记录老师状态]")
		}

		//检测作品成员列表
		//if wm.MScene.Running {
		//	if _, ok := wm.MScene.Players[username]; ok {
		//		//delete(wm.MScene.Players,pID)
		//		wm.MScene.Players[username] = 1
		//	}
		//}

	}

	wm.pLock.Unlock()

	fmt.Println("Players = ", wm.PMs)

	//wm.PrintUserList()

	//将player 添加到AOI网络规划中
	//wm.AoiMgr.AddToGrIDByPos(int(player.PID), player.X, player.Z)
}

// 提供添加一个玩家的的功能，将玩家添加进玩家信息表Players
func (wm *WorldManager) AddPlayer(player *Player) {
	//将player添加到 世界管理器中
	wm.pLock.Lock()
	wm.Players[player.PID] = player
	//wm.PMs[player.CID] = true
	wm.pLock.Unlock()

	//将player 添加到AOI网络规划中
	//wm.AoiMgr.AddToGrIDByPos(int(player.PID), player.X, player.Z)
}

// 开启多人课件
func (wm *WorldManager) SetupMScene(cid string, cType int32, mode int32) *pb.Tcp_Members {

	// 1. 初始化数据
	if wm.MScene != nil {
		wm.MScene.Players = make(map[string]int32)
		wm.MScene.JHObjects = make(map[int]*JHObject)
		wm.MScene.Questions = make(map[int]int32)
		wm.MScene.Steps = make(map[int]*Step)
		wm.MScene.TakeObjects = make(map[int]int32)
		wm.MScene.teacherReadyed = false
		wm.MScene.Running = true
		wm.MScene.CI = &CourseInfo{Cid: cid, CType: cType, Mode: mode, MainCtroller: ""}
	}

	if mode == 3 {
		wm.MScene.CI.MainCtroller = "teacher"
	}

	// 2. 封装用户(在线无故障的)
	call := &pb.Tcp_Members{}
	for username, pid := range wm.PMs {
		if p, ok := wm.Players[pid]; ok {
			if p.CDevice.Status == 1 {
				wm.MScene.Players[username] = 0 //0-未准备
				call.Ms = append(call.Ms, username)
			}
		}
	}
	return call
}

// 课件资源准备完成
func (wm *WorldManager) ReadyMScene(userName string) (bool, *pb.Tcp_CoursePlay) {

	// 1. 初始化数据
	if wm.MScene != nil && wm.MScene.Running {

		if len(wm.MScene.Players) == 0 && userName == "teacher" {

			return true, &pb.Tcp_CoursePlay{
				Cid:        wm.MScene.CI.Cid,
				CMode:      wm.MScene.CI.Mode,
				Contorller: wm.MScene.CI.MainCtroller,
				CType:      wm.MScene.CI.CType,
				PartNumber: 1,
			}
		}

		if userName == "teacher" {
			wm.MScene.teacherReadyed = true
		} else if _, ok := wm.MScene.Players[userName]; ok {
			wm.MScene.Players[userName] = 1
		}

		for _, state := range wm.MScene.Players {
			if state == 0 {
				return false, nil
			}
		}

		if !wm.MScene.teacherReadyed {
			return false, nil
		}

		return true, &pb.Tcp_CoursePlay{
			Cid:        wm.MScene.CI.Cid,
			CMode:      wm.MScene.CI.Mode,
			Contorller: wm.MScene.CI.MainCtroller,
			CType:      wm.MScene.CI.CType,
			PartNumber: int32(len(wm.MScene.Players)) + 1,
		}

	}

	return false, nil
}

// 结束课件
func (wm *WorldManager) CloseMScene() {

	// 1. 初始化数据
	if wm.MScene != nil {
		wm.MScene.Running = false
	}

}

func (wm *WorldManager) PrintUserList() {
	fmt.Println("[世界成员列表]")
	for _key, value := range wm.PMs {
		fmt.Println("->username=", _key, ",uid=", value)
	}

	if wm.MScene != nil && wm.MScene.Players != nil && len(wm.MScene.Players) > 0 {
		fmt.Println("[当前作品成员列表]")
		for _key, value := range wm.MScene.Players {
			fmt.Println("->username=", _key, ",状态(0-资源未加载，1-资源已加载)=", value)
		}
	} else {
		fmt.Println("[当前作品成员列表](无)")
	}
}

// 玩家掉线了
// code 0-掉线 1-被踢
func (wm *WorldManager) LostConnection(pID int32, userName string, code int) {

	wm.RemovePlayer(pID)
	wm.RemovePlayerUName(userName)
	wm.RemoveWorkPlayer(userName)

	wm.PrintUserList()

	//if code == 0 {
	//	if wm.MScene.Running {
	//		if wm.MScene.MainCtroller == userName {
	//			fmt.Println("[离线重新指定主控]原主控=", userName)
	//			wm.ChangeMainCtroller(1)
	//		} else {
	//			wm.CheckOver(1)
	//		}
	//	}
	//}
}

// code 0-用户中途离开指定主控 1-用户掉线指定主控
func (wm *WorldManager) ChangeMainCtroller(code int) bool {

	//重新指定一个主控
	for username, state := range wm.MScene.Players {

		if state == 1 {
			if uid, ok := wm.PMs[username]; ok {
				p := wm.GetPlayerByPID(uid)
				if p != nil {
					fmt.Println("[离线重新指定主控]新主控=", username)
					p.SendMsg(2, 10008, &pb.Tcp_Info{Code: 200})
				}
				return true
			}
		}
	}

	//没有人了，结束作品
	wm.WorkFinish(0)
	return false

}

// code 0-用户离开了检测结束 1-用户离线了检测结束
func (wm *WorldManager) CheckOver(code int) bool {

	//检测是否还有人
	for _, state := range wm.MScene.Players {

		if state == 1 {
			return false
		}
	}

	//没有找到有效的主控则结束课程
	wm.WorkFinish(0)

	return true
}

// 世界同步给所有参与人的消息
func (wm *WorldManager) WorldToa(msgID uint16, subID uint16, data proto.Message) {

	for _, player := range wm.Players {
		if player != nil {
			player.SendMsg(msgID, subID, data)
		}
	}
}

// 世界同步给所有参与人的消息
func (wm *WorldManager) WorldToo(msgID uint16, subID uint16, data proto.Message, out_userName string) {

	for userName, pid := range wm.PMs {
		if p, ok := wm.Players[pid]; ok {
			if userName != out_userName {
				fmt.Println("Too -> ", userName)
				p.SendMsg(msgID, subID, data)
			}
		}
	}
}

// 给没有故障的用户转发消息，包括老师
func (wm *WorldManager) Toa_NoGz(msgID uint16, subID uint16, data proto.Message) {

	fmt.Println("Toa_NoGz->消息ID=", msgID, ",SubId=", subID)
	for _, pid := range wm.PMs {
		if p, ok := wm.Players[pid]; ok {
			if p.CDevice.Status == 1 {
				p.SendMsg(msgID, subID, data)
			}
		}
	}
}

// 给没有故障的用户转发消息，不包括老师
func (wm *WorldManager) Toa_NoGzNoTeacher(msgID uint16, subID uint16, data proto.Message) {

	fmt.Println("Toa_NoGzNoTeacher->消息ID=", msgID, ",SubId=", subID)
	for _, pid := range wm.PMs {
		if p, ok := wm.Players[pid]; ok {
			if p.CDevice.Status == 1 && p.AccountType != 1 {
				p.SendMsg(msgID, subID, data)
			}
		}
	}
}

// 给没有故障的用户转发消息，不包括老师
func (wm *WorldManager) ShotDown(msgID uint16, subID uint16, data proto.Message) {

	fmt.Println("Toa_NoGzNoTeacher->消息ID=", msgID, ",SubId=", subID)
	for _, pid := range wm.PMs {
		if p, ok := wm.Players[pid]; ok {
			if p.CDevice.Status == 1 && p.AccountType != 1 {
				p.SendMsg(msgID, subID, data)
				p.LostConnection()
			}
		}
	}
}

// 给没有故障的用户转发消息，不包括老师
func (wm *WorldManager) Toa_NoGzNoTeacher_InScene(msgID uint16, subID uint16, data proto.Message) {

	fmt.Println("Toa_NoGzNoTeacher_InScene->消息ID=", msgID, ",SubId=", subID)
	if wm.MScene != nil && len(wm.MScene.Players) > 0 {
		for username, _ := range wm.MScene.Players {
			if pid, has := wm.PMs[username]; has {
				if p, ok := wm.Players[pid]; ok {
					p.SendMsg(msgID, subID, data)
				}
			}
		}
	}

}

// 给没有故障的用户转发消息，包括老师
func (wm *WorldManager) ToTeacher(msgID uint16, subID uint16, data proto.Message) {

	username := "teacher"
	if pid, has := wm.PMs[username]; has {
		if p, ok := wm.Players[pid]; ok {
			//fmt.Println("ToTeacher->消息ID=", msgID, ",SubId=", msgSub)
			p.SendMsg(msgID, subID, data)
		}
	}
}

// 给没有故障的用户转发消息，包括老师
func (wm *WorldManager) Toa_NoGz_InScene(msgID uint16, subID uint16, data proto.Message) {

	fmt.Println("Toa_NoGz_InScene->消息ID=", msgID, ",SubId=", subID)
	if wm.MScene != nil && len(wm.MScene.Players) > 0 {
		for username, _ := range wm.MScene.Players {
			if pid, has := wm.PMs[username]; has {
				if p, ok := wm.Players[pid]; ok {
					p.SendMsg(msgID, subID, data)
				}
			}
		}
	}

	username := "teacher"
	if pid, has := wm.PMs[username]; has {
		if p, ok := wm.Players[pid]; ok {
			p.SendMsg(msgID, subID, data)
		}
	}
}

// 给没有故障的用户转发消息，包括老师
func (wm *WorldManager) Too(un string, msgID uint16, subID uint16, data proto.Message) {

	fmt.Println("Too->消息ID=", msgID, ",SubId=", subID, " un = ", un)
	if wm.MScene != nil && len(wm.MScene.Players) > 0 {
		for username, _ := range wm.MScene.Players {
			if username != un {
				if pid, has := wm.PMs[username]; has {
					if p, ok := wm.Players[pid]; ok {
						p.SendMsg(msgID, subID, data)
					}
				}
			}
		}
	}

	username := "teacher"
	if username != un {
		if pid, has := wm.PMs[username]; has {
			if p, ok := wm.Players[pid]; ok {
				p.SendMsg(msgID, subID, data)
			}
		}
	}
}

// 作品内同步给所有参与人的消息
func (wm *WorldManager) Toa(msgID uint16, subID uint16, data proto.Message) {

	//fmt.Println("Toa->消息ID=",msgID,",SubId=",msgSub)
	//wm.PrintUserList()
	if wm.MScene != nil {
		if wm.MScene.Running {
			players := wm.MScene.Players
			if players != nil && len(players) > 0 {
				for username, _ := range players {
					//if state == 1 {
					p := wm.GetPlayerByUserName(username)
					if p != nil {
						p.SendMsg(msgID, subID, data)
					}
					//}
				}
			}
		}
	}
}

//// 作品内同步给其他参与人的消息
//func (wm *WorldManager) Too(uid int32, msgID uint32, msgSub uint32, data proto.Message) {
//	//fmt.Println("Too->消息ID=",msgID,",SubId=",msgSub)
//
//	if wm.MScene != nil {
//		if wm.MScene.Running {
//			players := wm.MScene.Players
//			if players != nil && len(players) > 0 {
//				for username, _ := range players {
//					p := wm.GetPlayerByUserName(username)
//					if p != nil {
//						if p.PID != uid {
//							//if state == 1 {
//							p.SendMsg(msgID, msgSub, data)
//							//}
//						}
//					}
//				}
//			}
//		}
//	}
//}

// 从玩家信息表中移除一个玩家
func (wm *WorldManager) RemovePlayer(pID int32) {
	wm.pLock.Lock()
	delete(wm.Players, pID)
	wm.pLock.Unlock()
}

func (wm *WorldManager) RemovePlayerUName(userName string) {
	wm.pLock.Lock()
	if _, ok := wm.PMs[userName]; ok {
		delete(wm.PMs, userName)
	}
	wm.pLock.Unlock()
}

// 移除作品中的玩家(假移除，可以让用户重连进作品)
func (wm *WorldManager) RemoveWorkPlayer(userName string) {
	wm.pLock.Lock()
	//移除参与玩家
	if wm.MScene.Running {
		if _, ok := wm.MScene.Players[userName]; ok {
			wm.MScene.Players[userName] = 2 //表示当前参与玩家离线了
		}
	}
	wm.pLock.Unlock()
}

func (wm *WorldManager) PullUserToCell() {
	if wm.Players != nil && len(wm.Players) > 0 {
		for username, _ := range wm.PMs {
			wm.MScene.Players[username] = 0
		}
	}
}

func (wm *WorldManager) PullListToCell(users []string) {

	for _, userName := range users {
		//fmt.Println("userName=",userName,wm.PMs)
		if _, ok := wm.PMs[userName]; ok {
			//fmt.Println("userName1=",userName)
			wm.MScene.Players[userName] = 1
		}
	}
}

func (wm *WorldManager) UpdateReadedStatus(pid int32) {
	if wm.MScene.Running {

		p := wm.GetPlayerByPID(pid)
		if p != nil {
			if _, ok := wm.MScene.Players[p.UserName]; ok {
				wm.MScene.Players[p.UserName] = 1 //表示已经准备好了的
				fmt.Println(p.UserName, "[准备完毕]")
			}
		}
	}

}

// 准备好的人数
func (wm *WorldManager) AllReaded() bool {

	if wm.MScene.Players != nil && len(wm.MScene.Players) > 0 {
		for _, readed := range wm.MScene.Players {
			if readed == 0 {
				fmt.Println("有未准备好的成员，全部人数：", len(wm.MScene.Players))
				return false
			}
		}
		return true
	}
	return false
}

// 通过玩家ID 获取对应玩家信息
func (wm *WorldManager) GetPlayerByPID(pID int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	return wm.Players[pID]
}

// 通过玩家ID 获取对应玩家信息
func (wm *WorldManager) GetPlayerByUserName(username string) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	if uid, ok := wm.PMs[username]; ok {
		return wm.Players[uid]
	}
	return nil
}

// 获取所有玩家的信息
func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	//创建返回的player集合切片
	players := make([]*Player, 0)

	//添加切片
	for _, v := range wm.Players {
		players = append(players, v)
	}

	//返回
	return players
}

// 获取指定gID中的所有player信息
func (wm *WorldManager) GetPlayersByGID(gID int) []*Player {
	//通过gID获取 对应 格子中的所有pID
	pIDs := wm.AoiMgr.grIDs[gID].GetPlyerIDs()

	//通过pID找到对应的player对象
	players := make([]*Player, 0, len(pIDs))
	wm.pLock.RLock()
	for _, pID := range pIDs {
		players = append(players, wm.Players[int32(pID)])
	}
	wm.pLock.RUnlock()

	return players
}

// 通过玩家ID 获取对应玩家信息
func (wm *WorldManager) HasLogined(cid string) bool {

	_, ok := wm.PMs[cid]

	return ok

}

// ==================业务====================
func (wm *WorldManager) RegisterObject(data *pb.TCP_RegisterObj) {

	if _, ok := wm.MScene.JHObjects[int(data.ObjId)]; !ok {
		obj := &JHObject{
			ObjId:           int(data.ObjId),
			InteractiveType: int(data.InteractiveType),
			Tb:              int(data.Tb),
			Visiable:        int(data.Visiable),
			X:               data.X,
			Y:               data.Y,
			Z:               data.Z,
			RX:              data.RX,
			RY:              data.RY,
			RZ:              data.RZ,
		}
		wm.MScene.JHObjects[int(data.ObjId)] = obj
	}

}

func (wm *WorldManager) UpdateObjectPos(data *pb.TCP_TbObj) {

	if _, ok := wm.MScene.JHObjects[int(data.ObjId)]; ok {
		wm.MScene.JHObjects[int(data.ObjId)].X = data.X
		wm.MScene.JHObjects[int(data.ObjId)].Y = data.Y
		wm.MScene.JHObjects[int(data.ObjId)].Z = data.Z
		wm.MScene.JHObjects[int(data.ObjId)].RX = data.RX
		wm.MScene.JHObjects[int(data.ObjId)].RY = data.RY
		wm.MScene.JHObjects[int(data.ObjId)].RZ = data.RZ
	}

}

func (wm *WorldManager) UpdateObjectStatus(data *pb.Tcp_ObjectStatus) {

	if _, ok := wm.MScene.JHObjects[int(data.ObjId)]; ok {
		wm.MScene.JHObjects[int(data.ObjId)].Visiable = int(data.Status)
	}

}

func (wm *WorldManager) RegisterStep(data *pb.Tcp_Step) {

	//_id := wm.MScene.GlobalStepId+1
	//wm.MScene.GlobalStepId = _id
	obj := &Step{
		StepId:    int(data.StepId),
		StepState: int(data.StepState),
		StepDate:  data.StepDate,
		UName:     data.UName,
	}
	wm.MScene.Steps[int(data.StepId)] = obj

}

// code 0-服务器检测没人了结束作品 1-老师控制结束作品
func (wm *WorldManager) WorkFinish(code int) {
	fmt.Println("[作品结束了]code=", code)
	if wm.MScene.Running {

		wm.MScene.Running = false
	}
	wm.PrintUserList()
}

// 个人离开活动
// _state //0-中途离开 1-结束离开(结束离开如果是多人场景不重新指定主控)
func (wm *WorldManager) WorkLeave(userName string, _state int) {

	//移除参与玩家
	if wm.MScene != nil && wm.MScene.Running {
		fmt.Println("[用户离开作品]userName=", userName)
		if _, ok := wm.MScene.Players[userName]; ok {
			delete(wm.MScene.Players, userName)
		}

		//if _state != 1 {
		//	//中途离开的话需要检测
		//	if wm.MScene.MainCtroller == userName {
		//		fmt.Println("[中途离开重新指定主控]原主控=", userName)
		//		wm.ChangeMainCtroller(1)
		//	} else {
		//		wm.CheckOver(1)
		//	}
		//} else {
		//	wm.CheckOver(1)
		//}
	}

}
