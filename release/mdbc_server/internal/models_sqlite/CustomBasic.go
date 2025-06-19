package models_sqlite

import "gorm.io/gorm"

type CustomBasic struct {
	gorm.Model
	Rid        string `gorm:"column:rid;type:varchar(100);uniqueIndex;not null" json:"rid"`
	Did        int    `gorm:"column:did;type:int;not null;default=1" json:"did"`               //所属目录
	RName      string `gorm:"column:rname;type:varchar(100)" json:"rname"`                     //课程名称
	CourseType int    `gorm:"column:coursetype;type:int;not null;default=1" json:"coursetype"` //课程类型
	Stereo     int    `gorm:"column:stereo;type:int;not null;default=1" json:"stereo"`         //视频类型
}

func (table *CustomBasic) TableName() string {
	return "custom_basic"
}
