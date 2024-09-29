package entity

import "github.com/vmihailenco/msgpack/v5"

var (
	ServerMessage  = 1
	ChatMessage    = 2
	CommandMessage = 3
)

type Message struct {
	Type      int    `msgpack:"type"`
	Author    string `msgpack:"author"`
	Content   string `msgpack:"content"`
	Room      string `msgpack:"room"`
	To        string `msgpack:"to"`
	Timestamp string `msgpack:"timestamp"`
	Color     string `msgpack:"color"`
}

func DecodeMessage(byteMessage []byte) (Message, error) {
	var decodedMessage Message

	err := msgpack.Unmarshal(byteMessage, &decodedMessage)
	if err != nil {
		return Message{}, err
	}

	return decodedMessage, nil
}
