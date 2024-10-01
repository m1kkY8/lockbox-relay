package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/m1kkY8/gochat-relay/src/entity"
	"github.com/m1kkY8/gochat-relay/src/ws"
)

func join(wsManager *ws.WebsocketManager, client *entity.ClientInfo, decodedMessage *entity.Message) error {
	parts := strings.Split(decodedMessage.Content, " ")

	// Ensure there's a room name provided after /join
	if len(parts) < 2 || parts[1] == "" {
		log.Println("Room name missing in join command")
		return fmt.Errorf("Room name missing in join command")
	}

	// Register the client to the room
	room := parts[1]

	// Remove client from current Room
	delete(wsManager.Rooms[client.Room], client.ClientID)

	client.Room = room
	wsManager.Register <- client

	return nil
}
