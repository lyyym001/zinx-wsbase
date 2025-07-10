package models_sqlite

import "gorm.io/gorm"

type DirBasic struct {
	gorm.Model
	Did   int    `gorm:"column:did;type:int;uniqueIndex;not null" json:"did"`
	Sort  int    `gorm:"column:sort;type:int;not null;default=1" json:"sort"` //排序
	DName string `gorm:"column:dname;type:varchar(100)" json:"dname"`         //IP地址
}

func (table *DirBasic) TableName() string {
	return "dir_basic"
}
