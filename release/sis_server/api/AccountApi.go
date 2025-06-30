package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lyyym/zinx-wsbase/global"
	"github.com/lyyym/zinx-wsbase/release/sis_server/core"
	"github.com/lyyym/zinx-wsbase/release/sis_server/pb"
	"github.com/lyyym/zinx-wsbase/ziface"
	"github.com/lyyym/zinx-wsbase/znet"
	"go.uber.org/zap"
	"log"
	"time"
)

type AccountApi struct {
	znet.BaseRouter
}

func (aa *AccountApi) Handle(request ziface.IRequest) {

	//1. 得到消息的Sub，用来细化业务实现
	sub := request.GetSubID()
	//sub := request.GetMsgSub()
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
	//fmt.Println("[Receive Account Msg] : Player = " , player.CID )

	switch sub {
	case 10013:
		fmt.Println("check alive!")
		break
	case 10012:
		aa.Handle_onRequest10002(player, request.GetData())
		break
	case 10002: //登录
		msg := &pb.SyncLogin{}
		json.Unmarshal(request.GetData(), msg)
		fmt.Println(" TID = ", msg.TID, " CID = ", msg.CID)

		player.Login(msg)

		//给老师回传 年级分类信息列表数据
		if player.TID == player.CID {
			aa.Handle_GetStus(player)
			aa.Handle_NjFl(player)
			aa.Handle_KcFl(player)
			aa.Handle_KcDatas(player) //课程数据
		}

		////这里限制用户不能重复登录
		//if core.WorldMgrObj.HasLogined(player.CID) {
		//
		//	data1,_ := json.Marshal(&pb.SyncLoginB{CtrlFlag:"0",Code:"error"})
		//	////发送数据给客户端
		//	player.SendMsg(1,10002, data1)
		//
		//	player.Conn.GetTCPConnection().Close()
		//
		//
		//}else {
		//	player.Login(msg)
		//	//给老师回传 年级分类信息列表数据
		//	if player.TID == player.CID {
		//		aa.Handle_NjFl(player)
		//		aa.Handle_KcFl(player)
		//		aa.Handle_KcDatas(player)
		//	}
		//}

		break
	case 10005: //本地课程数据
		//给老师回传 年级分类信息列表数据
		aa.Handle_KcFl(player)
		break
	case 10006: //本地课程数据
		//给老师回传 年级分类信息列表数据
		aa.Handle_KcDatas(player)
		break
	case 10009: //学生端心跳
		//给老师回传 年级分类信息列表数据
		player.BeatTime = int(time.Now().Unix())
		//fmt.Println("收到心跳 ： ", player.PID, player.CID, player.TID)
		player.SendMsg(1, 10009, request.GetData())
		break
	}

}

// 绑定
func (aa *AccountApi) Handle_onRequest10002(p *core.Player, data []byte) {
	msg := &pb.Tcp_Bind{}
	json.Unmarshal(data, msg)
	if msg == nil {
		fmt.Println("Tcp_Bind err")
		return
	}
	//1. 先绑定数据给玩家
	if p.Bind(msg.UserToken) {
		fmt.Println("bind ", p.CID)
		global.Glog.Info("Bind UserName = %s , Pid = %d",
			zap.String("UserName", p.CID),
			zap.Int32("PID", p.PID),
		)
		//2. 绑定玩家到世界
		core.WorldMgrObj.BindPlayer(p)

		//给老师回传 年级分类信息列表数据
		if p.TID == p.CID {
			aa.Handle_GetStus(p)
			aa.Handle_NjFl(p)
			aa.Handle_KcFl(p)
			aa.Handle_KcDatas(p) //课程数据
		}

	} else {
		global.Glog.Info("Bind Error")
	}
}

// 老师获取学生列表
func (aa *AccountApi) Handle_GetStus(p *core.Player) {

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
	var sData pb.AllStudentData
	sData.AllData = allStudent
	x, _ := json.Marshal(sData)
	p.SendMsg(1, 10008, x)
	fmt.Println("[1-10008]老师获取学生列表 student = ", sData)

}

