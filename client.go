package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Client model defines attributes of a client object.
type Client struct {
	id     string
	socket net.Conn
	data   chan []byte
}

func (client *Client) receive() {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			client.socket.Close()
			break
		}
		if length > 0 {
			messageParts := strings.Split(string(message), "$$$")
			clientID := messageParts[0]
			actualMessage := messageParts[1]
			fmt.Println("CLIENT:", clientID, "RECEIVED: "+string(actualMessage))
		}
	}
}

// StartClientMode method is used for starting a single client.
func StartClientMode() {

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "12345"
	}
	hosts, ok := os.LookupEnv("HOSTS")
	if !ok {
		hosts = "localhost"
	}

	connection, err := net.Dial("tcp", hosts+":"+port)
	if err != nil {
		fmt.Println(err)
	}

	clientID := connection.LocalAddr().String()
	client := &Client{id: clientID, socket: connection}
	fmt.Println("Starting client with ID -", clientID)

	go client.receive()
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		if strings.TrimRight(message, "\n") != "" {
			connection.Write([]byte(client.id + "$$$" + strings.TrimRight(message, "\n")))
		}
	}
}
