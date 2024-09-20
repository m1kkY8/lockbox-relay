package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m1kkY8/gochat-relay/src/handler"
	"github.com/m1kkY8/gochat-relay/src/ws"
)

// Povezivanje na endpoint servera

func main() {
	router := gin.Default()
	wsManager := ws.NewWebsocketManager()

	go wsManager.Start()

	router.GET("/ws", func(ctx *gin.Context) {
		conn, err := ws.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Println("error upgrading")
			return
		}

		go handler.EndpointHandler(wsManager, ctx, conn)
	})

	router.GET("/health", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{"status": "nikola gej"})
	})

	router.Run(":1337")
}
