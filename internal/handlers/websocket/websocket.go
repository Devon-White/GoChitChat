package websocket

import (
	ws "GoChitChat/internal/models/websocket"
	"GoChitChat/pkg"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsUpgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(w http.ResponseWriter, r *http.Request, hub *ws.Hub) {
	conn, err := wsUpgrade.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade WebSocket", http.StatusInternalServerError)
		return
	}

	clientID := pkg.GenerateUniqueId()
	client := &ws.Client{ID: clientID, Hub: hub, Conn: conn, Send: make(chan []byte, 256)}

	err = conn.WriteJSON(map[string]string{"event": "clientId", "clientId": clientID})
	if err != nil {
		// Handle the error
		log.Println(err)
	}

	hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
