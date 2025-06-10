package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/core"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/internal/models_sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func WorkList(c *gin.Context) {

	request_data := new(Tcp_Tj1)
	err := c.ShouldBindJSON(request_data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常,err = " + err.Error(),
		})
		return
	}

	// 1.封装一下条件
	//条件
	_condition := []string{"2", "3"}
	if request_data.Mode2 == 1 && request_data.Mode3 != 1 {
		_condition = []string{"2"}
	} else if request_data.Mode2 != 1 && request_data.Mode3 == 1 {
		_condition = []string{"3"}
	}

	//查询
	pageSize := 7
	offset := (int(request_data.Page) - 1) * pageSize
	var db_data []models_sqlite.WorkBasic
	err = models_sqlite.DB.Where("workid = ? and mode in (?)", request_data.Cid, _condition).Limit(pageSize).Offset(offset).Find(&db_data).Error //models_sqlite.DB.Where("workid = ? and mode in (?)", request_data.Cid, _condition).Limit(pageSize).Offset(offset).Find(db_data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "",
		})
	} else {

		response_data := Tcp_Tj1Data{}
		for _, v := range db_data {
			response_data.Data = append(response_data.Data, &Tcp_Tj1Info{
				Date:     v.Date,
				Mode:     int32(v.Mode),
				Number:   int32(v.Partnumber),
				MaxScore: float32(v.MaxScore),
				Score:    float32(v.Score),
				UniqueId: int32(v.Uniqueid),
			})
		}

		// 3. 查询一下最大页数
		var count int64
		err = models_sqlite.DB.Model(models_sqlite.WorkBasic{}).Count(&count).Error
		if err != nil {
			count = 0
		}
		var _maxPage int = 0
		if count%7 == 0 {
			_maxPage = int(count / 7)
		} else {
			_maxPage = int(count/7 + 1)
		}
		response_data.MaxPage = int32(_maxPage)

		jsonBytes, err := json.Marshal(response_data)
		if err != nil {
			log.Fatal(err)
		}
		jsonString := string(jsonBytes)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": jsonString,
		})

	}

}

func WorkScore(c *gin.Context) {

	request_data := new(Tcp_Tj2)
	err := c.ShouldBindJSON(request_data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常,err = " + err.Error(),
		})
		return
	}

	//查询
	pageSize := 9
	offset := (int(request_data.Page) - 1) * pageSize
	var db_data []models_sqlite.WorkRecordBasic

	tableName := fmt.Sprintf("workrecord_basic_%d", request_data.UniqueId)
	err = models_sqlite.DB.Table(tableName).Group("uname").Select("username,uname,sum(score) as score").Limit(pageSize).Offset(offset).Find(&db_data).Error
	//err = models_sqlite.DB.Find(db_data).Error //.Group("uname").Select("username,uname,sum(score)").Limit(pageSize).Offset(offset)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "",
		})
	} else {
		_maxCount := len(db_data)
		response_data := Tcp_Tj2Data{}
		for _, v := range db_data {
			response_data.Data = append(response_data.Data, &Tcp_Tj2Info{
				UName: v.Uname,
				Score: float32(v.Score),
			})
		}

		//读取一下最大页数
		var _maxPage int = 0
		if _maxCount%9 == 0 {
			_maxPage = _maxCount / 9
		} else {
			_maxPage = _maxCount/9 + 1
		}
		response_data.MaxNumber = int32(_maxCount)
		response_data.MaxPage = int32(_maxPage)

		jsonBytes, err := json.Marshal(response_data)
		if err != nil {
			log.Fatal(err)
		}
		jsonString := string(jsonBytes)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": jsonString,
		})

	}

}

func WorkDetail(c *gin.Context) {

	request_data := new(Tcp_Tj3)
	err := c.ShouldBindJSON(request_data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常,err = " + err.Error(),
		})
		return
	}

	//查询
	tableName := fmt.Sprintf("workrecord_basic_%d", request_data.UniqueId)
	var db_data []models_sqlite.WorkRecordBasic
	if len(request_data.UName) == 0 {
		err = models_sqlite.DB.Table(tableName).Find(&db_data).Error
	} else {
		err = models_sqlite.DB.Table(tableName).Where("uname = ?", request_data.UName).Find(&db_data).Error
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "",
		})
	} else {
		//_maxCount := len(db_data)
		response_data := Tcp_Tj3Data{}
		for _, v := range db_data {
			response_data.Data = append(response_data.Data, &Tcp_Tj3Info{
				UName:   v.Uname,
				Score:   float32(v.Score),
				Type:    int32(v.Type),
				Date:    v.Date,
				Content: v.Content,
				State:   int32(v.State),
			})
		}

		//补充request
		response_data.Number = request_data.Number
		response_data.WorkName = request_data.WorkName
		response_data.Mode = request_data.Mode
		response_data.MaxScore = float32(request_data.MaxScore)

		jsonBytes, err := json.Marshal(response_data)
		if err != nil {
			log.Fatal(err)
		}
		jsonString := string(jsonBytes)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": jsonString,
		})

	}

}

// 学生获取老师端状态
func WorkStatus(c *gin.Context) {

	response_data := Tcp_WorkStatus{}
	// 1. 查询老师状态
	if _, ok := core.WorldMgrObj.PMs["teacher"]; ok {
		//老师在线
		response_data.TeacherStatus = 1
	} else {
		response_data.TeacherStatus = 0
	}
	response_data.AutoStatus = core.WorldMgrObj.AutoStatus
	// 2. 推送消息
	jsonBytes, err := json.Marshal(response_data)
	if err != nil {
		log.Fatal(err)
	}
	jsonString := string(jsonBytes)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": jsonString,
	})

}

func WorkRecord(c *gin.Context) {

	request_data := new(Tcp_WorkInfoRecord)
	err := c.ShouldBindJSON(request_data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常,err = " + err.Error(),
		})
		return
	}

	//查询
	tableName := fmt.Sprintf("workrecord_basic_%d", request_data.UniqueId)
	exist := models_sqlite.DB.Migrator().HasTable(tableName)
	if exist {
		//存在

		// 1. 插入记录
		u := &models_sqlite.WorkRecordBasic{
			Username: request_data.Username,
			Uname:    request_data.Uname,
			Date:     request_data.Date,
			State:    int(request_data.State),
			Score:    int(request_data.Score),
			Type:     int(request_data.Type),
			Content:  request_data.Content,
		}
		if err := models_sqlite.DB.Table(tableName).Create(&u).Error; err != nil {
			fmt.Println("insert dir error")
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "操作记录插入失败" + err.Error(),
			})
		} else {
			// 1. 更新总表分数
			if request_data.Score > 0 {
				models_sqlite.DB.Model(&models_sqlite.WorkBasic{}).Where("uniqueid = ?", request_data.UniqueId).Update("score", gorm.Expr("score + ?", request_data.Score))
			}
			response_data := Tcp_ResponseRecordInfo{
				Code: 1,
			}
			jsonBytes, err := json.Marshal(response_data)
			if err != nil {
				log.Fatal(err)
			}
			jsonString := string(jsonBytes)
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"data": jsonString,
			})
		}
	}
}

func WorkQuestion(c *gin.Context) {

}
