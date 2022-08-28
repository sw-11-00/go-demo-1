package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HelloPage(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "welcome",
	})
}
