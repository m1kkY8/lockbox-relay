package entity

import "github.com/gorilla/websocket"

type ClientInfo struct {
	Conn     *websocket.Conn
	ClientID string
	Color    string
	Username string
	Pubkey   string
	Room     string
}
