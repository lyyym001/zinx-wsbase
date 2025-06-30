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
)

type CourseApi struct {
	znet.BaseRouter
}

func (ca *CourseApi) Handle(request ziface.IRequest) {

	//1. 得到消息的Sub，用来细化业务实现
	sub := request.GetSubID()
	//fmt.Println("CourseApi Do : msgID = " , request.GetMsgID() , " Sub = " , request.GetMsgSub() , " msgLength = " , len(request.GetData()))

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
	fmt.Println("[Receive Course Msg] : Player = ", player.CID)

	switch sub {

	case 10001: //登录
		ca.Handle_onGetReportBySnum(player, request.GetData())
		break
	case 10002: //本地课程数据
		ca.Handle_onGetReportBySname(player, request.GetData())
		break
	case 10003: //本地课程数据
		ca.Handle_onGetReportByCourseName(player, request.GetData())
		break
	case 10004: //本地课程数据
		ca.Handle_onGetReportByCourseType(player, request.GetData())
		break
	case 10005: //本地课程数据
		ca.Handle_onOpentheprojection(player, request.GetData())
		break
	case 10006: //本地课程数据
		ca.Handle_onReceivePic(player, request.GetData())
		break
	case 10007: //本地课程数据
		ca.Handle_onClosetheprojection(player, request.GetData())
		break
	case 10008: //本地课程数据
		ca.Handle_onAddLocalCourseData(player, request.GetData())
		break
	case 10009: //本地课程数据
		ca.Handle_onCreateCourseTypeInfo(player, request.GetData())
		break
	case 10010: //本地课程数据
		ca.Handle_onUpdateCourseTypeInfo(player, request.GetData())
		break
	case 10011: //本地课程数据
		ca.Handle_onDeleteCourseTypeInfo(player, request.GetData())
		break
	case 10012: //本地课程数据
		ca.Handle_onUpdateClassifyTypeInfo2(player, request.GetData())
		break
	case 10013: //本地课程数据
		ca.Handle_onUpdateCourseInfo(player, request.GetData())
		break
	case 10014: //本地课程数据
		ca.Handle_onUpdateOneCourse(player, request.GetData())
		break
	case 10015: //本地课程数据
		ca.Handle_CR_OpenCreateXCourse(player, request.GetData())
		break
	case 10016: //本地课程数据
		ca.Handle_onCreatexCourseReset(player, request.GetData())
		break
	case 10017: //本地课程数据
		ca.Handle_SendScoreToTeacher(player, request.GetData())
		break
	case 10018:
		ca.Handle_onAddCleverData(player, request.GetData())
		break
	case 10019:
		ca.Handle_UpdateCleverCourse(player, request.GetData())
		break
	}

}

// 老师获取某个学生的考试记录
func (aa *CourseApi) Handle_onGetReportBySnum(p *core.Player, data []byte) {

	var sData pb.Sync_GetReportBySnumData_Send
	getData := &pb.Sync_GetReportBySnumData{}
	json.Unmarshal(data, getData)

	playerT := core.RoomMgrObj.GetTPlayer(p.TID)
	if playerT != nil {

		var studyCount int
		var courseName, courseType string
		var averageTime, personalAverageScore, personalHighestScore, personalLowestScore, totalAverageScore, totalHighestScore float32
		var studentReport pb.StudentReportData
		var allStudentReport []pb.StudentReportData
		db := global.SqliteInst.GetDB()
		rows, err := db.Query("select CourseName,CourseType,count(id),avg(StuTimeLong),avg(Score),max(Score),min(Score) from tb_recore where SNum = ? and Tid = ? group by CourseName", getData.Snum, p.TID)
		defer rows.Close()
		if err == nil {
			for rows.Next() {
				if err := rows.Scan(&courseName, &courseType, &studyCount, &averageTime, &personalAverageScore, &personalHighestScore, &personalLowestScore); err == nil {
					studentReport.CourseName = courseName
					studentReport.CourseType = courseType
					studentReport.StudyCount = studyCount
					studentReport.AverageTime = averageTime
					studentReport.PersonalAverageScore = personalAverageScore
					studentReport.PersonalHighestScore = personalHighestScore
					studentReport.PersonalLowestScore = personalLowestScore
					rows, err := db.Query("select avg(Score),max(Score) from tb_recore where Tid = ? and CourseName = ?", p.TID, courseName)
					defer rows.Close()
					if err == nil {
						for rows.Next() {
							if err := rows.Scan(&totalAverageScore, &totalHighestScore); err == nil {
								studentReport.TotalAverageScore = totalAverageScore
								studentReport.TotalHighestScore = totalHighestScore
							}
						}
					}
					allStudentReport = append(allStudentReport, studentReport)
				} else {
					log.Println("Mysql_GetRecordBySnum,", err)
				}
			}
		}
		//数据
		sData.Snum = getData.Snum
		sData.ReportBySnumData = allStudentReport
		x, _ := json.Marshal(sData)
		p.SendMsg(3, 10001, x)
	}

	fmt.Println("[3-30001]老师获取某个学生的考试记录 student = ", sData)

}

