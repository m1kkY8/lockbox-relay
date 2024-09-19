package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/m1kkY8/gochat-relay/src/entity"
	"github.com/m1kkY8/gochat-relay/src/ws"
)

// Povezivanje na endpoint servera
func EndpointHandler(wsManager *ws.WebsocketManager, ctx *gin.Context, conn *websocket.Conn) {
	_, bytesHandshake, err := conn.ReadMessage()
	if err != nil {
		log.Println("error reading handshake")
		return
	}

	handshake, err := entity.DecodeHandshake(bytesHandshake)
	if err != nil {
		log.Println("error decoding handshake")
		return
	}

	client := &entity.ClientInfo{
		Conn:     conn,
		ClientID: handshake.ClientId,
		Color:    handshake.Color,
		Username: handshake.Username,
		Pubkey:   handshake.PublicKey,
	}

	wsManager.Register <- client

	// Citaj poruke koje klijent salje
	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			wsManager.Unregister <- client
			client.Conn.Close()
			break
		}
		wsManager.Broadcast <- message
	}
}

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

		go EndpointHandler(wsManager, ctx, conn)
	})

	router.GET("/health", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{"status": "nikola gej"})
	})

	router.Run(":1337")
}
