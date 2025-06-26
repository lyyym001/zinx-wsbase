package models

import "gorm.io/gorm"

type UserBasic struct {
	gorm.Model
	Username string `gorm:"column:username;type:varchar(100);uniqueIndex;not null" json:"username"`
	Password string `gorm:"column:password;type:varchar(36);not null" json:"password"`
	//Power       int32  `gorm:"column:power;type:uint;not null" json:"power"`
	Role        uint32 `gorm:"column:role;type:int;not null" json:"role"`
	Identify    string `gorm:"column:identify;type:varchar(100);not null" json:"uid"`
	Sdp         string `gorm:"column:sdp;type:text" json:"sdp"`
	AccountType int32  `gorm:"column:account_type;type:int;not null" json:"account_type"`
	NickName    string `gorm:"column:nick_name;type:varchar(100);not null" json:"nick_name"`
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}
