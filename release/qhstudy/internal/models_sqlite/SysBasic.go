package models_sqlite

import "gorm.io/gorm"

type SysBasic struct {
	gorm.Model
	Sid        int   `gorm:"column:sid;type:int;not null" json:"sid"`
	DirVersion int64 `gorm:"column:dirversion;type:int64;not null" json:"dirversion"`
}

func (table *SysBasic) TableName() string {
	return "sys_basic"
}