// 根据学生姓名获取学生报表
func (aa *CourseApi) Handle_onGetReportBySname(p *core.Player, data []byte) {

	var sData pb.Sync_GetReportBySnameData_Send
	getData := &pb.Sync_GetReportBySnameData{}
	json.Unmarshal(data, getData)

	playerT := core.RoomMgrObj.GetTPlayer(p.TID)
	if playerT != nil {
		db := global.SqliteInst.GetDB()

		snum := 0

		rows1, err1 := db.Query("select snum from tb_snum where Tid = ? and pname = ?", p.TID, getData.Sname)
		defer rows1.Close()
		if err1 == nil {
			for rows1.Next() {
				if err := rows1.Scan(&snum); err == nil {

				} else {
					log.Println("Handle_onGetReportBySname,", err)
				}
			}
		}

		var studyCount int
		var courseName, courseType string
		var averageTime, personalAverageScore, personalHighestScore, personalLowestScore, totalAverageScore, totalHighestScore float32
		var studentReport pb.StudentReportData
		var allStudentReport []pb.StudentReportData

		rows, err := db.Query("select CourseName,CourseType,count(id),avg(StuTimeLong),avg(Score),max(Score),min(Score) from tb_recore where SNum = ? and Tid = ? group by CourseName", snum, p.TID)
		defer rows.Close()
		if err == nil {
			for rows.Next() {
				if err := rows.Scan(&courseName, &courseType, &studyCount, &averageTime, &personalAverageScore, &personalHighestScore, &personalLowestScore); err == nil {
					studentReport.CourseName = courseName
					studentReport.CourseType = courseType
					studentReport.StudyCount = studyCount
					studentReport.AverageTime = averageTime
					studentReport.PersonalAverageScore = personalAverageScore
					studentReport.PersonalHighestScore = personalHighestScore
					studentReport.PersonalLowestScore = personalLowestScore
					rows, err := db.Query("select CourseName,avg(Score),max(Score) from tb_recore where Tid = ? and CourseName = ?", p.TID, courseName)
					defer rows.Close()
					if err == nil {
						for rows.Next() {
							if err := rows.Scan(&totalAverageScore, &totalHighestScore); err == nil {
								studentReport.TotalAverageScore = totalAverageScore
								studentReport.TotalHighestScore = totalHighestScore
							}
						}
					}
					allStudentReport = append(allStudentReport, studentReport)
				} else {
					log.Println("Mysql_GetRecordBySnum,", err)
				}
			}
		}
		//数据
		sData.ReportBySnumData = allStudentReport
		x, _ := json.Marshal(sData)
		p.SendMsg(3, 10002, x)
	}

	fmt.Println("[3-30002]根据学生姓名获取学生报表 student = ", sData)

}

