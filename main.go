package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Generate unique id for new client connection
func generateClientID() string {
	return uuid.New().String()
}

type message struct {
	Timestamp string `json:"timestamp"`
	Username  string `json:"username"`
	Message   string `json:"message"`
	To        string `json:"to"`
}

type WebsocketManager struct {
	Clients    map[*ClientInfo]bool
	Broadcast  chan []byte
	Register   chan *ClientInfo
	Unregister chan *ClientInfo
}

type ClientInfo struct {
	Conn     *websocket.Conn
	ClientID string
}

// Instancira novi Manager
func NewWebsocketManager() *WebsocketManager {
	return &WebsocketManager{
		Clients:    make(map[*ClientInfo]bool),
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
			wsManager.Clients[client] = true

			// Ukloni klienta

		case client := <-wsManager.Unregister:
			delete(wsManager.Clients, client)

			// Posalji poruku svim povezanim klijentima
		case message := <-wsManager.Broadcast:
			for client := range wsManager.Clients {
				if err := client.Conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
					fmt.Println("Error writing to websocket")
					client.Conn.Close()
					delete(wsManager.Clients, client)
				}
			}

		}
	}
}

func (wsManager *WebsocketManager) Shutdown() {
	for client := range wsManager.Clients {
		client.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server is shutting down"))
		client.Conn.Close()
	}
}

// Povezivanje na endpoint servera
func EndpointHandler(wsManager *WebsocketManager, ctx *gin.Context) {
	// Upgrade connection from http to ws
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("error upgrading")
		return
	}

	client := &ClientInfo{
		Conn: conn,
		// TODO: Sacuvati ID i poslati ga klientu radi logovanja
		ClientID: generateClientID(),
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

func main() {
	router := gin.Default()
	wsManager := NewWebsocketManager()

	go wsManager.Start()
	// defer wsManager.Shutdown()

	router.GET("/ws", func(ctx *gin.Context) {
		EndpointHandler(wsManager, ctx)
	})

	router.GET("/health", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.Run(":1337")
}
