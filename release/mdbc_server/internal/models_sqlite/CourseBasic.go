package models_sqlite

import "gorm.io/gorm"

type CourseBasic struct {
	gorm.Model
	Rid   string `gorm:"column:rid;type:varchar(100);uniqueIndex;not null" json:"rid"`
	Did   int    `gorm:"column:did;type:int;not null;default=1" json:"did"` //所属目录
	RName string `gorm:"column:rname;type:varchar(100)" json:"rname"`       //课程名称
}

func (table *CourseBasic) TableName() string {
	return "course_basic"
}
