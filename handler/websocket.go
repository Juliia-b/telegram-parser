package handler

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

//var upgrader = websocket.Upgrader{
//	ReadBufferSize:  1024,
//	WriteBufferSize: 1024,
//}

//func wbHandler(w http.ResponseWriter, r *http.Request) {
//	conn, err := upgrader.Upgrade(w, r, nil)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	//... Use conn to send and receive messages.
//
//	for {
//		messageType, p, err := conn.ReadMessage()
//		if err != nil {
//			log.Println(err)
//			return
//		}
//		if err := conn.WriteMessage(messageType, p); err != nil {
//			log.Println(err)
//			return
//		}
//	}
//}

// -------------------------------- messenger example --------------------------------

// contains all active connections to the client
var wsConnections []*websocket.Conn

type Client struct {
	wsConn *websocket.Conn
	id     string // id of user in database
}

func (h *handler) WsInit(w http.ResponseWriter, r *http.Request) {

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil) //ws
	if err != nil {
		log.Println(err)
	}

	WriteMsg(ws, []byte("Some string from server"))

	logrus.Info("сообщение отправлено")

	//var user *userIdentity
	////будет ждать ответа от ReadMsg пока не получит id
	////TODO нужен таймаут (около полуминуты или меньше тк всё это время юзер не сможет получить сообщения)
	//user, err = ReadMsg(ws)
	//if err != nil {
	//	panic(err) // TODO тут выкидывает ошибку, если соединение было закрыто => необходимо обработать
	//}
	////TODO обработать ошибку (отправить по ws?)
}

// разрешить проблему при повторном соединении (?) :
// repeated read on failed websocket connection
// home/julia/go/pkg/mod/github.com/gorilla/websocket@v1.4.2/conn.go:1001 (0xcf7506)
// (*Conn).NextReader: panic("repeated read on failed websocket connection")

// TODO должна вернуться ошибка, если невозможно записать данные. CheckChanges должна дождаться соединения
func WriteMsg(ws *websocket.Conn, message []byte) {
	// Отправить сообщение клиенту в виде массива байт
	if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Println(err)
		return
	}
}

func ReadMsg(ws *websocket.Conn) error {
	_, msg, err := ws.ReadMessage()
	if err != nil {
		return err
	}

	msg = msg

	//var ui userIdentity
	//if err := json.Unmarshal(msg, &ui); err != nil {
	//	return nil, err
	//}

	//fmt.Printf("------------------------------------- Websocket ReadMsg, ui: %#v \n", ui)

	return nil
}
