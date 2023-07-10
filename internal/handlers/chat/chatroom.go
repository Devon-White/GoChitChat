package chat

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ChatroomHandler(context *gin.Context) {
	context.HTML(http.StatusOK, "chatroom.html", gin.H{})
}
