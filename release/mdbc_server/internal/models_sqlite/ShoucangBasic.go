package models_sqlite

import "gorm.io/gorm"

type ShoucangBasic struct {
	gorm.Model
	Cid    string `gorm:"column:cid;type:varchar(100);uniqueIndex;not null" json:"cid"`
	Status int    `gorm:"column:status;type:int;not null" json:"status"` //0-未收藏 1-收藏
	Date   string `gorm:"column:date;type:varchar(100);not null" json:"date"`
	CName  string `gorm:"column:cname;type:varchar(100);not null" json:"cname"`
}

func (table *ShoucangBasic) TableName() string {
	return "shoucang_basic"
}
