package api

import (
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/core"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/internal/models_sqlite"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/pb"
	"google.golang.org/protobuf/proto"

	"fmt"
	//"github.com/fd/lframework/utils"
	"github.com/lyyym/zinx-wsbase/ziface"
	"github.com/lyyym/zinx-wsbase/znet"
	//"log"
)

type WorkApi struct {
	znet.BaseRouter
}

func (aa *WorkApi) Handle(request ziface.IRequest) {

	//1. 得到消息的Sub，用来细化业务实现
	sub := request.GetSubID()
	//fmt.Println("WorkApi Do : msgID = " , request.GetMsgID() , " Sub = " , request.GetMsgSub() , " msgLength = " , len(request.GetData()))

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

	//fmt.Println("[Receive WorkApi Msg] : Player = " , player.PID )

	switch sub {

	case 10000:
		//老师/学生播放互动内容 ？
		aa.Handle_onRequestCoursePlay(player, request.GetData())
		break
	case 10001: //资源注册同步给其他客户端(ok)
		aa.Handle_onRequest10001(player, request.GetData())
		break
	case 10002: //资源被拿起，同步给其他客户端(ok)
		aa.Handle_onRequest10002(player, request.GetData())
		break
	case 10003: //资源脱离，同步给其他客户端(ok)
		aa.Handle_onRequest10003(player, request.GetData())
		break
	case 10004: //资源位置同步给其他人(ok)
		aa.Handle_onRequest10004(player, request.GetData())
		break
	case 10005: //同步给其他人资源状态(ok)
		aa.Handle_onRequest10005(player, request.GetData())
		break
	case 10006: //资源被其他人抓取了，取消本次抓取
		//从10002接口判断返回
		break
	case 10007: //关键步骤得分(ok)
		aa.Handle_onRequest10007(player, request.GetData())
		break
	case 10009: //资源准备完毕
		//aa.Handle_onRequest10009(player, request.GetData())
		break
	case 10010: //老师控制离开(离开后无法继续进入)   ?
		aa.Handle_onRequest10010(player, request.GetData())
		break
	case 10011: //成员主动离开(离开后无法继续进入)(ok)
		aa.Handle_onRequest10011(player, request.GetData())
		break
	case 10012: //toa消息转发(ok)
		aa.Handle_onRequest10012(player, request.GetData())
		break
	case 10013: //too消息转发(ok)
		aa.Handle_onRequest10013(player, request.GetData())
		break
	case 10014: //作品播放结束(ok)
		aa.Handle_onRequest10014(player, request.GetData())
		break
	case 10016: //提交答案(ok)
		aa.Handle_onRequest10016(player, request.GetData())
		break
	case 10017: //老师注册多人交互课件(废弃)
		//aa.Handle_onRequest10017(player, request.GetData())
		break
	case 10018: //老师开启一场交互课件(ok)
		aa.Handle_onRequest10018(player, request.GetData())
		break
	case 10019: //提交日志记录
		//aa.Handle_onRequest10019(player, request.GetData())
		break
	case 10020: //作品总统计(分页)
		//aa.Handle_onRequest10020(player, request.GetData())
		break
	case 10021: //作品分统计(分页)
		//aa.Handle_onRequest10021(player, request.GetData())
		break
	case 10022: //作品详情统计
		//aa.Handle_onRequest10022(player, request.GetData())
		break
	}

}

