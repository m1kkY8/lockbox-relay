package commands

import (
	"github.com/m1kkY8/lockbox-relay/src/entity"
	"github.com/m1kkY8/lockbox-relay/src/ws"
)

func leave(wsManager *ws.WebsocketManager, client *entity.ClientInfo) {
	// Remove client from current Room
	delete(wsManager.Rooms[client.Room], client.ClientID)

	client.Room = ""
	wsManager.Register <- client
}
