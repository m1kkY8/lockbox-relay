package ws

import (
	"github.com/m1kkY8/lockbox-relay/src/entity"
	"github.com/m1kkY8/lockbox-relay/src/util"
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
	wsManager.PublicKeys[client.ClientID] = client.Pubkey

	util.BroadcastOnlineUsers(wsManager.Clients, &wsManager.Mutex)
	util.SendKeys(wsManager.Clients, wsManager.PublicKeys, &wsManager.Mutex)
	return nil
	// Log and exit if room is empty
}
