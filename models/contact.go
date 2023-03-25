package models

import (
	"fmt"
	"ginchat/utils"

	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	OwnerId  uint
	TargetId uint
	Type     int //對應類型:  1好友 2群組 3
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}

func SearchFriend(userId uint) []UserBasic {

	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	fmt.Println(contacts)
	utils.DB.Where("owner_id = ? and type=1", userId).Find(&contacts)
	fmt.Println(contacts)
	for _, v := range contacts {
		objIds = append(objIds, uint64(v.TargetId))
	}
	users := make([]UserBasic, 0)
	utils.DB.Where("id in ?", objIds).Find(&users)
	return users
}
