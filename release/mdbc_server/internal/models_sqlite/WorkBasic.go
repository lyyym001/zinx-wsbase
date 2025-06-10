package models_sqlite

import "gorm.io/gorm"

//workid, workname,`date`,mode,partnumber,uniqueid,maxscore,score

type WorkBasic struct {
	gorm.Model
	WorkId     string `gorm:"column:workid;type:varchar(100);not null" json:"Workid"`
	Workname   string `gorm:"column:workname;type:varchar(100);not null" json:"Workname"`
	Date       string `gorm:"column:date;type:varchar(100);not null" json:"Date"`
	Mode       int    `gorm:"column:mode;type:int;not null" json:"Mode"`
	Partnumber int    `gorm:"column:partnumber;type:int;not null" json:"Partnumber"`
	Uniqueid   int    `gorm:"column:uniqueid;type:int;not null" json:"Uniqueid"`
	MaxScore   int    `gorm:"column:maxscore;type:int;not null" json:"MaxScore"`
	Score      int    `gorm:"column:score;type:int;not null" json:"Score"`
}

func (table *WorkBasic) TableName() string {
	return "work_basic"
}
