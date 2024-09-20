package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/m1kkY8/gochat-relay/src/entity"
	"github.com/m1kkY8/gochat-relay/src/ws"
)

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
