package util

import (
	"log"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/m1kkY8/gochat-relay/src/entity"
	"github.com/vmihailenco/msgpack/v5"
)

func getAllUsers(clients map[string]*entity.ClientInfo) []byte {
	var usernames []string

	for _, client := range clients {
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

func BroadcastOnlineUsers(clients map[string]*entity.ClientInfo, mutex *sync.Mutex) {
	for _, client := range clients {

		mutex.Lock()
		names := getAllUsers(clients)

		if client.Conn == nil {
			continue
		}

		err := client.Conn.WriteMessage(websocket.BinaryMessage, names)
		if err != nil {
			log.Println("Error writing to websocket")
			client.Conn.Close()
			delete(clients, client.ClientID)

		}
		mutex.Unlock()
	}
}
