package main

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/m1kkY8/gochat-relay/src/entity"
	"github.com/vmihailenco/msgpack/v5"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    8096,
	WriteBufferSize:   8096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketManager struct {
	Clients    map[string]*ClientInfo
	Broadcast  chan []byte
	Register   chan *ClientInfo
	Unregister chan *ClientInfo
	Mutex      sync.Mutex
}

type ClientInfo struct {
	Conn     *websocket.Conn
	ClientID string
	Color    string
	Username string
	Pubkey   string
}

// Instancira novi Manager
func NewWebsocketManager() *WebsocketManager {
	return &WebsocketManager{
		Clients:    make(map[string]*ClientInfo),
		Broadcast:  make(chan []byte),
		Register:   make(chan *ClientInfo),
		Unregister: make(chan *ClientInfo),
	}
}

// Startuje event loop koji slusa poruke
func (wsManager *WebsocketManager) Start() {
	for {
		select {
		// Dodaj klienta
		case client := <-wsManager.Register:
			wsManager.Clients[client.ClientID] = client
			broadcastOnlineUsers(wsManager)

			// Ukloni klienta

		case client := <-wsManager.Unregister:
			delete(wsManager.Clients, client.ClientID)
			broadcastOnlineUsers(wsManager)

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

// Povezivanje na endpoint servera
func EndpointHandler(wsManager *WebsocketManager, ctx *gin.Context) {
	// Upgrade connection from http to ws
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("error upgrading")
		return
	}

	_, bytesHandshake, err := conn.ReadMessage()
	if err != nil {
		log.Println("error reading handshake")
		return
	}

	handshake, err := entity.DecodeHandshake(bytesHandshake)
	if err != nil {
		log.Println("error decoding handshake")
		return
	}

	client := &ClientInfo{
		Conn:     conn,
		ClientID: handshake.ClientId,
		Color:    handshake.Color,
		Username: handshake.Username,
		Pubkey:   handshake.PublicKey,
	}

	wsManager.Register <- client

	// Citaj poruke koje klijent salje
	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			wsManager.Unregister <- client
			client.Conn.Close()
			break
		}
		wsManager.Broadcast <- message
	}
}

func GetAllUsers(wsManager *WebsocketManager) []byte {
	var usernames []string

	for _, client := range wsManager.Clients {
		coloredName := client.Color + ":" + client.Username
		usernames = append(usernames, coloredName)
	}

	userString := strings.Join(usernames, " ")

	userlist := entity.Message{
		Type:      entity.ServerMessage,
		Author:    "",
		Content:   userString,
		Timestamp: "",
		To:        "",
	}

	encodedMessage, err := msgpack.Marshal(userlist)
	if err != nil {
		return nil
	}

	return encodedMessage
}

func broadcastOnlineUsers(wsManager *WebsocketManager) {
	for _, client := range wsManager.Clients {

		wsManager.Mutex.Lock()
		names := GetAllUsers(wsManager)

		if client.Conn == nil {
			continue
		}

		err := client.Conn.WriteMessage(websocket.BinaryMessage, names)
		if err != nil {
			log.Println("Error writing to websocket")
			client.Conn.Close()
			delete(wsManager.Clients, client.ClientID)
		}
		wsManager.Mutex.Unlock()
	}
}

func main() {
	router := gin.Default()
	wsManager := NewWebsocketManager()

	go wsManager.Start()

	router.GET("/ws", func(ctx *gin.Context) {
		EndpointHandler(wsManager, ctx)
	})

	router.GET("/health", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{"status": "nikola gej"})
	})

	router.Run(":1337")
}