//// /作品详情统计
//func (aa *WorkApi) Handle_onRequest10022(p *core.Player, data []byte) {
//
//	request_data := &pb.Tcp_Tj3{}
//	proto.Unmarshal(data, request_data)
//	response_data := &pb.Tcp_Tj3Data{}
//	response_dataInfo := &pb.Tcp_Tj3Info{}
//
//	var sql string = ""
//	tableName := fmt.Sprintf("%s%d", "tb_work_", request_data.UniqueId)
//	if len(request_data.UName) == 0 {
//		sql = fmt.Sprintf("select uname,`date`,`type`,state,score,content from %s;", tableName)
//	} else {
//		sql = fmt.Sprintf("select uname,`date`,`type`,state,score,content from %s where uname = '%s';", tableName, request_data.UName)
//	}
//	fmt.Println("sql", sql)
//	db := utils.GlobalObject.SqliteInst.GetDB()
//	rows, _ := db.Query(sql)
//	defer rows.Close()
//	for rows.Next() {
//		if err := rows.Scan(&response_dataInfo.UName, &response_dataInfo.Date, &response_dataInfo.Type, &response_dataInfo.State, &response_dataInfo.Score, &response_dataInfo.Content); err == nil {
//			response_data.Data = append(response_data.Data, response_dataInfo)
//			log.Println("[统计]总记录(读取)-Succ,UniqueId=", request_data.UniqueId, ",UName=", request_data.UName)
//		} else {
//			log.Println("[统计]总记录(读取),UniqueId=", request_data.UniqueId, ",UName=", request_data.UName, ",Err=", err)
//		}
//	}
//
//	response_data.WorkName = request_data.WorkName
//	response_data.Mode = request_data.Mode
//	response_data.MaxScore = request_data.MaxScore
//	response_data.Number = request_data.Number
//	//data1, _ := json.Marshal(response_data)
//	//fmt.Println(response_data)
//	//p.SendMsg(2, 10022, response_data)
//
//}

//// /作品分统计(分页)
//func (aa *WorkApi) Handle_onRequest10021(p *core.Player, data []byte) {
//
//	request_data := &pb.Tcp_Tj2{}
//	json.Unmarshal(data, request_data)
//	response_data := &pb.Tcp_Tj2Data{}
//	response_dataInfo := &pb.Tcp_Tj2Info{}
//
//	var _maxCount int = 0
//
//	tableName := fmt.Sprintf("%s%d", "tb_work_", request_data.UniqueId)
//	sql := fmt.Sprintf("select uname,sum(score) from %s group by (uname) limit %d,9;", tableName, (request_data.Page-1)*9)
//	db := utils.GlobalObject.SqliteInst.GetDB()
//	rows, _ := db.Query(sql)
//	defer rows.Close()
//	for rows.Next() {
//		if err := rows.Scan(&response_dataInfo.UName, &response_dataInfo.Score); err == nil {
//			response_data.Data = append(response_data.Data, response_dataInfo)
//			_maxCount = _maxCount + 1
//			log.Println("[统计]总记录(读取)-Succ,UniqueId=", request_data.UniqueId, ",page=", request_data.Page)
//		} else {
//			log.Println("[统计]总记录(读取),UniqueId=", request_data.UniqueId, ",page=", request_data.Page, ",Err=", err)
//		}
//	}
//
//	//读取一下最大页数
//	var _maxPage int = 0
//	if _maxCount%9 == 0 {
//		_maxPage = _maxCount / 9
//	} else {
//		_maxPage = _maxCount/9 + 1
//	}
//	response_data.MaxNumber = int32(_maxCount)
//	response_data.MaxPage = int32(_maxPage)
//	//data1, _ := json.Marshal(response_data)
//	//fmt.Println(response_data)
//	//p.SendMsg(2, 10021, data1)
//
//}

