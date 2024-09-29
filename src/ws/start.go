package ws

func (wsManager *WebsocketManager) Start() {
	for {
		select {

		case client := <-wsManager.Register:
			wsManager.registerClient(client)

		case client := <-wsManager.Unregister:
			wsManager.unregisterClient(client)

		case message := <-wsManager.Broadcast:
			err := wsManager.broadcast(message)
			if err != nil {
				continue
			}
		}
	}
}
