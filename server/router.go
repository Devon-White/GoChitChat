package server

import (
	"GoChitChat/internal/chat/handlers"
	gcc "GoChitChat/internal/handlers"
	"GoChitChat/internal/websocket"
	handlers2 "GoChitChat/internal/websocket/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {

	router := gin.Default()

	hub := websocket.NewHub()
	go hub.Run()

	// load HTML templates
	router.LoadHTMLGlob("./templates/*")

	router.Static("/static", "./static")
	router.GET("/", gcc.IndexHandler)
	router.GET("/chatroom", handlers.ChatroomHandler)
	router.GET("/ws", func(context *gin.Context) {
		handlers2.WsHandler(context.Writer, context.Request, hub)

	})

	return router
}