//// /作品总统计(分页)
//func (aa *WorkApi) Handle_onRequest10020(p *core.Player, data []byte) {
//
//	request_data := &pb.Tcp_Tj1{}
//	json.Unmarshal(data, request_data)
//	response_data := &pb.Tcp_Tj1Data{}
//	response_dataInfo := &pb.Tcp_Tj1Info{}
//	//条件
//	var _condition string = "2,3"
//	if request_data.Mode2 == 1 && request_data.Mode3 != 1 {
//		_condition = "2"
//	} else if request_data.Mode2 != 1 && request_data.Mode3 == 1 {
//		_condition = "3"
//	}
//	sql := fmt.Sprintf("select `date`,mode,partnumber,maxscore,score,uniqueid from tb_work where workid = %d and mode in (%s) order by `date` desc limit %d,7;", request_data.Cid, _condition, (request_data.Page-1)*7)
//	db := utils.GlobalObject.SqliteInst.GetDB()
//	rows, _ := db.Query(sql)
//	defer rows.Close()
//	for rows.Next() {
//		if err := rows.Scan(&response_dataInfo.Date, &response_dataInfo.Mode, &response_dataInfo.Number, &response_dataInfo.MaxScore, &response_dataInfo.Score, &response_dataInfo.UniqueId); err == nil {
//			response_data.Data = append(response_data.Data, response_dataInfo)
//			log.Println("[统计]总记录(读取)-Succ,cid=", request_data.Cid, ",page=", request_data.Page)
//		} else {
//			log.Println("[统计]总记录(读取),cid=", request_data.Cid, ",page=", request_data.Page, ",Err=", err)
//		}
//	}
//
//	//读取一下最大页数
//	var _maxCount int = 0
//	sql1 := fmt.Sprintf("select count(id) from tb_work where workid = %d and mode in (%s);", request_data.Cid, _condition)
//	rows1, _ := db.Query(sql1)
//	defer rows1.Close()
//	for rows1.Next() {
//		if err := rows1.Scan(&_maxCount); err == nil {
//
//			log.Println("[统计]总记录(最大页)-Succ,cid=", request_data.Cid, ",page=", request_data.Page)
//		} else {
//			log.Println("[统计]总记录(最大页),cid=", request_data.Cid, ",page=", request_data.Page, ",Err=", err)
//		}
//	}
//
//	var _maxPage int = 0
//	if _maxCount%7 == 0 {
//		_maxPage = _maxCount / 7
//	} else {
//		_maxPage = _maxCount/7 + 1
//	}
//	response_data.MaxPage = int32(_maxPage)
//	//data1, _ := json.Marshal(response_data)
//	//fmt.Println(response_data)
//	//p.SendMsg(2, 10020, data1)
//
//}

//// /提交日志记录
//func (aa *WorkApi) Handle_onRequest10019(p *core.Player, data []byte) {
//
//	request_data := &pb.Tcp_WorkInfoRecord{}
//	json.Unmarshal(data, request_data)
//	var response_data pb.Tcp_ResponseRecordInfo
//
//	//操作db
//	var name sql.NullString
//	//var maxScore int  = 0
//
//	tableName := fmt.Sprintf("%s%d", "tb_work_", request_data.UniqueId)
//	sql := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s'", tableName)
//	db := utils.GlobalObject.SqliteInst.GetDB()
//	rows, _ := db.Query(sql)
//	defer rows.Close()
//	for rows.Next() {
//		if err := rows.Scan(&name); err == nil {
//			response_data.Code = 1
//			log.Println("[统计]提交日志一场交互课件(记录表是否存在)-Succ,UniqueId=", request_data.UniqueId, ",记录ID=", request_data.SetId, ",时间=", request_data.Date)
//		} else {
//			log.Println("[统计]提交日志一场交互课件(记录表是否存在),UniqueId=", request_data.UniqueId, ",记录ID=", request_data.SetId, ",时间=", request_data.Date, ",Err=", err)
//			response_data.Code = 0
//		}
//	}
//
//	if response_data.Code == 1 {
//		if request_data.State == 1 && request_data.Score > 0 {
//			stmt, err := db.Prepare("update tb_work set score = score + ? where uniqueid = ?;")
//			defer stmt.Close()
//			if err == nil {
//				stmt.Exec(request_data.Score, request_data.UniqueId)
//				response_data.Code = 1
//				log.Println("[统计]提交日志一场交互课件(加分)-Succ")
//			} else {
//				log.Println("[统计]提交日志一场交互课件(加分)-err=", err)
//				response_data.Code = 0
//			}
//		}
//	}
//
//	if name.Valid != false {
//		//插入
//		sql := fmt.Sprintf("insert into '%s' (username, uname,date,`type`,state,score,content) VALUES (?,?,?,?,?,?,?);", tableName)
//		stmt, err := db.Prepare(sql)
//		defer stmt.Close()
//		if err == nil {
//			stmt.Exec(request_data.Username, request_data.Uname, request_data.Date, request_data.Type, request_data.State, request_data.Score, request_data.Content)
//			response_data.Code = 1
//			log.Println("[统计]提交日志一场交互课件(写入)-Succ")
//		} else {
//			log.Println("[统计]提交日志一场交互课件(写入)-err=", err)
//			response_data.Code = 0
//		}
//	}
//
//	//响应
//	//fmt.Println("[统计]注册课件，UniqueId = ",UniqueId , " 状态：",response_data.Code,"(0-失败，1-成功)")
//	//data1, _ := json.Marshal(response_data)
//	//p.SendMsg(2,10019,data1)
//
//}

