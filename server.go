package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

type ChatServer struct {
	mu       sync.Mutex
	clients  map[int]chan string
	nextID   int
	messages []string
}

type Message struct {
	Name string
	Text string
	ID   int
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		clients: make(map[int]chan string),
		nextID:  1,
	}
}

func (c *ChatServer) RegisterClient(dummy int, id *int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	clientID := c.nextID
	c.nextID++

	c.clients[clientID] = make(chan string, 20)
	*id = clientID

	for otherID, ch := range c.clients {
		if otherID != clientID {
			ch <- fmt.Sprintf("User [%d] joined", clientID)
		}
	}

	return nil
}

func (c *ChatServer) SendMessage(msg Message, reply *string) error {
	if msg.Text == "" {
		return errors.New("empty message")
	}

	formatted := fmt.Sprintf("%s: %s", msg.Name, msg.Text)
	c.mu.Lock()
	c.messages = append(c.messages, formatted)

	for id, ch := range c.clients {
		if id != msg.ID {
			ch <- formatted
		}
	}
	c.mu.Unlock()

	*reply = "ok"
	return nil
}

func (c *ChatServer) ReceiveMessages(clientID int, out *string) error {
	c.mu.Lock()
	ch, exists := c.clients[clientID]
	c.mu.Unlock()

	if !exists {
		*out = ""
		return nil
	}

	msg := <-ch
	*out = msg
	return nil
}

func main() {
	server := NewChatServer()
	rpc.Register(server)

	ln, err := net.Listen("tcp", ":6700")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Server running on port 6700...")

	for {
		conn, _ := ln.Accept()
		go rpc.ServeConn(conn)
	}
}
