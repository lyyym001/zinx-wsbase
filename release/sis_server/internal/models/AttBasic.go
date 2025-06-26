package models

import "gorm.io/gorm"

type AttBasic struct {
	gorm.Model
	Username string `gorm:"column:username;type:varchar(100);not null" json:"username"`
	Flag     int32  `gorm:"column:flag;type:int;not null" json:"flag"`
}

func (table *AttBasic) TableName() string {
	return "att_basic"
}
