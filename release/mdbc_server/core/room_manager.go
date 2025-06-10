package core

import (
	"sync"
)




/*
	当前游戏世界的房间管理模块
*/
type RoomManager struct {

	Rooms map[string]*Room //当前在线的玩家集合
	pLock   sync.RWMutex      //保护Players的互斥读写机制

}


//提供一个对外的房间管理模块句柄
var RoomMgrObj *RoomManager

//提供WorldManager 初始化方法
func init() {
	RoomMgrObj = &RoomManager{
		Rooms: make(map[string]*Room),
	}
}
//
////提供添加一个玩家的的功能，将玩家添加进玩家信息表Players
//func (rm *RoomManager) CreateRoom(rid string) {
//
//	rm.pLock.Lock()
//	room := CreateRoom(rid)
//	rm.Rooms[rid] = room
//	rm.pLock.Unlock()
//}

////提供添加一个玩家的的功能，将玩家添加进玩家信息表Players
//func (rm *RoomManager) AddPlayer(player *Player) {
//
//	rm.pLock.Lock()
//	tid := player.TID
//	if _,ok := rm.Rooms[tid]; !ok {
//		room := NewRoom(tid)
//		rm.Rooms[tid] = room
//	}
//
//	rm.Rooms[tid].AddPlayer(player)
//
//	rm.pLock.Unlock()
//}
//
////从房间中获取老师
//func (rm *RoomManager) GetRoom (tid string)  *Room{
//
//	rm.pLock.RLock()
//	defer rm.pLock.RUnlock()
//	if _,ok := rm.Rooms[tid];ok {
//		return rm.Rooms[tid]
//	}else {
//		fmt.Println("[GetRoom]  RoomNotExit->Room = " , tid)
//		return nil
//	}
//
//
//
//}
//
////玩家从房间中移出
//func (rm *RoomManager) RemovePlayerByPID(tid string , pID int32 , cid string) {
//	rm.pLock.Lock()
//	if _,ok := rm.Rooms[tid];ok {
//		rm.Rooms[tid].RemovePlayerByPID(pID)
//	}else {
//		fmt.Println("[RemovePlayer In Room] player = " , cid , " pid = " , pID , " RoomNotExit->Room = " , tid)
//	}
//	rm.pLock.Unlock()
//}
//
//
////从房间中获取老师
//func (rm *RoomManager) GetTPlayer (tid string)  *Player{
//
//	rm.pLock.RLock()
//	defer rm.pLock.RUnlock()
//	if _,ok := rm.Rooms[tid];ok {
//		return rm.Rooms[tid].GetTPlayer()
//	}else {
//		fmt.Println("[GetTPlayer In Room]  RoomNotExit->Room = " , tid)
//		return nil
//	}
//
//
//
//}
//
//
////从房间中获取玩家列表
//func (rm *RoomManager) GetAllPlayers (tid string)  []*Player{
//
//	rm.pLock.RLock()
//	defer rm.pLock.RUnlock()
//	if _,ok := rm.Rooms[tid];ok {
//		return rm.Rooms[tid].GetAllPlayers()
//	}else {
//		fmt.Println("[GetAllPlayers In Room]  RoomNotExit->Room = " , tid)
//		return nil
//	}
//
//
//
//}



