package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// ClientManager model defines attributes of a manager object.
type ClientManager struct {
	clients    map[*Client]string
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func (manager *ClientManager) start() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = connection.id
			fmt.Println("[NEW CONNECTION] Addr -", connection.id)
		case connection := <-manager.unregister:
			_, ok := manager.clients[connection]
			if ok {
				close(connection.data)
				delete(manager.clients, connection)
				fmt.Println("[CONNECTION CLOSED] Addr -", connection.id)
			}
		case message := <-manager.broadcast:

			messageParts := strings.Split(string(message), "$$$")
			messageClientID := messageParts[0]
			for connection := range manager.clients {
				if messageClientID == connection.id {
					continue
				}
				select {
				case connection.data <- message:
				default:
					close(connection.data)
					delete(manager.clients, connection)
				}
			}
		}
	}
}

func (manager *ClientManager) receive(client *Client) {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			manager.unregister <- client
			client.socket.Close()
			break
		}
		if length > 0 {
			messageParts := strings.Split(string(message), "$$$")
			clientID := messageParts[0]
			actualMessage := messageParts[1]
			fmt.Println("CLIENT:", clientID, "RECEIVED: "+string(actualMessage))
			manager.broadcast <- message
		}
	}
}

func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			}
			client.socket.Write(message)
		}
	}
}

// StartServerMode method is used for starting the server.
func StartServerMode() {

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "12345"
	}
	hosts, ok := os.LookupEnv("HOSTS")
	if !ok {
		hosts = "localhost"
	}

	fmt.Println("Starting server... Accepting connections on -", hosts+":"+port)
	listener, error := net.Listen("tcp", hosts+":"+port)
	if error != nil {
		fmt.Println(error)
	}

	manager := ClientManager{
		clients:    make(map[*Client]string),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	go manager.start()
	for {
		connection, _ := listener.Accept()
		if error != nil {
			fmt.Println(error)
		}
		client := &Client{id: connection.RemoteAddr().String(), socket: connection, data: make(chan []byte)}
		manager.register <- client
		go manager.receive(client)
		go manager.send(client)
	}
}
