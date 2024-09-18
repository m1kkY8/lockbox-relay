package entity

import "github.com/vmihailenco/msgpack/v5"

type Handshake struct {
	Username  string `msgpack:"username"`
	Color     string `msgpack:"color"`
	ClientId  string `msgpack:"client_id"`
	PublicKey string `msgpack:"pubkey"`
}

func DecodeHandshake(encodedHandshake []byte) (Handshake, error) {
	var handshake Handshake
	err := msgpack.Unmarshal(encodedHandshake, &handshake)
	if err != nil {
		return Handshake{}, err
	}

	return handshake, nil
}
