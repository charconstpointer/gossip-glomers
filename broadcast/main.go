package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Node struct {
	*maelstrom.Node
	mu       sync.RWMutex
	messages []int
}

func NewNode() *Node {
	return &Node{
		Node:     maelstrom.NewNode(),
		messages: []int{},
	}
}

func (n *Node) HandleEcho(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}
	body["type"] = "echo_ok"
	return n.Reply(msg, body)
}

func (n *Node) HandleBroadcast(msg maelstrom.Message) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	var res map[string]any
	if err := json.Unmarshal(msg.Body, &res); err != nil {
		return err
	}
	n.messages = append(n.messages, int(res["message"].(float64)))

	body := map[string]any{}
	body["type"] = "broadcast_ok"
	return n.Reply(msg, body)
}

func (n *Node) HandleRead(msg maelstrom.Message) error {
	n.mu.RLock()
	defer n.mu.RUnlock()
	body := map[string]any{}
	body["type"] = "read_ok"
	body["messages"] = n.messages
	return n.Reply(msg, body)
}

func (n *Node) HandleTopology(msg maelstrom.Message) error {
	body := map[string]any{}
	body["type"] = "topology_ok"
	return n.Reply(msg, body)
}

func main() {
	n := NewNode()
	n.Handle("echo", n.HandleEcho)
	n.Handle("broadcast", n.HandleBroadcast)
	n.Handle("read", n.HandleRead)
	n.Handle("topology", n.HandleTopology)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}