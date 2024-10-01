package ws

import (
	"runtime"

	"github.com/m1kkY8/gochat-relay/src/entity"
	"github.com/m1kkY8/gochat-relay/src/util"
)

func (wsManager *WebsocketManager) unregisterClient(client *entity.ClientInfo) {
	delete(wsManager.Clients, client.ClientID)
	delete(wsManager.Rooms[client.Room], client.ClientID)
	delete(wsManager.PublicKeys, client.ClientID)

	util.BroadcastOnlineUsers(wsManager.Clients, &wsManager.Mutex)
	util.SendKeys(wsManager.Clients, wsManager.PublicKeys, &wsManager.Mutex)

	client.Conn.Close()
	runtime.GC()
}
