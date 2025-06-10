package models_sqlite

import "gorm.io/gorm"

type DeviceBasic struct {
	gorm.Model
	Username string `gorm:"column:username;type:varchar(100);uniqueIndex;not null" json:"username"`
	Status   int32  `gorm:"column:status;type:int;not null;default=1" json:"status"` //0-故障 1-正常
	Ip       string `gorm:"column:ip;type:varchar(100)" json:"ip"`                   //IP地址
}

func (table *DeviceBasic) TableName() string {
	return "device_basic"
}
