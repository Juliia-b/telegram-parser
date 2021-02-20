package server

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

// TODO
// разрешить проблему при повторном соединении (?) :
// repeated read on failed websocket connection
// home/julia/go/pkg/mod/github.com/gorilla/websocket@v1.4.2/conn.go:1001 (0xcf7506)
// (*Conn).NextReader: panic("repeated read on failed websocket connection")

/*---------------------------------STRUCTURES----------------------------------------*/

// ws contains all valid websocket connections.
type ws struct {
	connections map[*websocket.Conn]int // Key is websocket connection to client, value is always 0.
	rwMutex     *sync.RWMutex
}

/*-----------------------------------METHODS-----------------------------------------*/

// UpgradeToWs upgrades the HTTP server connection to the WebSocket protocol.
func (h *handler) UpgradeToWs(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this HTTP connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error(err)
	}

	h.ws.setVal(ws)
}

// writeMsg sends a message to the client on the websocket.
func writeMsg(ws *websocket.Conn, message []byte) {
	if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
		logrus.Error(err)
		return
	}
}

// readMsg reads a message from a client received via a websocket.
func readMsg(ws *websocket.Conn) error {
	_, msg, err := ws.ReadMessage()
	if err != nil {
		return err
	}

	msg = msg

	//var ui userIdentity
	//if err := json.Unmarshal(msg, &ui); err != nil {
	//	return nil, err
	//}

	//fmt.Printf("------------------------------------- Websocket readMsg, ui: %#v \n", ui)

	return nil
}

/*-----------------------------------HELPERS-----------------------------------------*/

// setVal adds ws connection into the *ws.connections map.
func (w *ws) setVal(newWS *websocket.Conn) {
	w.rwMutex.Lock()
	w.connections[newWS] = 0
	w.rwMutex.Unlock()
}

// delete removes the connection from the *ws.connections map.
func (w *ws) delete(wsConn *websocket.Conn) {
	w.rwMutex.Lock()
	delete(w.connections, wsConn)
	w.rwMutex.Unlock()
}

// DEPRECATE
// getVal returns a websocket connection to the user by his cookie.
//func (w *ws) getVal(cookie string) (connection *websocket.Conn) {
//	w.rwMutex.RLock()
//	connection, ok := w.connections[cookie]
//	if !ok {
//		// TODO придумать что делать с отсутствием куки
//		logrus.Fatal("DOES NOT HAVE COOKIE")
//	}
//	w.rwMutex.RUnlock()
//
//	return connection
//}
