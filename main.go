package main

import (
	"github.com/gin-gonic/gin"
	"github.com/m1kkY8/lockbox-relay/src/handler"
	"github.com/m1kkY8/lockbox-relay/src/ws"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	wsManager := ws.NewWebsocketManager()
	defer wsManager.Shutdown()

	go wsManager.Start()

	router.GET("/chat", handler.EndpointHandler(wsManager))
	router.GET("/health", handler.Health())

	router.Run(":6969")
}
