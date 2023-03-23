package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//	 GetUserList
//	 @Summary 所有用戶
//		@Tags		用戶
//		@Success	200	{string}	json{"code","message"}
//		@Router		/user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		"message": data,
	})
}

//	 CreateUser
//	 @Summary 新增用戶
//		@Tags		用戶
//		@Param name query string false		"使用者名稱"
//		@Param password query string false		"使用者密碼"
//		@Param repassword query string false		"使用者確認密碼"
//		@Success	200	{string}	json{"code","message"}
//		@Router		/user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("repassword")
	fmt.Println(user.Name, "  >>>>>>>>>>>  ", password, repassword)
	salt := fmt.Sprintf("%06d", rand.Int31())

	data := models.FindUserByName(user.Name)
	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //0成功 -1失敗
			"message": "使用者名稱或密碼不能為空",
			"data":    data,
		})
		return
	}
	if data.Name != "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //0成功 -1失敗
			"message": "使用者名稱已註冊",
			"data":    data,
		})
		return
	}
	if password != repassword {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //0成功 -1失敗
			"message": "兩次密碼不一致",
		})
		return
	}

	user.Password = utils.MakePasssword(password, salt)
	user.Salt = salt
	models.CreateUser(user)
	c.JSON(http.StatusOK, gin.H{
		"code":    0, //0成功 -1失敗
		"message": "創建用戶成功",
		"data":    data,
	})
}

//	 DeleteUser
//	 @Summary 刪除用戶
//		@Tags		用戶
//		@Param id query string false		"使用者ID"
//		@Success	200	{string}	json{"code","message"}
//		@Router		/user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(http.StatusOK, gin.H{
		"code":    0, //0成功 -1失敗
		"message": "刪除用戶成功",
		"data":    user,
	})

}

//	 UpdateUser
//	 @Summary 編輯用戶
//		@Tags		用戶
//		@Param id formData string false		"使用者ID"
//		@Param name formData string false		"使用者名稱"
//		@Param password formData string false		"使用者密碼"
//		@Param phone formData string false		"電話號碼"
//		@Param email formData string false		"電子信箱"
//		@Success	200	{string}	json{"code","message"}
//		@Router		/user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	fmt.Println("update:", user)

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //0成功 -1失敗
			"message": "編輯用戶失敗",
			"data":    user,
		})
		return
	}
	models.UpdateUser(user)
	c.JSON(http.StatusOK, gin.H{
		"code":    0, //0成功 -1失敗
		"message": "編輯用戶成功",
		"data":    user,
	})

}

//	 FindUserByNameAndPwd
//	 @Summary 所有用戶
//		@Tags		用戶
//		@Param name query string false		"使用者名稱"
//		@Param password query string false		"使用者密碼"
//		@Success	200	{string}	json{"code","message"}
//		@Router		/user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}
	name := c.PostForm("name")
	password := c.PostForm("password")

	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //0成功 -1失敗
			"message": "該使用者不存在",
			"data":    data,
		})
		return
	}
	flag := utils.ValidPassword(password, user.Salt, user.Password)
	if !flag {
		c.JSON(http.StatusOK, gin.H{
			"message": "密碼不正確",
		})
		return
	}
	enCodePwd := utils.MakePasssword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, enCodePwd)
	c.JSON(http.StatusOK, gin.H{
		"code":    0, //0成功 -1失敗
		"message": "登入成功",
		"data":    data,
	})
}

// 防止跨域偽造請求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws, c)

}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	msg, err := utils.Subscribe(c, utils.PublishKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("發送訊息: ", msg)
	tm := time.Now().Format("YYYY-MM-DD HH:mm:ss")
	m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
	err = ws.WriteMessage(1, []byte(m))
	if err != nil {
		fmt.Println(err)
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
