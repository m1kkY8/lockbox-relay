package entity

import (
	"crypto/rsa"

	"github.com/gorilla/websocket"
	"github.com/vmihailenco/msgpack/v5"
)

type Handshake struct {
	Username  string         `msgpack:"username"`
	Color     string         `msgpack:"color"`
	ClientId  string         `msgpack:"client_id"`
	PublicKey *rsa.PublicKey `msgpack:"pubkey"`
}

func decodeHandshake(encodedHandshake []byte) (Handshake, error) {
	var handshake Handshake
	err := msgpack.Unmarshal(encodedHandshake, &handshake)
	if err != nil {
		return Handshake{}, err
	}

	return handshake, nil
}

// Read the handshake message from the client
func NewHandshake(conn *websocket.Conn) (Handshake, error) {
	_, bytesHandshake, err := conn.ReadMessage()
	if err != nil {
		return Handshake{}, err
	}

	handshake, err := decodeHandshake(bytesHandshake)
	if err != nil {
		return Handshake{}, err
	}

	return handshake, nil
}