// /老师开启一场交互课件
func (aa *WorkApi) Handle_onRequest10018(p *core.Player, data []byte) {

	request_data := &pb.Tcp_WorkRecord{}
	proto.Unmarshal(data, request_data)
	response_data := &pb.Tcp_ResponseRecord{}

	// 1.最大
	var newDid uint
	mode := &models_sqlite.WorkBasic{}
	err := models_sqlite.DB.Last(mode).Error
	if err != nil {
		newDid = 1
	} else {
		newDid = mode.ID + 1
	}
	var UniqueId int = int(newDid)
	// 2.插入记录
	u := &models_sqlite.WorkBasic{
		WorkId:     request_data.Workid,
		Workname:   request_data.Workname,
		Date:       request_data.Date,
		Mode:       int(request_data.Mode),
		Partnumber: int(request_data.Partnumber),
		MaxScore:   int(request_data.MaxScore),
		Uniqueid:   UniqueId,
	}
	if err := models_sqlite.DB.Create(&u).Error; err != nil {
		fmt.Println("insert dir error")
		response_data.Code = 0
	} else {
		response_data.Code = 1
		//3. 创建记录表
		record := &models_sqlite.WorkRecordBasic{}
		models_sqlite.DB.Table(record.TableName(UniqueId)).AutoMigrate(record)
		response_data.UniqueId = int32(UniqueId)
	}

	core.WorldMgrObj.Toa_NoGz_InScene(2, 10018, response_data)

	//if response_data.Code == 1 {
	//	//创建Table
	//	tableName := fmt.Sprintf("%s%d", "tb_work_", UniqueId)
	//	//fmt.Println("UniqueId = ",UniqueId)
	//	//fmt.Println(tableName)
	//	var sql string = fmt.Sprintf("CREATE TABLE %s ( id INTEGER PRIMARY KEY AUTOINCREMENT, username STRING NOT NULL,uname STRING NOT NULL,date STRING NOT NULL,type INT NOT NULL,state INT NOT NULL,score FLOAT NOT NULL,content STRING NOT NULL);", tableName)
	//	stmt, err := db.Prepare(sql)
	//	defer stmt.Close()
	//	if err == nil {
	//		stmt.Exec()
	//		response_data.Code = 1
	//		log.Println("[统计]老师开启一场交互课件(创建分表)-Succ")
	//	} else {
	//		log.Println("[统计]老师开启一场交互课件(创建分表)-err=", err)
	//		response_data.Code = 0
	//	}
	//}
	//
	////响应
	//if response_data.Code == 1 {
	//	response_data.UniqueId = UniqueId
	//}
	//fmt.Println("[统计]注册课件，UniqueId = ", UniqueId, " 状态：", response_data.Code, "(0-失败，1-成功)")
	////data1, _ := json.Marshal(response_data)
	////core.WorldMgrObj.WorldToa(2, 10018, data1)

}