// 通过名称获得课程报表
func (aa *CourseApi) Handle_onGetReportByCourseName(p *core.Player, data []byte) {

	var sData pb.Sync_GetReportByCnameData_Send
	getData := &pb.Sync_GetReportByCnameData{}
	json.Unmarshal(data, getData)

	playerT := core.RoomMgrObj.GetTPlayer(p.TID)
	if playerT != nil {
		db := global.SqliteInst.GetDB()

		var studyCount int
		var courseName string
		var averageTime, averageScore float32
		var courseReport pb.CourseReportData
		var allCourseReport []pb.CourseReportData
		var likename string = "%"
		likename += getData.Cname
		likename += "%"
		rows, err := db.Query("select CourseName,count(id),avg(StuTimeLong),avg(Score) from tb_recore where CourseName like ? and Tid = ? group by CourseName", likename, p.TID)
		defer rows.Close()
		if err == nil {
			for rows.Next() {
				if err := rows.Scan(&courseName, &studyCount, &averageTime, &averageScore); err == nil {
					courseReport.CourseName = courseName
					courseReport.StudyCount = studyCount
					courseReport.AverageTime = averageTime
					courseReport.AverageScore = averageScore
					allCourseReport = append(allCourseReport, courseReport)
				} else {
					log.Println("Mysql_GetReportByCourseName,", err)
				}
			}
		}

		//数据
		sData.ReportByCnameData = allCourseReport
		x, _ := json.Marshal(sData)
		p.SendMsg(3, 10003, x)
	}

	fmt.Println("[3-30003]通过名称获得课程报表 student = ", sData)

}

// 通过类型获得课程报表
func (aa *CourseApi) Handle_onGetReportByCourseType(p *core.Player, data []byte) {

	var sData pb.Sync_GetReportByCtypeData_Send
	getData := &pb.Sync_GetReportByCtypeData{}
	json.Unmarshal(data, getData)

	playerT := core.RoomMgrObj.GetTPlayer(p.TID)
	if playerT != nil {
		db := global.SqliteInst.GetDB()

		var studyCount int
		var courseName string
		var averageTime, averageScore float32
		var courseReport pb.CourseReportData
		var allCourseReport []pb.CourseReportData
		rows, err := db.Query("select CourseName,count(id),avg(StuTimeLong),avg(Score) from tb_recore where CourseType = ? and Tid = ? group by CourseName", getData.Ctype, p.TID)
		defer rows.Close()
		if err == nil {
			for rows.Next() {
				if err := rows.Scan(&courseName, &studyCount, &averageTime, &averageScore); err == nil {
					courseReport.CourseName = courseName
					courseReport.StudyCount = studyCount
					courseReport.AverageTime = averageTime
					courseReport.AverageScore = averageScore
					allCourseReport = append(allCourseReport, courseReport)
				} else {
					log.Println("Mysql_GetReportByCourseType,", err)
				}
			}
		}

		//数据
		sData.ReportByCtypeData = allCourseReport
		x, _ := json.Marshal(sData)
		p.SendMsg(3, 10004, x)
	}

	fmt.Println("[3-30004]通过类型获得课程报表 student = ", sData)

}

// 老师 通知一个学生打开视频投影
func (aa *CourseApi) Handle_onOpentheprojection(p *core.Player, data []byte) {

	getData := &pb.Sync_OpentheprojectionData{}
	json.Unmarshal(data, getData)

	players := core.RoomMgrObj.GetAllPlayers(p.TID)
	if players != nil {
		for _, player := range players {
			if player.CID == getData.StuAccountID {
				player.SendMsg(3, 10005, data)
			}
		}
	}

	fmt.Println("[3-30005]通过类型获得课程报表 student = ", getData)

}

// 收到学生图片，发送给老师
func (aa *CourseApi) Handle_onReceivePic(p *core.Player, data []byte) {

	playerT := core.RoomMgrObj.GetTPlayer(p.TID)
	if playerT != nil {
		playerT.SendMsg(3, 10006, data)
	}

	fmt.Println("[3-30005]收到学生图片，发送给老师 data = ", len(data))

}

