package models

import (
	"fmt"
	"ginchat/utils"

	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name    string
	OwnerId uint
	Img     string
	Desc    string
}

func CreateCommunity(community Community) (int, string) {
	if len(community.Name) == 0 {
		return -1, "群組名稱不能為空"
	}
	if community.OwnerId == 0 {
		return -1, "請先登入"
	}
	if err := utils.DB.Create(&community).Error; err != nil {
		fmt.Println(err)
		return -1, "建立群組失敗"
	}
	return 0, "建立群組成功"
}
