package models

type SocketMessage struct {
	Event string      `json:"event"`
	DATA  interface{} `json:"data"`
}
