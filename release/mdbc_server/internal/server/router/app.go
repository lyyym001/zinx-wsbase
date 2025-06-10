package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/internal/middlewares"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/internal/server/service"
	"github.com/lyyym/zinx-wsbase/ziface"
)

func Router(server ziface.IServer) *gin.Engine {

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//跨域
	r.Use(middlewares.Cors())

	/*
		websocket
	*/
	r.GET("/ws", server.Serve)

	/*
		注册
	*/
	r.POST("/user/register", service.UserRegister)
	/*
		登录
	*/
	r.POST("/user/login", service.UserLogin)
	/*
		鉴权
	*/
	auth := r.Group("/auth", middlewares.Auth())

	//读取用户配置数据
	auth.POST("/user/info", service.GetUserConfig)

	//设备状态
	auth.POST("/device/list", service.GetDeviceStatus)

	//设备详情
	auth.POST("/device/detail", service.GetDeviceDetail)

	//设备详情
	auth.POST("/dir/list", service.GetDirStatus)

	//目录改名
	auth.POST("/dir/altername", service.SetDirName)

	//目录创建
	auth.POST("/dir/create", service.DirCreate)

	//目录删除
	auth.POST("/dir/delete", service.DirDelete)

	//更改目录
	auth.POST("/course/changedir", service.CourseChangeDir)

	//课程收藏
	auth.POST("/course/shoucang", service.CourseShoucang)

	//课程收藏列表
	auth.POST("/course/shoucanglist", service.GetShoucang)

	//课程记录
	auth.POST("/course/record", service.CourseRecord)

	//课程播放记录列表
	auth.POST("/course/recordlist", service.GetRecord)

	//课程播放记录列表
	auth.POST("/user/online", service.GetUserStatus)

	//课件播放列表数据
	auth.POST("/work/recordlist", service.WorkList)

	//课件播放分数列表
	auth.POST("/work/recordscore", service.WorkScore)

	//课件播放详情数据
	auth.POST("/work/recorddetail", service.WorkDetail)

	//课件播放详情数据
	auth.POST("/work/record", service.WorkRecord)

	//学生获取老师端状态
	auth.POST("/work/status", service.WorkStatus)

	return r
}
