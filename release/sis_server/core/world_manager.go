package core

import (
	"encoding/json"
	"fmt"
	"github.com/lyyym/zinx-wsbase/global"
	"github.com/lyyym/zinx-wsbase/release/sis_server/pb"
	"net"
	"sync"
	"time"
)

/*
当前游戏世界的总管理模块
*/
type WorldManager struct {
	AoiMgr  *AOIManager       //当前世界地图的AOI规划管理器
	Players map[int32]*Player //当前在线的玩家集合
	PMs     map[string]int32  //记录登录的玩家，防止重复登录
	AppID   string            //用户登录标识，无此标识则不允许登录
	pLock   sync.RWMutex      //保护Players的互斥读写机制
}

// 提供一个对外的世界管理模块句柄
var WorldMgrObj *WorldManager

// 提供WorldManager 初始化方法
func init() {
	WorldMgrObj = &WorldManager{
		Players: make(map[int32]*Player),
		AoiMgr:  NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
		PMs:     make(map[string]int32),
	}
}

// 提供添加一个玩家的的功能，将玩家添加进玩家信息表Players
func (wm *WorldManager) BindPlayer(player *Player) {
	//将player添加到 世界管理器中
	wm.pLock.Lock()

	if _, ok := wm.Players[player.PID]; ok {

		//3. 记录真实玩家(已经发生登录的玩家)
		if _, has := wm.PMs[player.CID]; has {
			fmt.Println("已经有相同账号记录在服务器")
			//1. 有可能被顶号
			//2. 有可能重复绑定
			// 暂时不处理
		}
		wm.PMs[player.CID] = player.PID

		//if player.AccountType == 1 {
		//	//记录老师状态
		//	wm.TUserName = player.UserName
		//	wm.TUid = player.PID
		//	//fmt.Println("[记录老师状态]")
		//}
	}

	wm.pLock.Unlock()
	fmt.Println("Players = ", wm.PMs)
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

// 从玩家信息表中移除一个玩家
func (wm *WorldManager) RemovePlayerByPID(pID int32) {
	wm.pLock.Lock()
	delete(wm.PMs, wm.Players[pID].CID)
	delete(wm.Players, pID)
	wm.pLock.Unlock()
}

// 通过玩家ID 获取对应玩家信息
func (wm *WorldManager) GetPlayerByPID(pID int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	return wm.Players[pID]
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
