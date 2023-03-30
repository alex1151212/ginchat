package models

import (
	"encoding/json"
	"fmt"
	"ginchat/utils"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FromId   int64  //發送者
	TargetId int64  //接收者
	Type     int    //發送類型: 1.私訊聊天 2.群組聊天 3.廣播
	Media    int    //訊息類型: 1.文字 2.表情包 3.圖片 4.音訊
	Content  string //訊息內容
	Pic      string
	Url      string
	Desc     string
	Amount   int //其他數字統計
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// Map Contact
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 讀寫鎖
var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {
	// 1.獲取參數並進行token驗證
	// token := query.Get("token")
	query := request.URL.Query()
	id := query.Get("userId")
	userId, _ := strconv.ParseInt(id, 10, 64)
	// msgType := query.Get("type")
	// targetId := query.Get("targetId")
	// context := query.Get("context")
	isValida := true // TODO checkToken()

	conn, err := (&websocket.Upgrader{
		//token驗證
		CheckOrigin: func(r *http.Request) bool {
			return isValida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 2.獲取conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	// 3.使用者關係
	// 4.userId 跟 node 綁定並加鎖
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	// 5.完成發送邏輯
	go sendProc(node)
	// 6.完成接收邏輯
	go recvProc(node)

	sendMsg(userId, []byte("歡迎進入聊天室!"))

}
func sendProc(node *Node) {
	for {

		select {
		case data := <-node.DataQueue:
			fmt.Println("[ws] sendProc >>>> msg :", string(data))
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(data)
		broadMsg(data)
		fmt.Println("[ws] resvProc <<<<< ", string(data))
	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecvProc()
	fmt.Println("init goroutine")
}

// 完成udp資料發送協程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 255),
		Port: viper.GetInt("port.udp"),
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-udpsendChan:
			fmt.Println("udpSendProc data :", string(data))
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: viper.GetInt("port.udp"),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("udpRecvProc data :", string(buf[0:]))

		dispatch(buf[0:n])

	}
}

// 後端調度邏輯
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: //私訊
		fmt.Println("dispatch data :", string(data))
		sendMsg(msg.TargetId, data)
	case 2: //群組聊天
		sendGroupMsg(msg.TargetId, data)
		// case 3:sendAllMsg() //廣播
		// case 4:

	}
}

func sendGroupMsg(targetId int64, msg []byte) {
	fmt.Println("開始群組訊息發送")
	userIds := SearchUserByGroupId(uint(targetId))
	for _, userId := range userIds {
		if targetId != int64(userId) {
			sendMsg(int64(userId), msg)
		}
	}
}

func sendMsg(userId int64, msg []byte) {
	fmt.Println("sendMsg >>> userID: ", userId, " msg:", string(msg))
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}

func JoinGroup(userId uint, comId string) (int, string) {
	contact := Contact{}
	contact.OwnerId = userId
	contact.Type = 2
	community := Community{}

	utils.DB.Where("id = ? or name= ? ", comId, comId).Find(&community)
	if community.Name == "" {
		return -1, "沒有找到此群組"
	}
	utils.DB.Where("owner_id =  ? and target_id = ? and  type = 2 ", userId, comId).Find(&contact)
	if !contact.CreatedAt.IsZero() {
		return -1, "已加過此群組"
	} else {
		contact.TargetId = community.ID
		utils.DB.Create(&contact)
		return 0, "加入群組成功"
	}
}
