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
}

func NewClient(handshake Handshake, conn *websocket.Conn) *ClientInfo {
	return &ClientInfo{
		Conn:     conn,
		ClientID: handshake.ClientId,
		Color:    handshake.Color,
		Username: handshake.Username,
		Pubkey:   handshake.PublicKey,
		Room:     "",
	}
}
