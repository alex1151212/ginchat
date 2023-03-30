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

func AddFriend(userId uint, targetName string) (int, string) {
	// user := UserBasic{}
	if targetName != "" {
		targetUser := FindUserByName(targetName)
		if targetUser.Salt != "" {
			if userId == targetUser.ID {
				return -1, "無法新增自己好友"
			}
			contact0 := Contact{}
			utils.DB.Where("owner_id = ? and target_id = ? and type = 1 ", userId, targetUser.ID).Find(&contact0)
			if contact0.ID != 0 {
				return -1, "無法重複添加好友"
			}

			tx := utils.DB.Begin()
			// DB事務錯誤Rollback
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()
			contact := Contact{}
			contact.OwnerId = userId
			contact.TargetId = targetUser.ID
			contact.Type = 1
			if err := utils.DB.Create(&contact).Error; err != nil {
				tx.Rollback()
				return -1, "新增好友失敗"
			}
			contact1 := Contact{}
			contact1.OwnerId = targetUser.ID
			contact1.TargetId = userId
			contact1.Type = 1
			if err := utils.DB.Create(&contact1).Error; err != nil {
				tx.Rollback()
				return -1, "新增好友失敗"
			}
			tx.Commit()
			return 0, "新增好友成功"
		}
		return -1, "查無此用戶"
	}
	return -1, "好友名稱不能為空"
}

func SearchUserByGroupId(communityId uint) []uint {
	contacts := make([]Contact, 0)
	objIds := make([]uint, 0)
	utils.DB.Where("target_id = ? and type = 2", communityId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, uint(v.OwnerId))
	}
	return objIds
}
