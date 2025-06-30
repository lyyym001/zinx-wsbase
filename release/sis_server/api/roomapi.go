package api

import (
	"encoding/json"
	"fmt"
	"github.com/lyyym/zinx-wsbase/global"
	"github.com/lyyym/zinx-wsbase/release/sis_server/core"
	"github.com/lyyym/zinx-wsbase/release/sis_server/pb"
	"github.com/lyyym/zinx-wsbase/ziface"
	"github.com/lyyym/zinx-wsbase/znet"
	"log"
	"strings"
	"time"
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

	case 10001: //老师发送课程信息给局域网服务器

		ra.Handle_TeacherData(player, request.GetData())

		//给老师回传 年级分类信息列表数据
		//if player.TID == player.CID {
		//	aa.Handle_NjFl(player)
		//	aa.Handle_KcFl(player)
		//	aa.Handle_KcDatas(player)
		//}
		break
	case 10002: //删除一个本地数据库的的课程
		ra.Handle_DeleteCourse(player, request.GetData())
		break
	case 10003: //更新一个本地课程数据
		ra.Handle_UpdateCourse(player, request.GetData())
		break
	case 10004: //添加一个新的本地课程数据
		ra.Handle_AddCourse(player, request.GetData())
		break
	case 10005: //控制解控
		ra.Handle_Control(player, request.GetData())
		break
	case 10006: //打开课程
		ra.Handle_OpenCourse(player, request.GetData())
		break
	case 10007: //学生进入课程
		ra.Handle_EnterCourse(player, request.GetData())
		break
	case 10008: //学生离开课程
		ra.Handle_LeaveCourse(player, request.GetData())
		break
	case 10009:
		ra.Handle_CloseCourse(player, request.GetData())
		break
	case 10010:
		ra.Handle_RequestData(player, request.GetData())
		break
	case 10011:
		ra.Handle_ResponseData(player, request.GetData())
		break
	case 10012:
		ra.Handle_ReplayCourse(player, request.GetData())
		break
	case 10013:
		ra.Handle_UpdateC(player, request.GetData())
		break
	case 10014:
		ra.Handle_UpdateOver(player, request.GetData())
		break
	case 10015:
		ra.Handle_UpdateErr(player, request.GetData())
		break
	case 10016:
		ra.Handle_SendNo(player, request.GetData())
		break
	case 10017:
		ra.Handle_OpenNo(player, request.GetData())
		break
	case 10018:
		ra.Handle_GetStus(player, request.GetData())
		break
	case 10019:
		ra.Handle_GetStusWithCondition(player, request.GetData())
		break
	case 10020:
		ra.Handle_onStudentRecord(player, request.GetData())
		break
	case 10021:
		ra.Handle_onAddStudent(player, request.GetData())
		break
	case 10022:
		ra.Handle_onStudyRecordByCourse(player, request.GetData())
		break
	case 10023:
		ra.Handle_onUpdatePassword(player, request.GetData())
		break
	case 10024:
		ra.Handle_onClientShutdown(player, request.GetData())
		break
	case 10025:
		ra.Handle_onClientBatteryAndSpace(player, request.GetData())
		break
	case 10026:
		ra.Handle_onVideoPause(player, request.GetData())
		break
	case 10027:
		ra.Handle_onCourseNodeUpdate(player, request.GetData())
		break
	case 10028:
		ra.Handle_onGamePause(player, request.GetData())
		break
	case 10029:
		ra.Handle_DeleteStu(player, request.GetData())
		break
	case 10030:
		ra.Handle_AddUser(player, request.GetData())
		break
	case 10031:
		ra.Handle_onStudentData(player, request.GetData())
		break
	case 10032:
		ra.Handle_onClientHuyan(player, request.GetData())
		break
	case 10033:
		ra.Handle_onClientHuyanClose(player, request.GetData())
		break
	case 10034: //开启录屏
		ra.Handle_onClientLpStart(player, request.GetData())
		break
	case 10035: //关闭录屏
		ra.Handle_onClientLpClose(player, request.GetData())
		break
	case 10036: //开启/关闭 静音/黑屏/护眼
		ra.Handle_onTurnOnOrOffStatus(player, request.GetData())
		break
		//
	}

}

// 记录老师的课程信息
func (aa *RoomApi) Handle_TeacherData(p *core.Player, data []byte) {

	//
	rm := core.RoomMgrObj.GetRoom(p.TID)
	msg := &pb.Sync_LoginTeacher_Send{}
	json.Unmarshal(data, msg)
	rm.AllCourses = msg.AllCourses

	fmt.Println("[2-20001]记录老师的课程信息 课程数量 = ", len(msg.AllCourses))
}

// 老师删除一个第三方课程 237 - > 2 - 10002
func (aa *RoomApi) Handle_DeleteCourse(p *core.Player, data []byte) {

	msg := &pb.DeleteCourse{}
	json.Unmarshal(data, msg)

	// (2) 更新
	db := global.SqliteInst.GetDB()
	stmt, err := db.Prepare("DELETE FROM tb_course WHERE ID = ?")
	if err != nil {
		fmt.Println("Sqlite Handle_DeleteCourse Update DB Err")
	} else {
		defer stmt.Close()
		_, err := stmt.Exec(msg.ID)
		//affectNum, err := result.RowsAffected()
		if err != nil {
			fmt.Println("Sqlite Err = " + err.Error())
		}
		//fmt.Println("update affect rows is ", affectNum)
	}

	fmt.Println("[2-20002]老师删除一个第三方课程 data = ", msg)
}