// 关闭投影
func (aa *CourseApi) Handle_onClosetheprojection(p *core.Player, data []byte) {

	getData := &pb.Sync_ClosetheprojectionData{}
	json.Unmarshal(data, getData)

	players := core.RoomMgrObj.GetAllPlayers(p.TID)
	if players != nil {
		for _, player := range players {
			if player.CID == getData.StuAccountID {
				player.SendMsg(3, 10007, data)
			}
		}
	}

	fmt.Println("[3-30007]关闭投影 student = ", getData)

}

// 添加一个本地课程信息
func (aa *CourseApi) Handle_onAddLocalCourseData(p *core.Player, data []byte) {

	getData := &pb.CoursewareData{}
	json.Unmarshal(data, getData)

	db := global.SqliteInst.GetDB()
	stmt, err := db.Prepare("insert into tb_course(courseName,iconName,courseID,courseType,courseOwner,inCourseType,inCourseTypeSort,thirdType,md5,gameUrl,resVersion) values(?,?,?,?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	if err == nil {
		stmt.Exec(getData.Name, getData.IconName, getData.CourseID, getData.CourseType, getData.CourseOwner, getData.InCourseType, getData.InCourseTypeSort, getData.ThirdType, getData.Md5, getData.GameUrl, getData.ResVersion)
		//拿到插入的数据库ID
		/*	res, _ := stmt.Exec(data.TypeName)
			id, _ := res.LastInsertId()
			return id*/
	} else {
		log.Println("Mysql_AddLocalCourseData,", err)
		//	return 0
	}

	fmt.Println("[3-30008]添加一个本地课程信息 student = ", getData)

}

// 创建1个课程类型信息(目录编辑)
func (aa *CourseApi) Handle_onCreateCourseTypeInfo(p *core.Player, data []byte) {

	getData := &pb.LocalCourseTypeData{}
	json.Unmarshal(data, getData)
	var SendData pb.CreateCourseTypeData_Send

	db := global.SqliteInst.GetDB()
	stmt, err := db.Prepare("insert into tb_coursetype(typeName,inClassType,inClassTypeSort) values(?,?,?)")
	defer stmt.Close()
	if err == nil {
		res, _ := stmt.Exec(getData.TypeName, getData.InClassType, getData.InClassTypeSort)
		id, _ := res.LastInsertId()
		SendData.ID = id
	} else {
		log.Println("Mysql_CreateCourseTypeInfo,", err)

	}

	playerT := core.RoomMgrObj.GetTPlayer(p.TID)
	if playerT != nil {
		data1, _ := json.Marshal(SendData)
		playerT.SendMsg(3, 10009, data1)
	}

	fmt.Println("[3-30009]创建1个课程类型信息(目录编辑) student = ", getData)

}

// 更新某个课程类型信息(目录编辑)
func (aa *CourseApi) Handle_onUpdateCourseTypeInfo(p *core.Player, data []byte) {

	getData := &pb.LocalCourseTypeData{}
	json.Unmarshal(data, getData)

	db := global.SqliteInst.GetDB()
	stmt, err := db.Prepare("UPDATE tb_coursetype SET typeName=?, inClassTypeSort=? WHERE ID = ?")
	defer stmt.Close()
	if err == nil {
		stmt.Exec(getData.TypeName, getData.InClassTypeSort, getData.ID)

	} else {
		log.Println("Mysql_UpdateCourseTypeInfo,", err)

	}

	fmt.Println("[3-30010]更新某个课程类型信息(目录编辑) student = ", getData)

}

// 删除某个课程类型信息(目录编辑)
func (aa *CourseApi) Handle_onDeleteCourseTypeInfo(p *core.Player, data []byte) {

	getData := &pb.DeleteCourseTypeData{}
	json.Unmarshal(data, getData)

	db := global.SqliteInst.GetDB()
	stmt, err := db.Prepare("DELETE FROM tb_coursetype WHERE ID = ?")
	defer stmt.Close()
	if err == nil {
		stmt.Exec(getData.ID)
	} else {
		log.Println("Mysql_DeleteCourseTypeInfo,", err)
	}

	fmt.Println("[3-30011]删除某个课程类型信息(目录编辑) student = ", getData)

}

