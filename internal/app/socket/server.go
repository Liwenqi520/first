package socket

import (
	"encoding/json"
	"first/internal/schema"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Liwenqi520/errorx"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func Init(heart int) {
	heartBeat = time.Duration(heart) * time.Second
	WsSessions = make(map[string]*WsSession)
}

var (
	//心跳
	heartBeat = 60 * time.Second
)

var Upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     CheckOrigin,
}

func CheckOrigin(r *http.Request) bool {
	return true
}

type WsMessage struct {
	MessageType int
	Data        []byte
}

// WsSessions 存储在线的ws session
var WsSessions map[string]*WsSession

type WsSession struct {
	ID        string
	Conn      *websocket.Conn
	InChan    chan *WsMessage
	OutChan   chan *WsMessage
	mu        sync.Mutex
	isClosed  bool
	CloseChan chan byte
}

type Dealer interface {
	DealMsg(msg []byte) []byte
}

func (ws *WsSession) WsProcessLoop(dealerMsg Dealer) {
	for {
		//读取消息队列的消息
		msg, err := ws.readInChan()
		if err != nil {
			break
		}
		//处理消息
		var newMsg []byte
		if dealerMsg != nil {
			newMsg = dealerMsg.DealMsg(msg.Data)
		}

		err = ws.writeOutChan(msg.MessageType, newMsg)
		if err != nil {
			break
		}
	}

}

func (ws *WsSession) WsReadLoop() {

	for {
		fmt.Println("开始读取msg")
		//心跳
		if err := ws.Conn.SetReadDeadline(time.Now().Add(heartBeat)); err != nil {
			fmt.Println("心跳设置失败")
		}

		msgType, data, err := ws.Conn.ReadMessage()
		if err != nil {
			fmt.Println("读取消息失败", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseGoingAway) {
				fmt.Println("1001 错误 浏览器离开 或1006超时错误")
			}
			ws.close()
			return
		}
		//ping心跳 不处理通知
		if string(data) == "ping" {
			continue
		}
		msg := &WsMessage{
			MessageType: msgType,
			Data:        data,
		}
		//放入信道中
		select {
		case ws.InChan <- msg:
			fmt.Println("读取的消息放入了inChan", string(msg.Data))
		case <-ws.CloseChan:
			fmt.Println("放入信道时in接收到退出信号")
			return
		}
	}

}

func (ws *WsSession) WsWriteLoop() {
	for {
		select {
		case msg := <-ws.OutChan:
			fmt.Println("获取了outChan", string(msg.Data))
			err := ws.Conn.WriteMessage(msg.MessageType, msg.Data)
			if err != nil {
				fmt.Println("发送out数据给客户端错误", err)
				ws.close()
				return
			}
		case <-ws.CloseChan:
			fmt.Println("WsWriteLoop接收到退出信号")
			return
		}
	}

}

// 关闭连接
func (ws *WsSession) close() {
	fmt.Println("关闭连接被调用了")
	ws.Conn.Close()
	ws.mu.Lock()
	defer ws.mu.Unlock()
	if ws.isClosed == false {
		ws.isClosed = true
		// 删除这个连接的用户
		delete(WsSessions, ws.ID)
		close(ws.CloseChan)
	}
	fmt.Println("还剩下链接:", WsSessions)
}

func (ws *WsSession) readInChan() (msg *WsMessage, err error) {
	select {
	case msg := <-ws.InChan:
		fmt.Println("接收到inChan信息~", string(msg.Data))
		// 给客户端发送消息
		// ws.SendToClient(msg)
		return msg, nil
	case <-ws.CloseChan:
		fmt.Println("获取inChan时接收到退出信号")
		return nil, errorx.NewWithMsg(0, "获取inChan时接收到退出信号")
	}
}

func (ws *WsSession) writeOutChan(msgType int, data []byte) (err error) {
	select {
	case ws.OutChan <- &WsMessage{
		MessageType: msgType,
		Data:        data,
	}:
		return nil
	case <-ws.CloseChan:
		fmt.Println("写入outChan时接收到退出信号")
		return errorx.NewWithMsg(0, "写入outChan时接收到退出信号")
	}
}

// 提供方法下发消息
func SendMessage(id string, msg string) error {
	get, ok := WsSessions[id]
	if ok {
		var data schema.HandleMacheineRes
		data.Flag = 1
		data.Data.Type = msg
		data.Data.Time = time.Now().Unix()
		jsonBytes, err := json.Marshal(data)
		// err := get.writeOutChan(websocket.TextMessage, []byte(str))
		err = get.writeOutChan(websocket.TextMessage, jsonBytes)
		if err != nil {
			return err
		}
	} else {
		return errorx.NewWithMsg(0, "no ws session")
	}
	return nil
}

func Handler(c *gin.Context, id string, dealer Dealer) {
	conn, err := Upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("upgrade err", err)
		return
	}

	item := &WsSession{
		ID:        id,
		Conn:      conn,
		InChan:    make(chan *WsMessage),
		OutChan:   make(chan *WsMessage),
		CloseChan: make(chan byte),
	}
	//存在链接 通知断开 多系统？
	if v, ok := WsSessions[id]; ok {
		v.close()
	}

	WsSessions[id] = item
	fmt.Println("在线人数", len(WsSessions))
	fmt.Println("都有：", WsSessions)

	go item.WsProcessLoop(dealer)
	go item.WsReadLoop()
	go item.WsWriteLoop()

	fmt.Println("end?")
}

func (ws *WsSession) SendToClient(msg *WsMessage) error {
	if err := ws.Conn.WriteMessage(msg.MessageType, msg.Data); err != nil {
		return err
	}
	return nil
}
