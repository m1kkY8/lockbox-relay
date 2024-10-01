package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{"status": "alive"})
	}
}
