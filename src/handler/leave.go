package handler

import (
	"github.com/m1kkY8/gochat-relay/src/entity"
	"github.com/m1kkY8/gochat-relay/src/ws"
)

func Leave(wsManager *ws.WebsocketManager, client *entity.ClientInfo) {
	// Remove client from current Room
	delete(wsManager.Rooms[client.Room], client.ClientID)

	client.Room = ""
	wsManager.Register <- client
}
