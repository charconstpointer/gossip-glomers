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
	messages map[int]struct{}
}

func NewNode() *Node {
	return &Node{
		Node:     maelstrom.NewNode(),
		messages: make(map[int]struct{}),
	}
}

type broadcastMsg struct {
	Type    string `json:"type"`
	Message int    `json:"message"`
}

type topologyMsg struct {
	Topology map[string][]string `json:"topology"`
	Type     string              `json:"type"`
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
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	id := int(body["message"].(float64))
	n.mu.Lock()
	if _, exists := n.messages[id]; exists {
		n.mu.Unlock()
		return nil
	}
	n.messages[id] = struct{}{}
	n.mu.Unlock()

	for _, dst := range n.NodeIDs() {
		if dst == msg.Src || dst == n.ID() {
			continue
		}

		dst := dst
		go func() {
			if err := n.Send(dst, body); err != nil {
				panic(err)
			}
		}()
	}

	return n.Reply(msg, map[string]any{
		"type": "broadcast_ok",
	})
}

func (n *Node) HandleRead(msg maelstrom.Message) error {
	n.mu.RLock()
	messages := make([]int, 0, len(n.messages))
	for k := range n.messages {
		messages = append(messages, k)
	}
	n.mu.RUnlock()

	body := map[string]any{}
	body["type"] = "read_ok"
	body["messages"] = messages
	return n.Reply(msg, body)
}

func (n *Node) HandleTopology(msg maelstrom.Message) error {
	var top topologyMsg
	if err := json.Unmarshal(msg.Body, &top); err != nil {
		return err
	}

	return n.Reply(msg, map[string]any{
		"type": "topology_ok",
	})
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
