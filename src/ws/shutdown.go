package ws

import "github.com/gorilla/websocket"

func (wsManager *WebsocketManager) Shutdown() {
	for _, client := range wsManager.Clients {
		client.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server is shutting down"))
		client.Conn.Close()
	}
}
