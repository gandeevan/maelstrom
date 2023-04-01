package main

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	log "github.com/sirupsen/logrus"
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

func sendReliableMessage(n *maelstrom.Node, neighbor string, body map[string]any) {
	err := n.Send(neighbor, body)
	if err == nil {
		log.Info("error sending broadcast message to %s: %s", neighbor, err)
	}
}

func deserializeNeighborList(neigbhors []interface{}) []string {
	var result []string
	for _, neighbor := range neigbhors {
		result = append(result, neighbor.(string))
	}
	return result
}

func main() {
	var messageStore MessageStore

	n := maelstrom.NewNode()
	var neigbhors []string

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		err := json.Unmarshal(msg.Body, &body)
		if err != nil {
			log.Printf("error unmarshaling broadcast message: %s", err)
			return err
		}

		messageStore.AddMessage(body["message"].(float64))
		for _, neighbor := range neigbhors {
			go sendReliableMessage(n, neighbor, body)
		}
		return n.Reply(msg, map[string]string{"type": "broadcast_ok"})
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		body := make(map[string]any)
		body["type"] = "read_ok"
		body["messages"] = messageStore.GetMessages()
		return n.Reply(msg, body)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		log.Infof("received topology message: %s", msg.Body)
		requestBody := make(map[string]any)
		err := json.Unmarshal(msg.Body, &requestBody)
		if err != nil {
			log.Printf("error unmarshaling topology message: %s", err)
			return err
		}

		neigbhors = deserializeNeighborList(requestBody["topology"].(map[string]interface{})[n.ID()].([]interface{}))

		responseBody := make(map[string]any)
		responseBody["type"] = "topology_ok"
		return n.Reply(msg, responseBody)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