// 学生更新课程，通知老师
func (aa *RoomApi) Handle_UpdateCourse(p *core.Player, data []byte) {

	CourseData := &pb.UpdateCourse{}
	json.Unmarshal(data, CourseData)

	// (2) 更新
	db := global.SqliteInst.GetDB()
	stmt, err := db.Prepare("UPDATE tb_course SET md5=?,gameUrl=?,resVersion=?,iconName=?,courseName=?,courseType=? WHERE courseID = ?")
	if err != nil {
		fmt.Println("Sqlite Handle_UpdateCourse Update DB Err : ", err.Error())
	} else {
		defer stmt.Close()
		_, err := stmt.Exec(CourseData.Md5, CourseData.GameUrl, CourseData.ResVersion, CourseData.GameIcoPath, CourseData.GameName, CourseData.VersionNo, CourseData.CourseID)
		//affectNum, err := result.RowsAffected()
		if err != nil {
			fmt.Println("Sqlite Handle_UpdateCourse DB Err 1")
		}
		//fmt.Println("update affect rows is ", affectNum)
	}

	fmt.Println("[2-20003]老师更新一个第三方课程 data = ", CourseData)
}

// 添加一个本地课程信息
func (aa *RoomApi) Handle_AddCourse(p *core.Player, data1 []byte) {

	data := &pb.CoursewareData{}
	json.Unmarshal(data1, data)
	var sid string = ""
	call := &pb.Sync_CleverCall{Code: 0}
	db := global.SqliteInst.GetDB()
	rows1, err1 := db.Query("select courseID from tb_course where courseID = ? ", data.CourseID)
	defer rows1.Close()
	if err1 == nil {
		for rows1.Next() {
			if err := rows1.Scan(&sid); err == nil {

			} else {
				log.Println("Handle_onGetReportBySname,", err)
			}
		}
	}
	fmt.Println("[3-20004]添加灵创课程 sid = ,", sid)
	if len(sid) > 0 {
		fmt.Println("[3-20004]添加灵创课程 err = 作品已经存在,", data.CourseID)
		data1, _ := json.Marshal(call)
		p.SendMsg(2, 10004, data1)
		return //作品已经存在
	}

	// (2) 更新
	fmt.Println("data = ", data.ThirdMsg, " - md5 = ", data.Md5)
	stmt, err := db.Prepare("insert into tb_course(courseName,iconName,courseID,courseType,courseOwner,inCourseType,inCourseTypeSort,thirdType,md5,gameUrl,resVersion) values(?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Println("Sqlite Handle_UpdateCourse Update DB Err : 1 = ", err.Error())
	} else {

		res, err := stmt.Exec(data.Name, data.IconName, data.CourseID, data.CourseType, data.CourseOwner, data.InCourseType, data.InCourseTypeSort, data.ThirdType, data.ThirdMsg, data.GameUrl, data.ResVersion)
		//affectNum, err := result.RowsAffected()
		if err != nil {
			fmt.Println("Sqlite Handle_UpdateCourse DB Err : 2 = ", err.Error())
		} else {
			//拿到插入的数据库ID
			//res, _ := stmt.Exec(data.CourseID)
			fmt.Println("res = ", res)
			_id, _ := res.LastInsertId()
			call.Code = 1
			call.Id = _id
			fmt.Println("[3-20004]添加灵创课程 , 作品已添加成功 ， ID = ", _id, call)
			data1, _ := json.Marshal(call)
			p.SendMsg(2, 10004, data1)
		}
		//fmt.Println("update affect rows is ", affectNum)
	}

	fmt.Println("[2-20004]添加一个本地课程信息 data = ", data)
}

// 控制解控(老师发来的消息)
func (aa *RoomApi) Handle_Control(p *core.Player, data []byte) {

	CRoom := core.RoomMgrObj.GetRoom(p.TID)
	if CRoom != nil {
		CRoom.SendUdpBroadcastToStudent_CloseCourse()
	}

	msg := &pb.Sync_TeacherControlData{}
	json.Unmarshal(data, msg)

	//教室
	Room := core.RoomMgrObj.GetRoom(p.TID)
	if Room != nil {
		Room.CtrlFlag = msg.IsTeacherControl
	}

	//学生列表
	players := core.RoomMgrObj.GetAllPlayers(p.TID)
	if players != nil {
		for _, player := range players {
			player.SendMsg(2, 10005, data)
		}
	}

	fmt.Println("[2-20005]控制解控 msg = ", msg)
}

