package core

import (
	"fmt"
	"github.com/lyyym/zinx-wsbase/global"
	"github.com/lyyym/zinx-wsbase/release/sis_server/pb"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type Room struct {
	RID        int32  //房间ID
	TID        string //老师ID
	TPID       int32  //老师的PID
	AllCourses []pb.Sync_LoginTeacher_Info
	ClientSNum map[string]int    //学生学号记录
	Players    map[int32]*Player //当前在线的玩家集合
	pLock      sync.RWMutex      //保护Players的互斥读写机制
	CtrlFlag   string            //控制状态 0-未控制 1-控制

	ZjUid      string //用来记录临时发送数据的中间用户
	CourseID   string
	CourseMd5  string
	CourseMode string
}

/*
Player ID 生成器
*/
var RIDGen int32 = 1   //用来生成玩家ID的计数器
var RIDLock sync.Mutex //保护PIDGen的互斥机制

// 创建一个玩家对象
func NewRoom(tid string) *Room {
	//生成一个PID
	IDLock.Lock()
	ID := PIDGen
	PIDGen++
	IDLock.Unlock()

	p := &Room{
		RID:        ID,
		TID:        tid,
		Players:    make(map[int32]*Player),
		ClientSNum: make(map[string]int),
		CtrlFlag:   "0",
	}

	fmt.Println("[CreateRoom] - RID = ", tid)

	return p
}

// 提供添加一个玩家的的功能，将玩家添加进玩家信息表Players
func (rm *Room) AddPlayer(player *Player) {
	//将player添加到 世界管理器中
	rm.pLock.Lock()
	rm.Players[player.PID] = player

	if player.TID == player.CID {
		fmt.Println("老师 = ", player.CID, " 进入房间！")
		rm.TPID = player.PID
		//go rm.BroadcastTeacher()
		//改成服务启动从世界服务发送
	}

	rm.pLock.Unlock()

	fmt.Println("[AddPlayer( ", player.CID, " ) InRoom] Rid = ", rm.TID, " playerLength = ", len(rm.Players))

}

func (rm *Room) BroadcastTeacher() {

	t1 := time.NewTimer(time.Millisecond * 5000) //5s
L:
	for {
		if rm.TPID == 0 {
			break L
		}
		//fmt.Println("len(cRoom.ClientHandle), ", len(cRoom.ClientHandle), "cRoom.TeacherCli,", cRoom.TeacherCli)
		/*if len(cRoom.ClientHandle) == 0 && cRoom.TeacherCli == nil{
			goto ForEnd
		}*/
		select {
		case <-t1.C:
			t1.Reset(time.Millisecond * 5000)
			SendUdpBroadcastToStudent(rm.TID)
		}
	}

}

func SendUdpBroadcastToStudent(tid string) {
	//
	//laddrStu := net.UDPAddr{
	//	IP:   net.ParseIP(utils.GlobalObject.Host), //LocalIp(),
	//	Port: utils.GlobalObject.UdpPort,
	//}
	//
	//// 这里设置接收者的IP地址为广播地址
	//raddrStu := net.UDPAddr{
	//	IP:   net.IPv4(255, 255, 255, 255),
	//	Port: utils.GlobalObject.UdpPort,
	//}
	////fmt.Println("SendUdpBroadcast To Student raddrStu = ", laddrStu, " raddrStu = ", raddrStu)
	//connStu, err := net.DialUDP("udp", &laddrStu, &raddrStu)
	//if err != nil {
	//	println(err.Error())
	//	return
	//}
	//
	////fmt.Println("Room[",tid,"] localConfig = " , laddrStu)
	//
	//connStu.Write([]byte(tid))
	//connStu.Close()
}

func (rm *Room) SendUdpBroadcastToStudent_CloseCourse() {
	//var sData pb.Sync_Hello
	//sData.Ip = global.Object.Host
	//sData.Port = global.Object.TCPPort
	//sData.GinPort = global.Object.GinPort
	//zlog.Debugf("Broadcast Ip=%s,TcpPort=%d,GinPort = %d", sData.Ip, sData.Port, sData.GinPort)
	//global.Glog.Info("udp broadcast ", zap.String("Ip", sData.Ip), zap.Int("TcpPort", sData.Port), zap.Int("GinPort", sData.GinPort))
	// 这里设置接收者的IP地址为广播地址
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

	//fmt.Println("pBroadcast LanServer ", raddrStu, "data = ", sData)
	//x, _ := json.Marshal(sData)
	//connStu.Write(x)
	//connStu.Close()
	connStu.Write([]byte("QuitCourse"))
	connStu.Close()
}

func LocalIp() net.IP {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var ip net.IP = nil
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ip = ipnet.IP
						break
						//fmt.Println(ip,ipnet.String())
					}
				}
			}
		}
		if ip != nil {
			break
		}
	}

	return ip
}

// 从玩家信息表中移除一个玩家
func (rm *Room) RemovePlayerByPID(pID int32) {
	rm.pLock.Lock()

	if rm.TPID == pID {
		fmt.Println("老师 = ", pID, " 离开房间！")
		rm.TPID = 0
	}

	delete(rm.Players, pID)

	rm.pLock.Unlock()
}

// 获取房间的老师
func (rm *Room) GetTPlayer() *Player {
	rm.pLock.RLock()
	defer rm.pLock.RUnlock()

	if rm.TPID == 0 {
		fmt.Println("[GetTPlayer In Room]  老师不存在->Room = ", rm.TID)
		return nil
	}

	return rm.Players[rm.TPID]

}

// 获取所有玩家的信息
func (rm *Room) GetAllPlayers() []*Player {
	rm.pLock.RLock()
	defer rm.pLock.RUnlock()

	//创建返回的player集合切片
	players := make([]*Player, 0)

	//添加切片
	for _, v := range rm.Players {
		players = append(players, v)
	}

	//返回
	return players
}

// 根据课程Id 拿到课程MD5
func (cRoom *Room) GetMd5ByCourseId(courseId string) string {

	fmt.Println("courseID = ", courseId)
	fmt.Println("AllCourses = ", cRoom.AllCourses)
	for _, value := range cRoom.AllCourses {
		if value.CourseID == courseId {
			fmt.Println("当前课程MD5 ", value.Extras)
			return value.Extras
		}
	}

	//找不到从数据库在找
	var _md5 string = ""
	db := global.SqliteInst.GetDB()
	rows, _ := db.Query("select MD5 from tb_course where courseID = ?", courseId)
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&_md5); err == nil {
			//allStudent.Students = append(allStudent.Students, studentInfo)
		} else {
			log.Println("Mysql_GetTeacherSourse,", err)
		}
	}
	if len(_md5) == 0 {
		fmt.Println("找不到当前课程")
	}

	return _md5
}

func (cRoom *Room) GetCourseNameByCourseId(courseId string) string {
	for _, value := range cRoom.AllCourses {
		if value.CourseID == courseId {
			return value.Name
		}
	}
	return ""
}

// 根据课程Id 拿到课程类型
func (cRoom *Room) GetCourseNameByCourseType(courseId string) string {
	for _, value := range cRoom.AllCourses {
		if value.CourseID == courseId {
			return strings.Split(value.CourseSubType, "|")[0]
		}
	}
	return ""
}
