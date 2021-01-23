package ws

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type wsT struct {
	upgrader websocket.Upgrader
}

var ws = &wsT{upgrader: websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}}

func (ws *wsT) Conn(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return ws.upgrader.Upgrade(w, r, nil)
}