//// /老师注册多人交互课件
//func (aa *WorkApi) Handle_onRequest10017(p *core.Player, data []byte) {
//
//	request_data := &pb.Tcp_Shixun{}
//	proto.Unmarshal(data, request_data)
//	var response_data pb.Tcp_ResponseScene
//
//	core.WorldMgrObj.MScene.Running = true
//	core.WorldMgrObj.MScene.Mode = byte(request_data.Mode)
//	core.WorldMgrObj.MScene.SceneID = request_data.WorkId
//	core.WorldMgrObj.MScene.Players = make(map[string]int32)
//	//core.WorldMgrObj.PullListToCell(request_data.Users)
//	core.WorldMgrObj.MScene.TakeObjects = make(map[int]int32)
//	core.WorldMgrObj.MScene.Questions = make(map[int]int32)
//	core.WorldMgrObj.MScene.JHObjects = make(map[int]*core.JHObject)
//	core.WorldMgrObj.MScene.Steps = make(map[int]*core.Step)
//	core.WorldMgrObj.MScene.MainCtroller = "0" //默认老师作为主控
//
//	//响应
//	response_data.UNumber = int32(len(core.WorldMgrObj.MScene.Players)) //参与人数
//	response_data.MainCtroller = core.WorldMgrObj.MScene.MainCtroller   //主控用户
//	response_data.WorkId = request_data.WorkId
//	response_data.Mode = request_data.Mode //当前播放模式
//
//	fmt.Println("[老师端注册课件成功]请求成员组=", request_data, " 响应成员组=", core.WorldMgrObj.MScene.Players, " pid=", p.PID)
//	//data1, _ := json.Marshal(response_data)
//	//core.WorldMgrObj.Too(p.PID, 2, 10017, data1)
//	core.WorldMgrObj.PrintUserList()
//
//}

// /提交答案
func (aa *WorkApi) Handle_onRequest10016(p *core.Player, data []byte) {

	request_data := &pb.Tcp_Step{}
	proto.Unmarshal(data, request_data)
	response_data := &pb.Tcp_QuestionInfo{}

	fmt.Println("[提交答案]", request_data)
	if core.WorldMgrObj.MScene != nil && core.WorldMgrObj.MScene.Questions != nil {

		if _, ok := core.WorldMgrObj.MScene.Questions[int(request_data.StepId)]; ok {
			//已经被别人答了
			response_data.Code = 0
			//data1, _ := json.Marshal(response_data)
			p.SendMsg(2, 10016, response_data)
		} else {
			response_data.Code = 1
			response_data.Qid = request_data.StepId
			response_data.StepDate = request_data.StepDate
			response_data.StepState = request_data.StepState
			response_data.UName = request_data.UName
			//data1, _ := json.Marshal(response_data)
			core.WorldMgrObj.Toa_NoGz_InScene(2, 10016, response_data)
		}
	}
}

// /toa消息转发
func (aa *WorkApi) Handle_onRequest10014(p *core.Player, data []byte) {

	core.WorldMgrObj.Too(p.UserName, 2, 10014, nil)

}

// /toa消息转发
func (aa *WorkApi) Handle_onRequest10012(p *core.Player, data []byte) {
	request_data := &pb.Tcp_To{}
	proto.Unmarshal(data, request_data)
	core.WorldMgrObj.Toa_NoGz_InScene(2, 10012, request_data)

}

// /too消息转发
func (aa *WorkApi) Handle_onRequest10013(p *core.Player, data []byte) {
	request_data := &pb.Tcp_To{}
	proto.Unmarshal(data, request_data)
	core.WorldMgrObj.Too(p.UserName, 2, 10013, request_data)

}

// /成员主动离开(离开后无法继续进入)
func (aa *WorkApi) Handle_onRequest10011(p *core.Player, data []byte) {
	request_data := &pb.Tcp_Leave{}
	proto.Unmarshal(data, request_data)
	core.WorldMgrObj.WorkLeave(p.UserName, int(request_data.State))

}

// /老师控制离开(离开后无法继续进入)
func (aa *WorkApi) Handle_onRequest10010(p *core.Player, data []byte) {

	//通知当前所有用户结束
	//core.WorldMgrObj.Too(p.PID, 2, 10010, []byte("ok"))

	if core.WorldMgrObj.MScene != nil {
		core.WorldMgrObj.MScene.Players = make(map[string]int32)
	}

	core.WorldMgrObj.WorkFinish(1)
}

//// /资源准备完毕
//func (aa *WorkApi) Handle_onRequest10009(p *core.Player, data []byte) {
//
//	core.WorldMgrObj.UpdateReadedStatus(p.PID)
//	if core.WorldMgrObj.AllReaded() {
//		fmt.Println("[成员全部准备完毕]")
//		core.WorldMgrObj.PrintUserList()
//		//core.WorldMgrObj.Toa(2, 10009, []byte("ok"))
//	}
//}

