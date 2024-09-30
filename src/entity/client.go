package entity

import (
	"crypto/rsa"

	"github.com/gorilla/websocket"
)

type ClientInfo struct {
	Conn     *websocket.Conn
	ClientID string
	Color    string
	Username string
	Pubkey   *rsa.PublicKey
	Room     string
	Amogus   string
}