// 老师打开课程通知学生
func (aa *RoomApi) Handle_OpenCourse(p *core.Player, data []byte) {

	msg := &pb.Sync_GetCourse{}
	json.Unmarshal(data, msg)

	Room := core.RoomMgrObj.GetRoom(p.TID)
	if Room != nil {
		var sendData pb.Sync_SendCourse
		sendData.Md5 = Room.GetMd5ByCourseId(msg.CourseID)
		sendData.CourseID = msg.CourseID
		sendData.Mode = msg.Mode
		Room.CourseID = sendData.CourseID
		Room.CourseMd5 = sendData.Md5
		Room.CourseMode = sendData.Mode

		db := global.SqliteInst.GetDB()

		stmt, err := db.Prepare("UPDATE tb_user SET isTeacherClose=?")
		if err != nil {
			fmt.Println("Sqlite Handle_CloseCourse Update DB Err - 1 : ", err.Error())
		} else {
			defer stmt.Close()
			_, err := stmt.Exec(1, "0")
			//affectNum, err := result.RowsAffected()
			if err != nil {
				fmt.Println("Sqlite Handle_CloseCourse affect DB Err")
			}
			//fmt.Println("update affect rows is ", affectNum)
		}

		data1, _ := json.Marshal(sendData)
		fmt.Println("md5 = ", sendData.Md5, " courseID = ", sendData.CourseID)
		//学生列表

		//find := strings.Contains(sendData.Md5, ".")

		players := core.RoomMgrObj.GetAllPlayers(p.TID)
		if players != nil {
			for _, player := range players {
				player.SendMsg(2, 10006, data1)

			}
			//for _, player := range players {
			//
			//	if find {
			//		fmt.Println("player[", player.PID, "] lost")
			//		player.LostConnection()
			//		player.Conn.Stop()
			//	}
			//}

		}

	}

	fmt.Println("[2-20006]打开课程 msg = ", data)
}

// 学生进入课程通知老师
func (aa *RoomApi) Handle_EnterCourse(p *core.Player, data []byte) {

	msg := &pb.Sync_SelectCourseware_Send{}
	json.Unmarshal(data, msg)

	p.CourseSetupDate = time.Now().Unix()
	p.CourseCore = 0
	p.CourseId = msg.CourseID
	p.CourseMode = msg.CourseMode

	fmt.Println("进入课程 = ", msg)

	//学生列表
	player := core.RoomMgrObj.GetTPlayer(p.TID)
	if player != nil {
		var sendData pb.StuName
		sendData.StuUserName = p.CID
		//序列化数据
		data1, _ := json.Marshal(sendData)
		player.SendMsg(2, 10007, data1)
	}

	fmt.Println("[2-20007]进入课程 msg = ", data)
}

// 学生离开课程通知老师
func (aa *RoomApi) Handle_LeaveCourse(p *core.Player, data []byte) {

	msg := &pb.Sync_LeaveData{}
	json.Unmarshal(data, msg)
	//学生列表
	player := core.RoomMgrObj.GetTPlayer(p.TID)
	if player != nil {
		var sendData pb.StuName
		sendData.StuUserName = p.CID
		//序列化数据
		data1, _ := json.Marshal(sendData)
		player.SendMsg(2, 10008, data1)
		fmt.Println("StuSendToTeacher", player.CID, time.Now())
	} else {
		fmt.Println("学生离开课程[", msg.LeaveType, "] 老师不存在。")
	}

	fmt.Println("离开课程 = ", p.SNum, " p.CourseId = ", p.CourseId)

	if p.SNum != "" && p.CourseId != "" {

		cRoom := core.RoomMgrObj.GetRoom(p.TID)
		if cRoom != nil {
			db := global.SqliteInst.GetDB()
			cname := cRoom.GetCourseNameByCourseId(p.CourseId)
			ctype := cRoom.GetCourseNameByCourseType(p.CourseId)

			var fuName string = cname
			if ctype == "1" {
				if p.CourseMode == "1" {
					fuName += "-教学模式"
				} else if p.CourseMode == "2" {
					fuName += "-实践模式"
				} else {
					fuName += "-考核模式"
				}
			}

			var dbdata pb.DB_SaveData
			dbdata.Tid = p.TID
			dbdata.SNum = p.SNum
			dbdata.CourseId = p.CourseId
			dbdata.CourseName = fuName
			dbdata.CourseMode = p.CourseMode
			dbdata.CourseType = ctype
			dbdata.StuTimeLong = time.Now().Unix() - p.CourseSetupDate
			dbdata.Score = p.CourseCore
			dbdata.Ability = p.CourseAbility
			dbdata.SetupDate = p.CourseSetupDate
			dbdata.LeaveType = msg.LeaveType

			//fmt.Println("20008:dbdata = ",dbdata)

			stmt, err := db.Prepare("INSERT INTO tb_recore(Tid, CourseId, CourseName, CourseMode, CourseType, StuTimeLong, SNum, Score, SetupDate, LeaveType, Ability) VALUES(?,?,?,?,?,?,?,?,?,?,?)")
			defer stmt.Close()
			if err == nil {
				timeRecore := time.Unix(dbdata.SetupDate, 0).Format("2006-01-02 15:04:05")
				stmt.Exec(dbdata.Tid, dbdata.CourseId, dbdata.CourseName, dbdata.CourseMode, dbdata.CourseType, dbdata.StuTimeLong, dbdata.SNum, dbdata.Score, timeRecore, dbdata.LeaveType, dbdata.Ability)
			}

		}
	}

	fmt.Println("[2-20008]离开课程 msg = ", len(data))
}

