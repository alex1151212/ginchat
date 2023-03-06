package models

import (
	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	OwnerId  uint
	TargetId uint
	Type     int //對應類型: 0 1 2
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}