// 年纪分类数据
func (aa *AccountApi) Handle_NjFl(p *core.Player) {

	//
	var allLocalGradeType pb.Sync_GetLocalGradeTypeList
	var localGradeTypeData pb.LocalGradeTypeData
	//从数据库获取数据
	//获取本地年级分类信息列表
	db := global.SqliteInst.GetDB()
	rows, err := db.Query("select ID,typeName, visible from tb_classtype")
	if err != nil {
		fmt.Println("Sqlite Test Query DB Err")
	} else {
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&localGradeTypeData.ID, &localGradeTypeData.TypeName, &localGradeTypeData.Visible); err == nil {
				allLocalGradeType.LocalGradeTypeArr = append(allLocalGradeType.LocalGradeTypeArr, localGradeTypeData)
			} else {
				log.Println("Mysql_GetLocalGradeTypeData,", err)
			}
		}
	}

	//序列化数据
	data, _ := json.Marshal(allLocalGradeType)
	fmt.Println("[10004]老师回执年纪分类: dataLength = ", len(data))
	//发送
	p.SendMsg(1, 10004, data)

}

// 获取本地课程分类信息列表
func (aa *AccountApi) Handle_KcFl(p *core.Player) {

	//

	var allLocalCourseType pb.Sync_GetLocalCourseTypeList
	var localCourseTypeData pb.LocalCourseTypeData
	//从数据库获取数据
	//获取本地年级分类信息列表
	db := global.SqliteInst.GetDB()
	rows, err := db.Query("select ID,typeName, bEdit,bNotClassified,inClassType,inClassTypeSort from tb_coursetype")
	if err != nil {
		fmt.Println("Sqlite Handle_KcFl Query DB Err")
	} else {
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&localCourseTypeData.ID, &localCourseTypeData.TypeName, &localCourseTypeData.BEdit, &localCourseTypeData.BNotClassified, &localCourseTypeData.InClassType, &localCourseTypeData.InClassTypeSort); err == nil {
				allLocalCourseType.LocalCourseTypeArr = append(allLocalCourseType.LocalCourseTypeArr, localCourseTypeData)
			} else {
				log.Println("Mysql_GetLocalCourseTypeData,", err)
			}
		}
	}

	//x1, _ := json.Marshal(rData)
	//jsonstr1 := string(x1)

	//序列化数据
	data, _ := json.Marshal(allLocalCourseType)

	//fmt.Println("老师回执课程分类信息列表 = " , allLocalCourseType,len(data))
	fmt.Println("[10005]老师回执课程分类信息列表: dataLength = ", len(data))
	//发送
	p.SendMsg(1, 10005, data)

}

// 获取本地课程信息列表
func (aa *AccountApi) Handle_KcDatas(p *core.Player) {

	//

	var allLocalCourse pb.Sync_GetLocalCourseList
	var localCourseData pb.CoursewareData
	//从数据库获取数据
	//获取本地年级分类信息列表
	var g_md5 sql.NullString
	var g_gameUrl sql.NullString
	var g_resVersion sql.NullString
	db := global.SqliteInst.GetDB()
	rows, err := db.Query("select ID,courseName, iconName,courseID,courseType," +
		"courseOwner,inCourseType,inCourseTypeSort,thirdType,ThirdMsg,md5,gameUrl,resVersion from tb_course")
	if err != nil {
		fmt.Println("Sqlite Handle_KcFl Query DB Err")
	} else {
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&localCourseData.ID, &localCourseData.Name, &localCourseData.IconName, &localCourseData.CourseID, &localCourseData.CourseType,
				&localCourseData.CourseOwner,
				&localCourseData.InCourseType, &localCourseData.InCourseTypeSort, &localCourseData.ThirdType, &localCourseData.ThirdMsg, &g_md5,
				&g_gameUrl, &g_resVersion); err == nil {

				localCourseData.Md5 = g_md5.String
				localCourseData.GameUrl = g_gameUrl.String
				localCourseData.ResVersion = g_resVersion.String

				allLocalCourse.LocalCourseArr = append(allLocalCourse.LocalCourseArr, localCourseData)
			} else {
				log.Println("Mysql_GetLocalCourseData ===== ,", err.Error())
			}
		}
	}

	fmt.Println("老师回执本地课程信息列表 = ", allLocalCourse)

	//序列化数据
	data, _ := json.Marshal(allLocalCourse)
	fmt.Println("[10006]老师回执本地课程信息列表: dataLength = ", len(data))
	//发送
	p.SendMsg(1, 10006, data)

}