// 老师关闭课程，通知学生
func (aa *RoomApi) Handle_CloseCourse(p *core.Player, data []byte) {

	msg := &pb.Sync_CloseCourse{}
	json.Unmarshal(data, msg)

	//发送UDP - 关闭课程
	CRoom := core.RoomMgrObj.GetRoom(p.TID)
	if CRoom != nil {
		CRoom.SendUdpBroadcastToStudent_CloseCourse()
		CRoom.SendUdpBroadcastToStudent_CloseCourse()
		CRoom.SendUdpBroadcastToStudent_CloseCourse()
		CRoom.SendUdpBroadcastToStudent_CloseCourse()
		CRoom.SendUdpBroadcastToStudent_CloseCourse()
	}

	players := core.RoomMgrObj.GetAllPlayers(p.TID)
	if msg.StuAccountID == "" {

		//全部结束课程
		db := global.SqliteInst.GetDB()

		stmt, err := db.Prepare("UPDATE tb_user SET isTeacherClose=?")
		if err != nil {
			fmt.Println("Sqlite Handle_CloseCourse Update DB Err")
		} else {
			defer stmt.Close()
			_, err := stmt.Exec(1, "1")
			//affectNum, err := result.RowsAffected()
			if err != nil {
				fmt.Println("Sqlite Handle_CloseCourse affect DB Err")
			}
			//fmt.Println("update affect rows is ", affectNum)
		}
		//butflySql.Mysql_UpdateStudentCourseCloseState("1")
		//学生列表

		if players != nil {
			for _, player := range players {
				player.SendMsg(2, 10009, data)
				fmt.Println("TeacherSendToStu", player.CID, time.Now())
			}
		} else {
			fmt.Println("老师关闭课程 教室没有学生。")
		}

	} else { //单独结束课程

		cRoom := core.RoomMgrObj.GetRoom(p.TID)
		if cRoom != nil {

			if players != nil {
				for _, player := range players {
					if player.CID == msg.StuAccountID {
						player.SendMsg(2, 10009, data)
						break
					}
				}
			} else {
				fmt.Println("老师关闭课程 教室没有学生。")
			}

		}
	}

	fmt.Println("[2-20009]老师端结束课程 msg = ", data)
}

// 学生端请求本地数据
func (aa *RoomApi) Handle_RequestData(p *core.Player, data []byte) {

	player := core.RoomMgrObj.GetTPlayer(p.TID)
	if player != nil {
		Room := core.RoomMgrObj.GetRoom(p.TID)
		if Room != nil {
			Room.ZjUid = p.CID
			player.SendMsg(2, 10010, data)
		}
	}
	fmt.Println("[2-20010]学生端请求本地数据 msg = ", data)
}

// 老师端回执本地数据
func (aa *RoomApi) Handle_ResponseData(p *core.Player, data []byte) {

	msg := &pb.Sync_GetLocalGradeTypeList{}
	json.Unmarshal(data, msg)

	players := core.RoomMgrObj.GetAllPlayers(p.TID)
	if players != nil {
		Room := core.RoomMgrObj.GetRoom(p.TID)
		if Room != nil {
			for _, player := range players {
				if player.CID == Room.ZjUid {
					player.SendMsg(2, 10011, data)
				}
			}
		}
	}
	fmt.Println("[2-20011]老师端回执本地数据 msg = ", len(data), msg)
}

// 老师端重启课程
func (aa *RoomApi) Handle_ReplayCourse(p *core.Player, data []byte) {

	msg := &pb.Sync_TInputSnum{}
	json.Unmarshal(data, msg)

	Room := core.RoomMgrObj.GetRoom(p.TID)
	if Room != nil {
		//要重启的课程数据
		var sendData pb.Sync_SendCourse
		sendData.CourseID = Room.CourseID
		sendData.Md5 = Room.CourseMd5
		sendData.Mode = Room.CourseMode
		data1, _ := json.Marshal(sendData)

		players := core.RoomMgrObj.GetAllPlayers(p.TID)
		if players != nil {

			for _, player := range players {
				if msg.StuAccountID == "" {
					player.SendMsg(2, 10012, data1)
				} else if msg.StuAccountID == player.CID {
					player.SendMsg(2, 10012, data1)
					break
				}

			}
		}
	}

	fmt.Println("[2-20012]老师端重启课程 msg = ", data)
}

// 学生更新课程通知老师
func (aa *RoomApi) Handle_UpdateC(p *core.Player, data []byte) {

	var sendData pb.Sync_TInputSnum
	sendData.StuAccountID = p.CID

	player := core.RoomMgrObj.GetTPlayer(p.TID)
	if player != nil {
		data1, _ := json.Marshal(sendData)
		player.SendMsg(2, 10013, data1)
	}

	fmt.Println("[2-20013]学生更新课程通知老师 sendData = ", sendData)
}

// 学生更新课程成功通知老师
func (aa *RoomApi) Handle_UpdateOver(p *core.Player, data []byte) {

	var sendData pb.Sync_TInputSnum
	sendData.StuAccountID = p.CID

	player := core.RoomMgrObj.GetTPlayer(p.TID)
	if player != nil {
		data1, _ := json.Marshal(sendData)
		player.SendMsg(2, 10014, data1)
	}

	fmt.Println("[2-20014]学生更新课程[成功]通知老师 sendData = ", sendData)
}

// 学生更新课程失败通知老师
func (aa *RoomApi) Handle_UpdateErr(p *core.Player, data []byte) {

	var sendData pb.Sync_TInputSnum
	sendData.StuAccountID = p.CID

	player := core.RoomMgrObj.GetTPlayer(p.TID)
	if player != nil {
		data1, _ := json.Marshal(sendData)
		player.SendMsg(2, 10015, data1)
	}

	fmt.Println("[2-20015]学生更新课程[失败]通知老师 sendData = ", sendData)
}

