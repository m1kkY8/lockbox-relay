package commands

import (
	"log"
	"strings"

	"github.com/m1kkY8/lockbox-relay/src/entity"
	"github.com/m1kkY8/lockbox-relay/src/ws"
)

func HandleCommandMessage(wsManager *ws.WebsocketManager, client *entity.ClientInfo, decodedMessage *entity.Message) {
	prefix := strings.Split(decodedMessage.Content, " ")

	switch prefix[0] {
	case "/join":
		err := join(wsManager, client, decodedMessage)
		if err != nil {
			log.Println("Error joining room:", err)
		}
	case "/leave":
		leave(wsManager, client)
	default:
		log.Println("Unknown command:", prefix[0])
	}
}
