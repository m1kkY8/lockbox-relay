package ws

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func (wsManager *WebsocketManager) broadcast(message *Room) error {
	log.Printf("Broadcast channel length: %d, capacity: %d", len(wsManager.Broadcast), cap(wsManager.Broadcast))
	roomName := message.Name

	if roomName == "" {
		return fmt.Errorf("User didnt join any room")
	}
	room, ok := wsManager.Rooms[roomName]
	if !ok {
		return fmt.Errorf("Room does not exist")
	}

	for _, client := range room {
		err := client.Conn.WriteMessage(websocket.BinaryMessage, message.Message)
		if err != nil {
			wsManager.Unregister <- client
			client.Conn.Close()
		}
	}

	return nil
}
