package ws

import (
	"github.com/m1kkY8/gochat-relay/src/entity"
	"github.com/m1kkY8/gochat-relay/src/util"
)

func (wsManager *WebsocketManager) unregisterClient(client *entity.ClientInfo) {
	delete(wsManager.Clients, client.ClientID)
	delete(wsManager.Rooms[client.Room], client.ClientID)
	util.BroadcastOnlineUsers(wsManager.Clients, &wsManager.Mutex)
}
