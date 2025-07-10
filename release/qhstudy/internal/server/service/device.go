package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/lyyym/zinx-wsbase/release/qhstudy/core"
	"github.com/lyyym/zinx-wsbase/release/qhstudy/internal/models_sqlite"
	"log"
	"net/http"
)

func GetDeviceStatus(c *gin.Context) {

	var ds DeviceListResponse
	ds.Ds = map[string]DeviceInfoResponse{}
	//获取设备清单
	var db_data []models_sqlite.DeviceBasic
	models_sqlite.DB.Find(&db_data)
	if db_data != nil {
		for _, device := range db_data {
			dinfo := DeviceInfoResponse{
				Username: device.Username,
				Status:   device.Status,
				Battery:  0,
				Free:     0,
				//Ip:       device.Ip,
			}
			p := core.WorldMgrObj.GetPlayerByUserName(device.Username)
			if p != nil {
				//在线
				dinfo.Ip = p.CDevice.Ip
				dinfo.Battery = p.CDevice.Battery
				dinfo.Free = p.CDevice.Free
				if dinfo.Status != 0 {
					dinfo.Status = 2
				}
			}
			ds.Ds[dinfo.Username] = dinfo
		}
	}
	//fmt.Println("设备列表=", ds.Ds)

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

func GetDeviceDetail(c *gin.Context) {

	in := new(DeviceDetailRequest)
	//fmt.Println("c = ", c.Request)
	err := c.ShouldBindJSON(in)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常,err = " + err.Error(),
		})
		return
	}

	var ds DeviceDetailResponse
	ds.Username = in.Username
	p := core.WorldMgrObj.GetPlayerByUserName(in.Username)
	if p != nil {
		if p.CDevice.Dir == models_sqlite.DirVersion {
			ds.Dir = 2
		} else {
			ds.Dir = 1
		}
		ds.Ip = p.CDevice.Ip
	} else {
		ds.Dir = 0
		ds.Ip = "未读取到"
	}

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