// 学生发送学号给老师
func (aa *RoomApi) Handle_SendNo(p *core.Player, data []byte) {

	getData := &pb.Sync_SInputSnum{}
	json.Unmarshal(data, getData)
	fmt.Println("Handle_SendNo", getData.SNum)
	//学生输入学号
	var sendData pb.Sync_SNumReturn
	var sendData1 pb.Sync_SNumToTeacher

	cRoom := core.RoomMgrObj.GetRoom(p.TID)
	if cRoom != nil {

		//去验证一下学号是否存在
		if pid, ok := cRoom.ClientSNum[getData.SNum]; ok {
			//重复输入
			sendData.Result = "srep"
			x, _ := json.Marshal(sendData)

			p.SendMsg(2, 10016, x)

			fmt.Println("R", cRoom.TID, "学生", p.CID, "ID", p.PID, "学号重复输入--", pid, "已经输入这个学号")
		} else {

			var pName string = ""

			db := global.SqliteInst.GetDB()
			rows, _ := db.Query("select pname from tb_snum where Tid = ? and snum = ?", p.TID, getData.SNum)
			defer rows.Close()
			for rows.Next() {
				if err := rows.Scan(&pName); err == nil {
					//allStudent.Students = append(allStudent.Students, studentInfo)
				} else {
					log.Println("Mysql_GetTeacherSourse,", err)
				}
			}

			fmt.Println("cli.Tid    ", p.TID, getData.SNum, pName)
			if pName != "" {
				//
				cRoom.ClientSNum[getData.SNum] = int(p.PID)
				sendData1.SNum = getData.SNum
				sendData1.PName = pName
				sendData1.StuUserName = p.CID

				p.SNum = getData.SNum
				p.PName = pName

				x, _ := json.Marshal(sendData1)
				playerT := core.RoomMgrObj.GetTPlayer(p.TID)
				if playerT != nil {
					playerT.SendMsg(2, 10016, x)
				}
				sendData.Result = "success"
				x1, _ := json.Marshal(sendData)

				p.SendMsg(2, 10016, x1)

				//fmt.Println("学号验证成功", cli.UserName, pname, getData.SNum)
				fmt.Println("学生", p.CID, "ID", p.PID, "学号验证成功", pName, getData.SNum)
			} else {
				//学号不存在
				//log.Println("R", cRoom.Rid, "学生", cli.UserName, "ID", cli.CID, " 学号不存在")
				sendData.Result = "snull"
				x, _ := json.Marshal(sendData)

				p.SendMsg(2, 10016, x)

				//fmt.Println("学号不存在", cli.UserName)
				fmt.Println("R", p.TID, "学生", p.CID, "ID", p.PID, " 输入学号不存在")
			}

		}

	}

	fmt.Println("[2-20016]学生发送学号给老师-通知老师 sendData = ", sendData)
}

// 老师通知学生输学号
func (aa *RoomApi) Handle_OpenNo(p *core.Player, data []byte) {

	getData := &pb.Sync_TInputSnum{}
	json.Unmarshal(data, getData)
	cRoom := core.RoomMgrObj.GetRoom(p.TID)

	players := core.RoomMgrObj.GetAllPlayers(p.TID)
	if getData.StuAccountID == "" {

		if cRoom != nil {
			cRoom.ClientSNum = make(map[string]int)
		}

		if players != nil {
			for _, player := range players {
				player.SendMsg(2, 10017, data)
			}
		} else {
			fmt.Println("老师通知学生输学号 教室没有学生。")
		}

	} else { //单独结束课程

		if cRoom != nil {

			if players != nil {
				for _, player := range players {
					if player.CID == getData.StuAccountID {
						delete(cRoom.ClientSNum, player.SNum)
						player.SendMsg(2, 10017, data)
						break
					}
				}
			} else {
				fmt.Println("老师通知学生输学号 教室没有学生。")
			}

		}
	}

	fmt.Println("[2-20017]老师通知学生输学号 sendData = ", getData)
}

// 老师获取学生列表
func (aa *RoomApi) Handle_GetStus(p *core.Player, data []byte) {

	var student pb.StudentData
	var allStudent []pb.StudentData

	db := global.SqliteInst.GetDB()
	rows, err := db.Query("select snum,pname,class from tb_snum where Tid=?", p.TID)
	if err != nil {
		fmt.Println("Sqlite Handle_GetStus Query DB Err", err.Error())
	} else {
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&student.Snum, &student.Sname, &student.Sclass); err == nil {
				allStudent = append(allStudent, student)
			} else {
				log.Println("Mysql_GetLocalGradeTypeData,", err)
			}
		}
	}

	//数据
	x, _ := json.Marshal(allStudent)
	p.SendMsg(2, 10018, x)
	fmt.Println("[2-20018]老师获取学生列表 student = ", allStudent)

}

