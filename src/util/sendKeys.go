package util

import (
	"crypto/rsa"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/m1kkY8/gochat-relay/src/entity"
	"github.com/vmihailenco/msgpack/v5"
)

func SendKeys(clients map[string]*entity.ClientInfo, keys map[string]*rsa.PublicKey, mutex *sync.Mutex) {
	for _, client := range clients {

		var keyArray []*rsa.PublicKey
		mutex.Lock()
		for _, key := range keys {
			keyArray = append(keyArray, key)
		}

		if client.Conn == nil {
			continue
		}
		var keyMessage entity.PublicKeys
		keyMessage.Type = 4
		keyMessage.PublicKeys = keyArray
		byteMessage, err := msgpack.Marshal(keyMessage)
		if err != nil {
			continue
		}

		err = client.Conn.WriteMessage(websocket.BinaryMessage, byteMessage)
		if err != nil {
			log.Println("Error writing to websocket")
			client.Conn.Close()
			delete(clients, client.ClientID)

		}
		mutex.Unlock()
	}
}
