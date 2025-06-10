package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/internal/models_sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

func GetDirStatus(c *gin.Context) {

	var ds DirListResponse
	ds.Ds = map[string]DirInfoResponse{}
	ds.DirVersion = models_sqlite.DirVersion
	//获取目录清单
	var db_data []models_sqlite.DirBasic
	models_sqlite.DB.Find(&db_data)
	if db_data != nil {
		for _, device := range db_data {
			dinfo := DirInfoResponse{
				Did:   strconv.Itoa(device.Did),
				Sort:  device.Sort,
				DName: device.DName,
				Cs:    make(map[string]CourseInfoResponse),
			}
			ds.Ds[dinfo.Did] = dinfo
		}
	}

	//读取自定义课程清单
	var db_cdata []models_sqlite.CourseBasic
	models_sqlite.DB.Find(&db_cdata)
	if db_cdata != nil {
		for _, course := range db_cdata {
			dinfo := CourseInfoResponse{
				Rid:   course.Rid,
				RName: course.RName,
			}

			if _, ok := ds.Ds[strconv.Itoa(course.Did)]; ok {
				ds.Ds[strconv.Itoa(course.Did)].Cs[dinfo.Rid] = dinfo
			}
		}
	}

	fmt.Println("目录列表=", ds.Ds)

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

func SetDirName(c *gin.Context) {

	in := new(DirNameSetRequest)
	//fmt.Println("c = ", c.Request)
	err := c.ShouldBindJSON(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常,err = " + err.Error(),
		})
		return
	}

	if len(in.DName) < 1 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "改明后的名字不能为空",
		})
		return
	}

	//更新名称
	err = models_sqlite.DB.Model(&models_sqlite.DirBasic{}).Where("did = ?", in.Did).Update("dname", in.DName).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "",
		})
	} else {

		//更新目录版本
		models_sqlite.DB.Model(&models_sqlite.SysBasic{}).
			Where("sid = ?", 1).
			Update("dirversion", gorm.Expr("dirversion + ?", 1))
		models_sqlite.DirVersion += 1

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

func DirCreate(c *gin.Context) {

	//获取ID
	var newDid uint
	mode := &models_sqlite.DirBasic{}
	err := models_sqlite.DB.Last(mode).Error
	if err != nil {
		newDid = 1
	} else {
		newDid = mode.ID + 1
	}

	u := &models_sqlite.DirBasic{
		Did:   int(newDid),
		DName: "未命名",
		Sort:  0,
	}
	if err := models_sqlite.DB.Create(&u).Error; err != nil {
		fmt.Println("insert dir error")
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "创建目录失败" + err.Error(),
		})
	} else {

		//更新目录版本
		models_sqlite.DB.Model(&models_sqlite.SysBasic{}).
			Where("sid = ?", 1).
			Update("dirversion", gorm.Expr("dirversion + ?", 1))
		models_sqlite.DirVersion += 1

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": strconv.Itoa(int(newDid)),
		})
	}

}

func DirDelete(c *gin.Context) {

	in := new(DirNameSetRequest)
	//fmt.Println("c = ", c.Request)
	err := c.ShouldBindJSON(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常,err = " + err.Error(),
		})
		return
	}

	//删除
	err = models_sqlite.DB.Where("did = ?", in.Did).Unscoped().Delete(&models_sqlite.DirBasic{}).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "",
		})
	} else {

		//更新目录版本
		models_sqlite.DB.Model(&models_sqlite.SysBasic{}).
			Where("sid = ?", 1).
			Update("dirversion", gorm.Expr("dirversion + ?", 1))
		models_sqlite.DirVersion += 1

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": in.Did,
		})
	}
}

func CourseChangeDir(c *gin.Context) {

	in := new(DirCourseRequest)
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
	err = models_sqlite.DB.Model(models_sqlite.CourseBasic{}).Where("rid = ?", in.Cid).Count(&count).Error
	if err != nil || count == 0 {
		//创建
		if err := models_sqlite.DB.Create(&models_sqlite.CourseBasic{
			Did: in.Did,
			Rid: in.Cid,
		}).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"data": "",
			})
		} else {

			//更新目录版本
			models_sqlite.DB.Model(&models_sqlite.SysBasic{}).
				Where("sid = ?", 1).
				Update("dirversion", gorm.Expr("dirversion + ?", 1))
			models_sqlite.DirVersion += 1
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"data": "",
			})
		}
	} else {
		//更新
		err = models_sqlite.DB.Model(&models_sqlite.CourseBasic{}).Where("rid = ?", in.Cid).Update("did", in.Did).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"data": "",
			})
		} else {

			//更新目录版本
			models_sqlite.DB.Model(&models_sqlite.SysBasic{}).
				Where("sid = ?", 1).
				Update("dirversion", gorm.Expr("dirversion + ?", 1))
			models_sqlite.DirVersion += 1
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"data": "",
			})
		}
	}

}