// 更新年级类型信息(更新类型显隐)
func (aa *CourseApi) Handle_onUpdateClassifyTypeInfo2(p *core.Player, data []byte) {

	getData := &pb.LocalGradeTypeData{}
	json.Unmarshal(data, getData)

	db := global.SqliteInst.GetDB()
	stmt, err := db.Prepare("UPDATE tb_classtype SET visible=? WHERE ID = ?")
	defer stmt.Close()
	if err == nil {
		stmt.Exec(getData.Visible, getData.ID)
	} else {
		log.Println("Mysql_UpdateClassifyTypeInfo,", err)
	}

	fmt.Println("[3-30012]更新年级类型信息(更新类型显隐) student = ", getData)

}

// 更新某个课程信息(目录编辑)
func (aa *CourseApi) Handle_onUpdateCourseInfo(p *core.Player, data []byte) {

	getData := &pb.CoursewareData{}
	json.Unmarshal(data, getData)

	db := global.SqliteInst.GetDB()
	stmt, err := db.Prepare("UPDATE tb_course SET inCourseType=?, inCourseTypeSort=? WHERE ID = ?")
	defer stmt.Close()
	if err == nil {
		stmt.Exec(getData.InCourseType, getData.InCourseTypeSort, getData.ID)
	} else {
		log.Println("Mysql_UpdateCourseInfo,", err)
	}

	fmt.Println("[3-30013]更新某个课程信息(目录编辑) student = ", getData)

}

// 更新某个课程信息(目录编辑)
func (aa *CourseApi) Handle_onUpdateOneCourse(p *core.Player, data []byte) {

	getData := &pb.Sync_UpdataCourse{}
	json.Unmarshal(data, getData)

	players := core.RoomMgrObj.GetAllPlayers(p.TID)
	if players != nil {
		for _, player := range players {
			if getData.StuAccountID == "" {
				player.SendMsg(3, 10014, data)
			} else if getData.StuAccountID == player.CID {
				player.SendMsg(3, 10014, data)
				break
			}
		}
	}

	fmt.Println("[3-30014]更新某个课程信息(目录编辑) student = ", getData)

}

// 老师端开启Create-x课件
func (aa *CourseApi) Handle_CR_OpenCreateXCourse(p *core.Player, data []byte) {

	getData := &pb.Sync_GetCreatexCourse{}
	json.Unmarshal(data, getData)

	db := global.SqliteInst.GetDB()
	stmt, err := db.Prepare("UPDATE tb_user SET isTeacherClose=?")
	defer stmt.Close()
	if err == nil {
		if _, err := stmt.Exec("0"); err == nil {

		}
	}

	var sendmsg pb.Sync_SetCreatexCourse
	sendmsg.Uid = getData.Uid
	sendmsg.Pid = getData.Pid
	data1, _ := json.Marshal(sendmsg)

	players := core.RoomMgrObj.GetAllPlayers(p.TID)
	if players != nil {
		for _, player := range players {
			player.SendMsg(3, 10015, data1)
		}
	}

	fmt.Println("[3-30015]老师端开启Create-x课件 student = ", getData)

}

// 重启Create-x课程
func (aa *CourseApi) Handle_onCreatexCourseReset(p *core.Player, data []byte) {

	getData := &pb.Sync_GetRestarCreatex{}
	json.Unmarshal(data, getData)
	var sendData pb.Sync_SetCreatexCourse

	sendData.Uid = getData.Uid
	sendData.Pid = getData.Pid
	data1, _ := json.Marshal(sendData)

	players := core.RoomMgrObj.GetAllPlayers(p.TID)
	if players != nil {
		for _, player := range players {
			if getData.StuAccountID == "" {
				player.SendMsg(3, 10016, data1)
			} else if getData.StuAccountID == player.CID {
				player.SendMsg(3, 10016, data1)
				break
			}
		}
	}

	fmt.Println("[3-30016]重启Create-x课程 student = ", getData)

}

