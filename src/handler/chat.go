package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/m1kkY8/gochat-relay/src/commands"
	"github.com/m1kkY8/gochat-relay/src/entity"
	"github.com/m1kkY8/gochat-relay/src/ws"
)

func EndpointHandler(wsManager *ws.WebsocketManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		conn, err := ws.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Println("error upgrading")
			return
		}

		handshake, err := entity.NewHandshake(conn)
		if err != nil {
			log.Println("Error reading handshake:", err)
			conn.Close()
			return
		}

		client := entity.NewClient(handshake, conn)
		wsManager.Register <- client

		defer func() {
			wsManager.Unregister <- client
			client.Conn.Close()
		}()

		// Citaj poruke koje klijent salje
		for {
			_, message, err := client.Conn.ReadMessage()
			if err != nil {
				wsManager.Unregister <- client
				client.Conn.Close()
				break
			}

			// Decode the message
			decodedMessage, err := entity.DecodeMessage(message)
			if err != nil {
				log.Println("Error decoding message:", err)
				break
			}

			switch decodedMessage.Type {
			// Handle command message

			case entity.ChatMessage:
				wsManager.Broadcast <- &ws.Room{
					Name:    client.Room,
					Message: message,
				}

			case entity.CommandMessage:
				commands.HandleCommandMessage(wsManager, client, &decodedMessage)
			}
		}
	}
}
