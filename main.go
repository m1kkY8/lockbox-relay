package main

import (
	"encoding/json"
	"fmt"
	"log"
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
	Username string `json:"username"`
	Message  string `json:"message"`
	To       string `json:"to"`
}

type WebsocketManager struct {
	Clients    map[*ClientInfo]bool
	Broadcast  chan interface{}
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
		Broadcast:  make(chan interface{}),
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
			fmt.Printf("%s connected\n", client.ClientID)

			// Ukloni klienta
		case client := <-wsManager.Unregister:
			delete(wsManager.Clients, client)

			// Posalji poruku svim povezanim klijentima
		case message := <-wsManager.Broadcast:
			msg, err := json.Marshal(message)
			if err != nil {
				log.Println("error marshaling message")
				continue
			}
			fmt.Println()

			for client := range wsManager.Clients {
				if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					fmt.Println("Error writing to websocket")
					client.Conn.Close()
					delete(wsManager.Clients, client)
				}
			}

		}
	}
}

// Ovo trenutno ne radi nista
func (wsManager *WebsocketManager) Shutdown() {
	for client := range wsManager.Clients {
		client.Conn.Close()
	}
}

// Povezivanje na endpoint servera
func EndpointHandler(wsManager *WebsocketManager, ctx *gin.Context) {
	// Upgrade connection from http to ws
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("error upgrading")
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
		// DEBUG
		fmt.Println(string(message))

		// Ovo nije u funkciji
		// wsManager.Broadcast <- message
	}
}

func main() {
	router := gin.Default()
	wsManager := NewWebsocketManager()

	go wsManager.Start()

	router.GET("/ws", func(ctx *gin.Context) {
		EndpointHandler(wsManager, ctx)
	})

	router.Run(":42069")
}