// /关键步骤得分
func (aa *WorkApi) Handle_onRequest10007(p *core.Player, data []byte) {

	request_data := &pb.Tcp_Step{}
	proto.Unmarshal(data, request_data)
	core.WorldMgrObj.RegisterStep(request_data)
	//给参与作品的其他人转发消息
	core.WorldMgrObj.Too(p.UserName, 2, 10007, request_data)
}

// /同步给其他人资源状态
func (aa *WorkApi) Handle_onRequest10005(p *core.Player, data []byte) {

	request_data := &pb.Tcp_ObjectStatus{}
	proto.Unmarshal(data, request_data)
	core.WorldMgrObj.UpdateObjectStatus(request_data)
	//给参与作品的其他人转发消息
	core.WorldMgrObj.Too(p.UserName, 2, 10005, request_data)
}

// /资源位置同步给其他人
func (aa *WorkApi) Handle_onRequest10004(p *core.Player, data []byte) {

	request_data := &pb.TCP_TbObj{}
	proto.Unmarshal(data, request_data)
	core.WorldMgrObj.UpdateObjectPos(request_data)
	//给参与作品的其他人转发消息
	core.WorldMgrObj.Too(p.UserName, 2, 10004, request_data)
}

// /资源脱离，同步给其他客户端
func (aa *WorkApi) Handle_onRequest10003(p *core.Player, data []byte) {

	request_data := &pb.Tcp_Object{}
	proto.Unmarshal(data, request_data)
	if _, ok := core.WorldMgrObj.MScene.TakeObjects[int(request_data.ObjId)]; ok {
		//资源脱离
		delete(core.WorldMgrObj.MScene.TakeObjects, int(request_data.ObjId))
		fmt.Println("[资源脱离成功],id=", request_data.ObjId, "记录的脱离后数据=", core.WorldMgrObj.MScene.TakeObjects)
	} else {
		fmt.Println("[资源脱离失败],id=", request_data.ObjId, "未找到服务器记录，记录的数据=", core.WorldMgrObj.MScene.TakeObjects)
	}
	//给参与作品的其他人转发消息
	core.WorldMgrObj.Too(p.UserName, 2, 10003, request_data)
}

// /资源被拿起，同步给其他客户端
func (aa *WorkApi) Handle_onRequest10002(p *core.Player, data []byte) {

	request_data := &pb.TCP_Grab_Call{}
	proto.Unmarshal(data, request_data)
	if _, ok := core.WorldMgrObj.MScene.TakeObjects[int(request_data.ObjId)]; ok {
		//资源已经被拿取了
		if p != nil {
			p.SendMsg(2, 10006, request_data)
			fmt.Println("[资源拿起失败],id=", request_data.ObjId, "记录的拿起数据=", core.WorldMgrObj.MScene.TakeObjects)
		}
	} else {
		core.WorldMgrObj.MScene.TakeObjects[int(request_data.ObjId)] = p.PID
		//给参与作品的其他人转发消息
		core.WorldMgrObj.Too(p.UserName, 2, 10002, request_data)
		fmt.Println("[资源拿起成功],id=", request_data.ObjId, "记录的拿起数据=", core.WorldMgrObj.MScene.TakeObjects)
	}

}

// /资源注册同步给其他客户端
func (aa *WorkApi) Handle_onRequest10001(p *core.Player, data []byte) {

	request_data := &pb.TCP_RegisterObj{}
	proto.Unmarshal(data, request_data)
	core.WorldMgrObj.RegisterObject(request_data)
	fmt.Println("给其他人同步注册的资源：", request_data)
	//给参与作品的其他人转发消息
	core.WorldMgrObj.Too(p.UserName, 2, 10001, request_data)
}

