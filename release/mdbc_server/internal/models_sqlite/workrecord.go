package models_sqlite

import (
	"fmt"
	"gorm.io/gorm"
)

//id INTEGER PRIMARY KEY AUTOINCREMENT,
//username STRING NOT NULL,
//uname STRING NOT NULL,
//date STRING NOT NULL,
//type INT NOT NULL,
//state INT NOT NULL,
//score FLOAT NOT NULL,
//content STRING NOT NULL

type WorkRecordBasic struct {
	gorm.Model
	Username string `gorm:"column:username;type:varchar(100);not null" json:"Username"`
	Uname    string `gorm:"column:uname;type:varchar(100);not null" json:"Uname"`
	Date     string `gorm:"column:date;type:varchar(100);not null" json:"Date"`
	State    int    `gorm:"column:state;type:int;not null" json:"Dtate"`
	Score    int    `gorm:"column:score;type:int;not null" json:"Score"`
	Type     int    `gorm:"column:type;type:int;not null" json:"Type"`
	Content  string `gorm:"column:content;type:varchar(500);not null" json:"Content"`
}

type WorkRecordInfoBasic struct {
	Username string
	Uname    string
	Score    float32
}

func (table *WorkRecordBasic) TableName(uniqueid int) string {

	return fmt.Sprintf("workrecord_basic_%d", uniqueid)
}
