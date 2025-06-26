package models

import (
	"github.com/lyyym/zinx-wsbase/release/sis_server/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func NewDB() {
	//dsn := string(config.YamlConfig.Mysql.Dns)
	//"root:root@tcp(127.0.0.1:3306)/mspace?charset=utf8mb4&parseTime=True&loc=Local
	//"root:root@tcp(127.0.0.1:3306)/mspace?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(config.YamlConfig.Mysql.Dns), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&UserBasic{}, &AttBasic{}) //&RoomBasic{}, &RoomUser{},
	DB = db
}
