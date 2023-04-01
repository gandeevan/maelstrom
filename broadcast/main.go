package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type MessageStore struct {
	Messages []float64
}

func (ms *MessageStore) AddMessage(msg float64) {
	ms.Messages = append(ms.Messages, msg)
}

func (ms *MessageStore) GetMessages() []float64 {
	return ms.Messages
}

func main() {
	var messageStore MessageStore

	n := maelstrom.NewNode()

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		err := json.Unmarshal(msg.Body, &body)
		if err != nil {
			log.Printf("error unmarshaling broadcast message: %s", err)
			return err
		}
		messageStore.AddMessage(body["message"].(float64))
		return n.Reply(msg, map[string]string{"type": "broadcast_ok"})
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		body := make(map[string]any)
		body["type"] = "read_ok"
		body["messages"] = messageStore.GetMessages()
		return n.Reply(msg, body)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		body := make(map[string]any)
		body["type"] = "topology_ok"
		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