// 按条件获得当前老师某个学生信息
func (aa *RoomApi) Handle_GetStusWithCondition(p *core.Player, data []byte) {

	var sData pb.StudentData
	getData := &pb.GetConditionData{}
	json.Unmarshal(data, getData)

	db := global.SqliteInst.GetDB()
	rows, err := db.Query("select snum,pname,class from tb_snum where Tid=? and (snum=? or pname=?)", p.TID, getData.Condition, getData.Condition)
	if err != nil {
		fmt.Println("Sqlite Handle_GetStusWithCondition Query DB Err")
	} else {
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&sData.Snum, &sData.Sname, &sData.Sclass); err == nil {

			} else {
				log.Println("Handle_GetStusWithCondition,", err)
			}
		}
	}

	//数据
	x, _ := json.Marshal(sData)
	p.SendMsg(2, 10019, x)
	fmt.Println("[2-20019]老师获取学生列表 student = ", sData)

}

// 老师获取某个学生的考试记录
func (aa *RoomApi) Handle_onStudentRecord(p *core.Player, data []byte) {

	var sData pb.AllStudentRecordData
	getData := &pb.Sync_SInputSnum{}
	json.Unmarshal(data, getData)

	fmt.Println("[2-20020]Handle_onStudentRecord sum = ", getData.SNum, " tid = ", p.TID)
	var studentRecord pb.StudentRecordData
	var allStudentRecord []pb.StudentRecordData
	db := global.SqliteInst.GetDB()
	rows, err := db.Query("select ts.pname,ts.class,tr.CourseId,tr.CourseName,tr.SetupDate,tr.CourseMode,tr.StuTimeLong,tr.Score,tr.Ability from tb_recore tr inner join tb_snum ts where tr.SNum=? and tr.Tid=? and tr.tid = ts.Tid order by tr.CreateDate desc limit 10", getData.SNum, p.TID)
	if err != nil {
		fmt.Println("Sqlite Handle_onStudentRecord Query DB Err")
	} else {
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&studentRecord.StudentSnum, &studentRecord.StudentClass, &studentRecord.CourseID, &studentRecord.CourseName, &studentRecord.StudyTime, &studentRecord.StudyMode, &studentRecord.StudyTotalTime, &studentRecord.StudyScore, &studentRecord.StudyAbility); err == nil {
				allStudentRecord = append(allStudentRecord, studentRecord)
				fmt.Println("[2-20020]Handle_onStudentRecord studentRecord = ", studentRecord)
			} else {
				fmt.Println("[2-20020]Handle_GetStusWithCondition,", err)
			}
		}
	}

	//数据
	sData.AllRecordData = allStudentRecord
	x, _ := json.Marshal(sData)
	p.SendMsg(2, 10020, x)
	fmt.Println("[2-20020]Handle_onStudentRecord student = ", sData)

}

// 老师端添加学生
func (aa *RoomApi) Handle_onAddStudent(p *core.Player, data []byte) {

	getData := &pb.StudentData{}
	json.Unmarshal(data, getData)

	db := global.SqliteInst.GetDB()
	stmt, err := db.Prepare("insert into tb_snum(snum, pname, Tid, class) values(?,?,?,?)")
	defer stmt.Close()
	if err == nil {
		stmt.Exec(getData.Snum, getData.Sname, p.TID, getData.Sclass)
	} else {
		log.Println("Mysql_AddStudent,", err)
	}

	//数据
	fmt.Println("[2-20021]老师端添加学生 student = ", getData.Sname, getData.Snum, getData.Sclass)

}

// 根据课程获得学生考试记录
func (aa *RoomApi) Handle_onStudyRecordByCourse(p *core.Player, data []byte) {

	var setupDate, courseMode string

	var sData pb.AllStudyInfoData
	getData := &pb.Sync_SCourseId{}
	json.Unmarshal(data, getData)

	cRoom := core.RoomMgrObj.GetRoom(p.TID)
	if cRoom != nil {
		courseName := cRoom.GetCourseNameByCourseId(getData.CourseID)

		fmt.Println("cRoom.AllCourses = ", cRoom.AllCourses, "getData.CourseID = ", getData.CourseID, "getData.CourseID = ", getData.CourseID, " p.TID = ", p.TID)

		var studentRecord pb.SingleStudyInfoData
		var allStudentRecord []pb.SingleStudyInfoData
		db := global.SqliteInst.GetDB()
		rows, err := db.Query("select SetupDate,CourseMode from tb_recore where CourseId=? and Tid=? group by SetupDate order by CreateDate desc limit 10", getData.CourseID, p.TID)
		if err != nil {
			fmt.Println("Sqlite Handle_onStudentRecord Query DB Err")
		} else {
			defer rows.Close()
			for rows.Next() {
				if err := rows.Scan(&setupDate, &courseMode); err == nil {
					setupDate = strings.Replace(setupDate, "T", " ", -1)
					setupDate = strings.Replace(setupDate, "Z", "", -1)
					fmt.Println("setupDate = ", setupDate)
					var studyTotalTime, studyScore int
					var pname, snum, sclass, courseId, studyTime, studyMode, studyAbility string
					var srData pb.StudentRecordData
					var allSrData []pb.StudentRecordData
					rows, err := db.Query("select ts.snum,ts.pname,ts.class,tr.CourseId,tr.SetupDate,tr.CourseMode,tr.StuTimeLong,tr.Score,tr.Ability from tb_recore tr,tb_snum ts where tr.SetupDate=? and tr.Tid=? and ts.SNum=tr.SNum", setupDate, p.TID)
					defer rows.Close()
					if err == nil {
						for rows.Next() {
							if err := rows.Scan(&snum, &pname, &sclass, &courseId, &studyTime, &studyMode, &studyTotalTime, &studyScore, &studyAbility); err == nil {
								srData.StudentSnum = snum
								srData.StudentName = pname
								srData.StudentClass = sclass
								srData.CourseID = courseId
								srData.StudyTime = studyTime
								srData.CourseName = courseName
								srData.StudyMode = studyMode
								srData.StudyTotalTime = studyTotalTime
								srData.StudyScore = studyScore
								srData.StudyAbility = studyAbility
								allSrData = append(allSrData, srData)
							}
						}
					}
					if allSrData != nil {
						studentRecord.StudyTime = setupDate
						studentRecord.StudyCourseMode = courseMode
						studentRecord.StudentRecordData = allSrData
						allStudentRecord = append(allStudentRecord, studentRecord)
					}

				} else {
					fmt.Println("Handle_onStudyRecordByCourse,", err)
				}
			}
		}

		//
		sData.StudyInfoData = allStudentRecord
		x, _ := json.Marshal(sData)
		p.SendMsg(2, 10022, x)
	}

	//数据
	fmt.Println("[2-20022]根据课程获得学生考试记录 sData = ", sData)

}

