package models

import (
	"fmt"
	"ginchat/utils"
	"time"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string
	Password      string
	Phone         string `valid:"matches(^09\\d{8}$)"`
	Email         string `valid:"email"`
	Identity      string
	ClientIp      string
	ClientPort    string
	Salt          string
	LoginTime     *time.Time
	HeartbeatTime *time.Time
	LogoutTime    *time.Time
	IsLogout      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}
func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}
func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("Name = ?", name).First(&user)

	// Encode token
	str := fmt.Sprintf("%d", time.Now().Unix())
	token := utils.Md5Encode(str)

	utils.DB.Model(&user).Where("Id = ? ", user.ID).Update("identity", token)
	return user
}
func FindUserByNameAndPwd(name string, password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("Name = ? and password = ? ", name, password).First(&user)
	return user
}
func FindUserByPhone(phone string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("Phone = ?", phone).First(&user)
	return user
}
func FindUserByEmail(email string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("Email = ?", email).First(&user)
	return user
}

func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}
func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{
		Name:     user.Name,
		Password: user.Password,
		Phone:    user.Phone,
		Email:    user.Email,
	})
}
