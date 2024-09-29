package ws

import (
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

type Room struct {
	Name    string
	Message []byte
}

type WebsocketManager struct {
	Clients    map[string]*entity.ClientInfo
	Rooms      map[string]map[string]*entity.ClientInfo
	Broadcast  chan *Room
	Register   chan *entity.ClientInfo
	Unregister chan *entity.ClientInfo
	Mutex      sync.Mutex
}

func NewWebsocketManager() *WebsocketManager {
	return &WebsocketManager{
		Clients:    make(map[string]*entity.ClientInfo),
		Rooms:      make(map[string]map[string]*entity.ClientInfo),
		Broadcast:  make(chan *Room),
		Register:   make(chan *entity.ClientInfo),
		Unregister: make(chan *entity.ClientInfo),
	}
}
