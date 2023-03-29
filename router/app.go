package router

import (
	docs "ginchat/docs"
	"ginchat/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {

	r := gin.Default()

	//swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//靜態資源
	r.Static("/asset", "asset/")
	r.LoadHTMLGlob("views/**/*")

	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/toRegister", service.ToRegister)
	r.GET("/toChat", service.ToChat)
	r.GET("/chat", service.Chat)

	r.POST("/searchFriends", service.SearchFriends)

	r.GET("/user/getUserList", service.GetUserList)
	r.POST("/user/createUser", service.CreateUser)
	r.DELETE("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("/user/findUserByNameAndPwd", service.FindUserByNameAndPwd)

	//發送訊息
	r.GET("/user/sendMsg", service.SendMsg)
	r.GET("/user/sendUserMsg", service.SendUserMsg)

	//上傳文件
	r.POST("/attach/upload", service.Upload)

	//添加好友
	r.POST("/contact/addfriend", service.AddFriend)

	//創建群組
	r.POST("/contact/createCommunity", service.CreateCommunity)
	return r
}
