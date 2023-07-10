package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func IndexHandler(context *gin.Context) {
	context.HTML(http.StatusOK, "index.html", gin.H{})
}
