package service

import (
	"ginchat/models"
	"html/template"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetIndex
//
//	@Tags		前端網頁
//	@Success	200	{string}	welcome
//	@Router		/index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "index")
}

// ToRegister
//
//	@Tags		前端網頁
//	@Success	200	{string}	welcome
//	@Router		/toRegister [get]
func ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "register")
}

// ToChat
//
//	@Tags		前端網頁
//	@Success	200	{string}	welcome
//	@Router		/toChat [get]
func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/foot.html",
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/createcom.html",
		"views/chat/userinfo.html",
		"views/chat/main.html")
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token
	ind.Execute(c.Writer, user)
}
