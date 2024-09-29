package ws

import (
	"github.com/m1kkY8/gochat-relay/src/entity"
	"github.com/m1kkY8/gochat-relay/src/util"
)

func (wsManager *WebsocketManager) registerClient(client *entity.ClientInfo) error {
	wsManager.Clients[client.ClientID] = client

	// Create room if it doesn't exist
	if client.Room == "" {
		return nil
	}
	if _, ok := wsManager.Rooms[client.Room]; !ok {
		wsManager.Rooms[client.Room] = make(map[string]*entity.ClientInfo)
	}

	// Add client to Room
	wsManager.Rooms[client.Room][client.ClientID] = client

	util.BroadcastOnlineUsers(wsManager.Clients, &wsManager.Mutex)

	return nil
	// Log and exit if room is empty
}
