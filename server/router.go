package server

import (
	gcc "GoChitChat/internal/handlers"
	"GoChitChat/internal/handlers/chat"
	"GoChitChat/internal/handlers/websocket"
	websocket2 "GoChitChat/internal/models/websocket"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {

	router := gin.Default()

	hub := websocket2.NewHub()
	go hub.Run()

	// load HTML templates
	router.LoadHTMLGlob("./templates/*")

	router.Static("/static", "./static")
	router.GET("/", gcc.IndexHandler)
	router.GET("/chatroom", chat.ChatroomHandler)
	router.GET("/ws", func(context *gin.Context) {
		websocket.WsHandler(context.Writer, context.Request, hub)

	})

	return router
}