// 学生更改密码
func (aa *RoomApi) Handle_onUpdatePassword(p *core.Player, data []byte) {

	getData := &pb.UpdatePasswordData{}
	json.Unmarshal(data, getData)

	db := global.SqliteInst.GetDB()
	stmt, err := db.Prepare("UPDATE account_config SET password=? where account=?")
	defer stmt.Close()
	if err == nil {
		stmt.Exec(p.CID, getData.Password)
	} else {
		log.Println("Mysql_AddStudent,", err)
	}

	//数据
	fmt.Println("[2-20023]学生更改密码 student = ", getData)

}

// 老师端控制学生端关机
func (aa *RoomApi) Handle_onClientShutdown(p *core.Player, data []byte) {

	players := core.RoomMgrObj.GetAllPlayers(p.TID)
	if players != nil {
		for _, player := range players {
			fmt.Println("关机 cid = ", player.CID)
			player.SendMsg(2, 10024, data)
		}
	}

	//数据
	fmt.Println("[2-20024]老师端控制学生端关机 data = ", data)

}

// 221 - 服务器转发给老师端电量存储空间信息
func (aa *RoomApi) Handle_onClientBatteryAndSpace(p *core.Player, data []byte) {

	rData := &pb.Sync_BatteryAndSpaceData{}
	json.Unmarshal(data, rData)

	var sData pb.Sync_BatteryAndSpaceData_Send
	sData.StuName = p.CID
	sData.BatteryLevel = rData.BatteryLevel
	sData.TotalSpace = rData.TotalSpace
	sData.AvailableSpace = rData.AvailableSpace
	x, _ := json.Marshal(sData)
	player := core.RoomMgrObj.GetTPlayer(p.TID)
	if player != nil {
		player.SendMsg(2, 10025, x)
	}

	//数据
	fmt.Println("[2-20025]更新电池电量 rData = ", rData)

}

// 接收到老师端暂停视频的命令
func (aa *RoomApi) Handle_onVideoPause(p *core.Player, data []byte) {

	players := core.RoomMgrObj.GetAllPlayers(p.TID)
	if players != nil {
		for _, player := range players {
			player.SendMsg(2, 10026, data)
		}
	}

	//数据
	fmt.Println("[2-20026]接收到老师端暂停视频的命令 data = ", data)

}

// 服务器转发给老师端学生进度信息
func (aa *RoomApi) Handle_onCourseNodeUpdate(p *core.Player, data []byte) {

	rData := &pb.Sync_CourseNodeData{}
	json.Unmarshal(data, rData)

	var sData pb.Sync_CourseNodeData_Send
	sData.StuName = p.CID
	sData.NodeName = rData.NodeName
	sData.NodeIndex = rData.NodeIndex
	sData.NodeTotal = rData.NodeTotal
	x, _ := json.Marshal(sData)

	player := core.RoomMgrObj.GetTPlayer(p.TID)
	if player != nil {
		player.SendMsg(2, 10027, x)
	}

	//数据
	fmt.Println("[2-20027]服务器转发给老师端学生进度信息 rData = ", rData)

}

// 发送给学生端游戏暂停命令
func (aa *RoomApi) Handle_onGamePause(p *core.Player, data []byte) {

	players := core.RoomMgrObj.GetAllPlayers(p.TID)
	if players != nil {
		for _, player := range players {
			player.SendMsg(2, 10028, data)
		}
	}

	//数据
	fmt.Println("[2-20028]发送给学生端游戏暂停命令 data = ", data)

}

// 删除学生
func (aa *RoomApi) Handle_DeleteStu(p *core.Player, data []byte) {

	getData := &pb.StudentData{}
	json.Unmarshal(data, getData)

	db := global.SqliteInst.GetDB()
	stmt, err := db.Prepare("DELETE FROM tb_snum WHERE snum = ?")
	defer stmt.Close()
	if err == nil {
		stmt.Exec(getData.Snum)
	} else {
		log.Println("Mysql_AddStudent,", err)
	}

	//数据
	fmt.Println("[2-20029]老师端删除学生 student = ", getData)

}