// 重启Create-x课程
func (aa *CourseApi) Handle_SendScoreToTeacher(p *core.Player, data []byte) {

	getData := &pb.Sync_ScoreGet{}
	json.Unmarshal(data, getData)
	var sendData pb.Sync_ScoreSend
	sendData.AccountName = p.CID
	sendData.SNum = getData.SNum
	sendData.CourseID = getData.CourseID
	sendData.Score = getData.Score
	sendData.Ability = getData.Ability
	data1, _ := json.Marshal(sendData)

	p.CourseCore = getData.Score
	p.CourseAbility = getData.Ability

	playerT := core.RoomMgrObj.GetTPlayer(p.TID)
	if playerT != nil {
		playerT.SendMsg(3, 10017, data1)
	}
	fmt.Println("[3-30017]重启Create-x课程 分数 = ", len(data1))

}

// 添加灵创课程
func (aa *CourseApi) Handle_onAddCleverData(p *core.Player, data []byte) {

	getData := &pb.CleverData{}
	json.Unmarshal(data, getData)
	var sid string = ""
	call := &pb.Sync_CleverCall{Code: 0}

	db := global.SqliteInst.GetDB()
	rows1, err1 := db.Query("select courseID from tb_course where courseID = ? ", getData.CourseID)
	defer rows1.Close()
	if err1 == nil {
		for rows1.Next() {
			if err := rows1.Scan(&sid); err == nil {

			} else {
				log.Println("Handle_onGetReportBySname,", err)
			}
		}
	}
	fmt.Println("[3-30018]添加灵创课程 sid = ,", sid)
	if len(sid) > 0 {
		fmt.Println("[3-30018]添加灵创课程 err = 作品已经存在,", getData.CourseID)
		data, _ := json.Marshal(call)
		p.SendMsg(3, 10018, data)
		return //作品已经存在
	}

	stmt, err := db.Prepare("insert into tb_course(courseName,courseID,courseType,inCourseType,InCourseTypeSort,md5,resVersion,gameUrl) values(?,?,?,?,?,?,?,?)")
	if err == nil {
		stmt.Exec(getData.CourseName, getData.CourseID, 999, "1,2,3,4", "56,112,63,10", "", 1, "")
		//拿到插入的数据库ID
		res, _ := stmt.Exec(getData.CourseID)
		_id, _ := res.LastInsertId()
		call.Code = 1
		call.Id = _id
		fmt.Println("[3-30018]添加灵创课程 err = 作品已添加成功 ， ID = ", _id)
		data, _ := json.Marshal(call)
		p.SendMsg(3, 10018, data)
	} else {
		log.Println("Mysql_AddLocalCourseData,", err)
		//	return 0
	}

	fmt.Println("[3-30018]添加灵创课程 回执 = ", call)

}

// 更新灵创课程
func (aa *CourseApi) Handle_UpdateCleverCourse(p *core.Player, data []byte) {

	CourseData := &pb.CleverData_Update{}
	json.Unmarshal(data, CourseData)
	call := &pb.Sync_CleverCall{Code: 0}
	call.Id = CourseData.Id
	fmt.Println("[3-30019]更新灵创课程 请求参数 = ", CourseData)
	// (2) 更新
	db := global.SqliteInst.GetDB()
	stmt, err := db.Prepare("UPDATE tb_course SET courseID=?,courseName=? WHERE ID = ?")
	if err != nil {
		fmt.Println("Sqlite Handle_UpdateCourse Update DB Err : ", err.Error())
		data, _ := json.Marshal(call)
		p.SendMsg(3, 10019, data)
		return
	} else {
		defer stmt.Close()
		_, err := stmt.Exec(CourseData.CourseID, CourseData.CourseName, CourseData.Id)
		//affectNum, err := result.RowsAffected()
		if err != nil {
			fmt.Println("Sqlite Handle_UpdateCourse DB Err 1")
			data, _ := json.Marshal(call)
			p.SendMsg(3, 10019, data)
			return
		}
		//fmt.Println("update affect rows is ", affectNum)
	}
	call.Code = 1
	data2, _ := json.Marshal(call)
	p.SendMsg(3, 10019, data2)
	fmt.Println("[3-30019]更新灵创课程 回执 = ", call)
}
