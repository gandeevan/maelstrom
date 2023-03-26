package main

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
	n.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any
		err := json.Unmarshal(msg.Body, &body)
		if err != nil {
			log.Printf("error unmarshaling echo message: %s", err)
			return err
		}

		uuid, err := uuid.NewRandom()
		if err != nil {
			log.Printf("error generating uuid: %s", err)
			return err
		}

		body["type"] = "generate_ok"
		body["id"] = uuid.String()
		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