// 添加学生
func (aa *RoomApi) Handle_AddUser(p *core.Player, data []byte) {

	var students []pb.StudentInfo
	json.Unmarshal(data, students)

	db := global.SqliteInst.GetDB()
	for _, s := range students {
		stmt, err := db.Prepare("insert into tb_user(userName, isTeacherClose) values(?,'1')")
		defer stmt.Close()
		if err == nil {
			stmt.Exec(s.StuUserName)
		} else {
			log.Println("Mysql_AddUser,", err)
		}
	}

	//数据
	fmt.Println("[2-20030]老师端添加学生 student = ", students)

}

// 收到老师端发来的学生信息
func (aa *RoomApi) Handle_onStudentData(p *core.Player, data []byte) {

	var students pb.Sync_EnterRoom
	json.Unmarshal(data, students)

	cRoom := core.RoomMgrObj.GetRoom(p.TID)
	if cRoom != nil {
		//cRoom.AllStudents = lData.Students

		db := global.SqliteInst.GetDB()
		players := cRoom.GetAllPlayers()

		for _, p := range players {
			if Mysql_HasUser(p.CID) {
				continue
			}
			stmt, err := db.Prepare("insert into tb_user(userName, isTeacherClose) values(?,'1')")
			defer stmt.Close()
			if err == nil {
				stmt.Exec(p.CID)
			} else {
				log.Println("Mysql_AddUser,", err)
			}
		}

	}

	//数据
	fmt.Println("[2-20031]收到老师端发来的学生信息 student = ", students)

}

func Mysql_HasUser(userName string) bool {
	var UserName string
	db := global.SqliteInst.GetDB()
	rows, err := db.Query("select userName from tb_user where userName = ?", userName)
	defer rows.Close()
	if err == nil {
		for rows.Next() {
			if err := rows.Scan(&UserName); err == nil {
				return true
			} else {
				return false
			}
		}
	}
	return false
}

// 开启护眼模式
func (aa *RoomApi) Handle_onClientHuyan(p *core.Player, data []byte) {

	players := core.RoomMgrObj.GetAllPlayers(p.TID)
	if players != nil {
		for _, player := range players {
			player.SendMsg(2, 10032, data)
		}
	}

	//数据
	fmt.Println("[2-20032]通知学生开启护眼 data = ", data)

}

// 关闭护眼模式
func (aa *RoomApi) Handle_onClientHuyanClose(p *core.Player, data []byte) {

	players := core.RoomMgrObj.GetAllPlayers(p.TID)
	if players != nil {
		for _, player := range players {
			player.SendMsg(2, 10033, data)
		}
	}

	//数据
	fmt.Println("[2-20033]通知学生关闭护眼 data = ", data)

}

// 开启录屏
func (aa *RoomApi) Handle_onClientLpStart(p *core.Player, data []byte) {

	msg := &pb.Sync_UpdataCourse{}
	json.Unmarshal(data, msg)

	players := core.RoomMgrObj.GetAllPlayers(p.TID)
	if players != nil {
		for _, player := range players {
			if player.CID == msg.StuAccountID {
				player.SendMsg(2, 10034, data)
			}
		}
	}

	//数据
	fmt.Println("[2-20034]通知学生开启录屏 data = ", data)

}

// 开启录屏
func (aa *RoomApi) Handle_onClientLpClose(p *core.Player, data []byte) {

	msg := &pb.Sync_UpdataCourse{}
	json.Unmarshal(data, msg)

	players := core.RoomMgrObj.GetAllPlayers(p.TID)
	if players != nil {
		for _, player := range players {
			if player.CID == msg.StuAccountID {
				player.SendMsg(2, 10035, data)
			}
		}
	}

	//数据
	fmt.Println("[2-20035]通知学生结束录屏 data = ", data)

}

// 开启/关闭游戏状态
func (aa *RoomApi) Handle_onTurnOnOrOffStatus(p *core.Player, data []byte) {

	msg := &pb.Sync_GameStatus{}
	json.Unmarshal(data, msg)
	call := &pb.Sync_GameStatusBack{Code: 1}
	fmt.Println("[2-20036] param = ", msg, " serverStatus = ", core.RoomMgrObj.GlobalStatus)
	if msg != nil {
		if msg.Type < 1 || msg.Status < 0 {
			call.Code = -1
			callData, _ := json.Marshal(call)
			fmt.Println("[2-20036] err code = -1")
			p.SendMsg(2, 10036, callData)
			return
		}
		code := core.RoomMgrObj.RefrushGameStatus(msg)
		fmt.Println("[2-20036] refrush code = ", code)
		//通知老师操作结果
		call.Code = code
		callData, _ := json.Marshal(call)
		p.SendMsg(2, 10036, callData)
		if code == 1 {
			//通知学生操作结果
			players := core.RoomMgrObj.GetAllPlayers(p.TID)
			if players != nil {
				for _, player := range players {
					fmt.Println("[2-20036] checkSend = ", player.CID, p.CID)
					if player.CID != p.CID {
						fmt.Println("[2-20036] sendmsg player = ", data)
						player.SendMsg(2, 10036, data)
					}
				}
			}
		}
	}
	//数据
	fmt.Println("[2-20036]更新学生全局状态 data = ", data)

}
