package util

import (
	"log"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/m1kkY8/gochat-relay/src/entity"
	"github.com/m1kkY8/gochat-relay/src/ws"
	"github.com/vmihailenco/msgpack/v5"
)

func getAllUsers(wsManager *ws.WebsocketManager) []byte {
	var usernames []string

	for _, client := range wsManager.Clients {
		coloredName := client.Color + ":" + client.Username
		usernames = append(usernames, coloredName)
	}

	userString := strings.Join(usernames, " ")

	userlist := entity.Message{
		Type:      entity.ServerMessage,
		Author:    "",
		Content:   userString,
		Timestamp: "",
		To:        "",
	}

	encodedMessage, err := msgpack.Marshal(userlist)
	if err != nil {
		return nil
	}

	return encodedMessage
}

func BroadcastOnlineUsers(wsManager *ws.WebsocketManager) {
	for _, client := range wsManager.Clients {

		wsManager.Mutex.Lock()
		names := getAllUsers(wsManager)

		if client.Conn == nil {
			continue
		}

		err := client.Conn.WriteMessage(websocket.BinaryMessage, names)
		if err != nil {
			log.Println("Error writing to websocket")
			client.Conn.Close()
			delete(wsManager.Clients, client.ClientID)
		}
		wsManager.Mutex.Unlock()
	}
}
