package models_sqlite

import "gorm.io/gorm"

type RecordBasic struct {
	gorm.Model
	Cid   string `gorm:"column:cid;type:varchar(100);not null" json:"cid"`
	Date  string `gorm:"column:date;type:varchar(100);not null" json:"date"`
	CName string `gorm:"column:cname;type:varchar(100);not null" json:"cname"`
}

func (table *RecordBasic) TableName() string {
	return "record_basic"
}
