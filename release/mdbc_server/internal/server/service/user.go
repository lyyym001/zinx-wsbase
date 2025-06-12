package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lyyym/zinx-wsbase/global"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/core"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/internal/helper"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/internal/models"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/internal/models_sqlite"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/pb"
	"log"
	"net/http"
	"strings"
)

func UserRegister(c *gin.Context) {
	in := new(UserRegisterRequest)
	//fmt.Println("c = ", c.Request)
	err := c.ShouldBindJSON(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常,err = " + err.Error(),
		})
		return
	}
	if in.Username == "" || in.Password == "" || in.NickName == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "必填信息为空",
		})
		return
	}
	in.Password = helper.GetMd5(in.Password)

	//写入对抗数据
	u := models.UserBasic{
		Username:    in.Username,
		Password:    in.Password,
		NickName:    in.NickName,
		Role:        in.Role,
		AccountType: 0,
		Identify:    in.Username,
	}
	if err := models.DB.Create(&u).Where("username != ", in.Username).Error; err != nil {
		fmt.Println("insert user error")
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "注册失败" + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "注册成功",
		})
	}

}

func UserLogin(c *gin.Context) {
	in := new(UserLoginRequest)
	//fmt.Println("c = ", c.Request)
	err := c.ShouldBindJSON(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常,err = " + err.Error(),
		})
		return
	}
	if in.Username == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "必填信息为空",
		})
		return
	}

	token, err := helper.GenerateToken(in.Username, in.AccountType, in.NickName, in.Ip)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "GenerateToken Error:" + err.Error(),
		})
		return
	}

	//insertUsers
	if in.AccountType == 1 && len(in.Members) > 0 {
		//老师
		fmt.Println("老师删除不存在的users", in.Members)
		//var dd []models_sqlite.DeviceBasic
		models_sqlite.DB.Where("username not in (?) ", in.Members).Unscoped().Delete(&models_sqlite.DeviceBasic{})

		//写入设备清单
		devices := strings.Split(in.Members, ",")
		for _, value := range devices {
			u := models_sqlite.DeviceBasic{
				Username: value,
				Status:   1,
			}
			if err := models_sqlite.DB.Create(&u).Where("username != ", value).Error; err != nil {
				fmt.Println("insert device error")
			}
		}

		//写入目录清单
		var count int64
		models_sqlite.DB.Model(&models_sqlite.DirBasic{}).Count(&count)
		fmt.Println("count = ", count)
		if count == 0 && in.Dirs != nil && len(in.Dirs) > 0 {
			for _, value := range in.Dirs {
				u := models_sqlite.DirBasic{
					Did:   value.Did,
					Sort:  value.Sort,
					DName: value.DName,
				}
				if err := models_sqlite.DB.Create(&u).Where("did != ", value.Did).Error; err != nil {
					fmt.Println("insert device error")
				}
			}
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": token,
	})

}

func GetUserConfig(c *gin.Context) {

	//uc := c.MustGet("user_claims").(*helper.UserClaims)
	//if uc == nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"code": -1,
	//		"msg":  "tokenError",
	//	})
	//	return
	//}
	//
	//data := new(models.UserBasic)
	//err := models.DB.Where("id = ?", uc.Id).First(&data).Error
	//if err != nil {
	//	if err == gorm.ErrRecordNotFound {
	//		c.JSON(http.StatusOK, gin.H{
	//			"code": -1,
	//			"msg":  "id不存在",
	//		})
	//		return
	//	}
	//	c.JSON(http.StatusOK, gin.H{
	//		"code": -1,
	//		"msg":  "Get UserBasic Error:" + err.Error(),
	//	})
	//	return
	//}

	jsonData := pb.UserConfig{
		RtmpHost:     global.Object.Conf.RtmpHost, // config.YamlConfig.Conf.RtmpHost,
		StreamingUri: global.Object.Conf.StreamingUri,
		RtmpChannel:  global.Object.Conf.RtmpChannel,
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		log.Fatal(err)
	}
	jsonString := string(jsonBytes)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": jsonString,
	})
}

func GetUserStatus(c *gin.Context) {

	var ds UserListResponse
	//获取成员
	userList := core.WorldMgrObj.GetAllPlayers()
	for _, user := range userList {
		if user.CDevice.Status == 1 {
			ds.Us = append(ds.Us, user.UserName)
		}
	}

	//fmt.Println("在线设备列表(不包含故障)=", ds.Us)

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
