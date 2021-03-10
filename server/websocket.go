package server

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

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
	upgrade, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error(err)
	}

	h.ws.setVal(upgrade)
}

// writeMsg sends a message to the client on the websocket.
func writeMsg(wsConn *websocket.Conn, message []byte) (err error) {
	if err := wsConn.WriteMessage(websocket.TextMessage, message); err != nil {
		return err
	}
	return nil
}

// readMsg reads a message from a client received via a websocket.
func readMsg(ws *websocket.Conn) error {
	_, msg, err := ws.ReadMessage()
	if err != nil {
		//logrus.Errorf("Failed to read message from a websocket connection with error 'v'.", err.Error())  убрать отсюда
		return err
	}

	logrus.Infof("Received new msg from ws : '%v'.", string(msg))
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
