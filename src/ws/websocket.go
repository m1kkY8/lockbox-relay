package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/m1kkY8/gochat-relay/src/entity"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:    8096,
	WriteBufferSize:   8096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketManager struct {
	Clients    map[string]*entity.ClientInfo
	Broadcast  chan []byte
	Register   chan *entity.ClientInfo
	Unregister chan *entity.ClientInfo
	Mutex      sync.Mutex
}

func NewWebsocketManager() *WebsocketManager {
	return &WebsocketManager{
		Clients:    make(map[string]*entity.ClientInfo),
		Broadcast:  make(chan []byte),
		Register:   make(chan *entity.ClientInfo),
		Unregister: make(chan *entity.ClientInfo),
	}
}

func (wsManager *WebsocketManager) Start() {
	for {
		select {
		// Dodaj klienta
		case client := <-wsManager.Register:
			wsManager.Clients[client.ClientID] = client

			// Ukloni klienta

		case client := <-wsManager.Unregister:
			delete(wsManager.Clients, client.ClientID)

			// Posalji poruku svim povezanim klijentima
		case message := <-wsManager.Broadcast:
			for _, client := range wsManager.Clients {

				err := client.Conn.WriteMessage(websocket.BinaryMessage, message)
				if err != nil {
					log.Println("Error writing to websocket")
					client.Conn.Close()
					delete(wsManager.Clients, client.ClientID)
				}
			}
		}
	}
}

func (wsManager *WebsocketManager) Shutdown() {
	for _, client := range wsManager.Clients {
		client.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server is shutting down"))
		client.Conn.Close()
	}
}
