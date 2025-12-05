package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strings"
)

type Message struct {
	Name string
	Text string
	ID   int
}

func main() {
	client, err := rpc.Dial("tcp", "localhost:6700")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer client.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	var id int
	client.Call("ChatServer.RegisterClient", 0, &id)

	fmt.Printf("You joined as User [%d]\n", id)

	go func() {
		for {
			var msg string
			err := client.Call("ChatServer.ReceiveMessages", id, &msg)
			if err == nil && msg != "" {
				fmt.Println("\n> " + msg)
			}
		}
	}()

	for {
		fmt.Print("Enter message (or exit): ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "exit" {
			fmt.Println("Bye!")
			return
		}

		msg := Message{Name: name, Text: text, ID: id}
		var reply string
		client.Call("ChatServer.SendMessage", msg, &reply)
	}
}
