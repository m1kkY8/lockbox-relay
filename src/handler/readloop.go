package handler

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/m1kkY8/gochat-relay/src/entity"
	"github.com/m1kkY8/gochat-relay/src/ws"
)

type Commands interface {
	Join(wsManager *ws.WebsocketManager, client *entity.ClientInfo, decodedMessage entity.Message) error
	Leave(wsManager *ws.WebsocketManager, client *entity.ClientInfo) error
}

// Read the handshake message from the client
func readHandshake(conn *websocket.Conn) (*entity.Handshake, error) {
	_, bytesHandshake, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	handshake, err := entity.DecodeHandshake(bytesHandshake)
	if err != nil {
		return nil, err
	}

	return &handshake, nil
}

func createClient(handshake *entity.Handshake, conn *websocket.Conn) *entity.ClientInfo {
	return &entity.ClientInfo{
		Conn:     conn,
		ClientID: handshake.ClientId,
		Color:    handshake.Color,
		Username: handshake.Username,
		Pubkey:   handshake.PublicKey,
		Room:     "",
	}
}

func EndpointHandler(wsManager *ws.WebsocketManager, ctx *gin.Context, conn *websocket.Conn) {
	handshake, err := readHandshake(conn)
	if err != nil {
		log.Println("Error reading handshake:", err)
		return
	}

	client := createClient(handshake, conn)
	wsManager.Register <- client

	// Citaj poruke koje klijent salje

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			wsManager.Unregister <- client
			client.Conn.Close()
			break
		}

		// Decode the message
		// TODO: Implement end to end encryption so that server can't read the messages
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
			// Split the message into parts
			prefix := strings.Split(decodedMessage.Content, " ")

			switch prefix[0] {
			case "/join":
				err := Join(wsManager, client, decodedMessage)
				if err != nil {
					break
				}
				break
			case "/leave":
				Leave(wsManager, client)
			}

			// Handle regular chat messages
		}
	}
}
