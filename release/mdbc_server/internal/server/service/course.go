package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/internal/models_sqlite"
	"log"
	"net/http"
)

func CourseShoucang(c *gin.Context) {

	in := new(ScCourseRequest)
	//fmt.Println("c = ", c.Request)
	err := c.ShouldBindJSON(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常,err = " + err.Error(),
		})
		return
	}

	var count int64
	err = models_sqlite.DB.Model(models_sqlite.ShoucangBasic{}).Where("cid = ?", in.Cid).Count(&count).Error
	if err != nil || count == 0 {
		//收藏
		if in.Status != 1 {
			return
		}
		//创建
		if err := models_sqlite.DB.Create(&models_sqlite.ShoucangBasic{
			Cid:    in.Cid,
			Status: 1,
			Date:   in.Date,
			CName:  in.CName,
		}).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"data": "",
			})
		} else {
			jsonBytes, err := json.Marshal(in)
			if err != nil {
				log.Fatal(err)
			}
			jsonString := string(jsonBytes)
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"data": jsonString,
			})
		}
	} else {
		//取消收藏
		if in.Status != 0 {
			return
		}
		err = models_sqlite.DB.Where("cid = ?", in.Cid).Unscoped().Delete(&models_sqlite.ShoucangBasic{}).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"data": "",
			})
		} else {
			jsonBytes, err := json.Marshal(in)
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

func GetShoucang(c *gin.Context) {

	var ds ScCourseListRequest
	//获取目录清单
	var db_data []models_sqlite.ShoucangBasic
	models_sqlite.DB.Find(&db_data)
	if db_data != nil {
		for _, device := range db_data {
			dinfo := ScCourseRequest{
				Cid:   device.Cid,
				CName: device.CName,
				Date:  device.Date,
			}
			ds.Scs = append(ds.Scs, dinfo)
		}
	}

	fmt.Println("收藏列表=", ds.Scs)

	jsonBytes, err := json.Marshal(ds)
	if err != nil {
		log.Fatal(err)
	}
	jsonString := string(jsonBytes)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": jsonString,
	})
}

func CourseRecord(c *gin.Context) {

	in := new(ScCourseRequest)
	//fmt.Println("c = ", c.Request)
	err := c.ShouldBindJSON(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常,err = " + err.Error(),
		})
		return
	}

	//创建
	if err := models_sqlite.DB.Create(&models_sqlite.RecordBasic{
		Cid:   in.Cid,
		Date:  in.Date,
		CName: in.CName,
	}).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "",
		})
	} else {
		jsonBytes, err := json.Marshal(in)
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

func GetRecord(c *gin.Context) {

	var ds ScCourseListRequest
	//获取目录清单
	var db_data []models_sqlite.RecordBasic
	models_sqlite.DB.Order("date desc").Limit(20).Find(&db_data)
	if db_data != nil {
		for _, device := range db_data {
			dinfo := ScCourseRequest{
				Cid:   device.Cid,
				CName: device.CName,
				Date:  device.Date,
			}
			ds.Scs = append(ds.Scs, dinfo)
		}
	}

	fmt.Println("播放列表=", ds.Scs)

	jsonBytes, err := json.Marshal(ds)
	if err != nil {
		log.Fatal(err)
	}
	jsonString := string(jsonBytes)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": jsonString,
	})
}

func CourseCreate(c *gin.Context) {

	in := new(NewCourseRequest)
	//fmt.Println("c = ", c.Request)
	err := c.ShouldBindJSON(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常,err = " + err.Error(),
		})
		return
	}
	fmt.Println("in = ", in)
	//in.Cid = "CUSTOM" + in.Cid
	//查询一下是否存在该ID
	var count int64
	err = models_sqlite.DB.Model(models_sqlite.CustomBasic{}).Where("RID = ?", in.Cid).Count(&count).Error
	if err != nil {
		count = 0
	}
	u := &models_sqlite.CustomBasic{
		Rid:        in.Cid,
		RName:      in.CName,
		CourseType: in.CType,
		Stereo:     in.VType,
		Did:        999,
	}
	jsonBytes, err := json.Marshal(in)
	if err != nil {
		log.Fatal(err)
	}
	jsonString := string(jsonBytes)
	fmt.Println("count = ", count)
	//修改数据
	if count == 0 {
		//插入
		if err := models_sqlite.DB.Create(&u).Error; err != nil {
			fmt.Println("insert course error")
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "创建课程失败" + err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "创建课程成功",
				"data": jsonString,
			})
		}
	} else {
		err = models_sqlite.DB.Model(&models_sqlite.CustomBasic{}).Where("rid = ?", in.Cid).
			Update("RName", in.CName).
			Update("CourseType", in.CType).
			Update("Stereo", in.VType).
			Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "创建课程失败" + err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "创建课程成功",
				"data": jsonString,
			})
		}
	}
}

func GetCustomCourse(c *gin.Context) {

	in := new(NewCourseRequest)
	//fmt.Println("c = ", c.Request)
	err := c.ShouldBindJSON(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常,err = " + err.Error(),
		})
		return
	}
	fmt.Println("in = ", in)
	//in.Cid = "CUSTOM" + in.Cid
	//查询一下是否存在该ID
	var count int64
	err = models_sqlite.DB.Model(models_sqlite.CustomBasic{}).Where("RID = ?", in.Cid).Count(&count).Error
	if err != nil {
		count = 0
	}
	fmt.Println("count = ", count)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "课程查找成功",
		"data": count,
	})
}