// /老师/学生 播放互动内容
func (aa *WorkApi) Handle_onRequestCoursePlay(p *core.Player, data []byte) {

	//var response_data pb.Tcp_ResponseScene
	//request_data := &pb.Tcp_RequestScene{}
	////request_data.mode = 1-学练模式 2-考评模式 3-协同模式
	//err := proto.Unmarshal(data, request_data)
	//if err != nil {
	//	fmt.Println("proto.Unmarshal err", err)
	//	return
	//}
	//fmt.Println("[作品播控鉴权]作品ID=", request_data.SceneId, "模式=", request_data.Mode)
	//
	////学生不鉴权-单机播放
	//if p.AccountType == 0 {
	//	response_data.Code = 0
	//} else {
	//	if core.WorldMgrObj.MScene != nil && core.WorldMgrObj.MScene.Running {
	//		response_data.Code = -4 //作品正在进行，需要先结束作品
	//	} else {
	//		//fmt.Println("交互内容配置，scene.conf=", gameutils.GlobalScene.S)
	//		if _, ok := gameutils.GlobalScene.S[request_data.SceneId]; ok {
	//			//1.参数验证-模式选择是否正确(1-学练模式 2-考评模式 3-协同模式)
	//			//学练模式-不限制
	//			//考核模式-需验证是否支持
	//			//协同模式-需验证是否支持
	//			if request_data.Mode == 2 && gameutils.GlobalScene.S[request_data.SceneId].Assessment != 1 {
	//				response_data.Code = -1 //模式异常，不支持考核模式
	//			} else if request_data.Mode == 3 && gameutils.GlobalScene.S[request_data.SceneId].Cooperate != 1 {
	//				response_data.Code = -2 //模式异常，不支持协同模式
	//			} else {
	//				//模式正确，开启新的作品
	//
	//				//2.记录全局数据
	//				core.WorldMgrObj.MScene.Running = true
	//				core.WorldMgrObj.MScene.Mode = byte(request_data.Mode)
	//				core.WorldMgrObj.MScene.SceneID = request_data.SceneId
	//				core.WorldMgrObj.MScene.Players = make(map[string]int32)
	//				core.WorldMgrObj.MScene.Questions = make(map[int]int32)
	//				core.WorldMgrObj.PullUserToCell()
	//				core.WorldMgrObj.MScene.TakeObjects = make(map[int]int32)
	//				core.WorldMgrObj.MScene.JHObjects = make(map[int]*core.JHObject)
	//				core.WorldMgrObj.MScene.Steps = make(map[int]*core.Step)
	//				if request_data.Mode == 3 {
	//					//协作模式，多人
	//					core.WorldMgrObj.MScene.MainCtroller = core.WorldMgrObj.TUserName   //默认老师作为主控
	//					response_data.UNumber = int32(len(core.WorldMgrObj.MScene.Players)) //参与人数
	//
	//				} else {
	//					core.WorldMgrObj.MScene.MainCtroller = ""
	//					response_data.UNumber = 1
	//				}
	//
	//				//core.WorldMgrObj.MScene.GlobalStepId = 0
	//				//fmt.Println("MSceneData = ",core.WorldMgrObj.MScene)
	//				//3.响应数据补充
	//				response_data.Code = 1
	//				response_data.MainCtroller = core.WorldMgrObj.MScene.MainCtroller //主控用户
	//				response_data.WorkId = request_data.SceneId
	//				response_data.Mode = request_data.Mode //当前播放模式
	//
	//				//data2,_ :=json.Marshal(gameutils.GlobalScene.ComputePos(request_data.SceneId))
	//				//response_data.Pos = string(data2)
	//			}
	//		} else {
	//			response_data.Code = -3 //鉴权失败，作品未授权控制
	//		}
	//	}
	//}
	//
	//if response_data.Code != 1 {
	//	//回执
	//	//data1, _ := json.Marshal(response_data)
	//	//fmt.Println("[鉴权失败]code=", response_data.Code)
	//	//p.SendMsg(2, 10000, data1)
	//} else {
	//	//广播
	//	fmt.Println("[鉴权成功]")
	//	//data1, _ := json.Marshal(response_data)
	//	//core.WorldMgrObj.Toa(2, 10000, data1)
	//	//core.WorldMgrObj.PrintUserList()
	//}
}
