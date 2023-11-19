package server

import (
	"log"
	"miniIRC/cmd/handler"
)

func Server(messages chan handler.Message) {
	clients := map[string]*handler.Client{} // Holds all the connections

	for {
		msg := <-messages
		switch msg.Type {
		case handler.ClientConnected:
			client, err := handler.HandleNewConnection(msg)
			if err != nil {
				log.Println("Failed Getting user connection:", err)
				continue
			}
			clients[msg.Connection.RemoteAddr().String()] = client
		case handler.NewMessage:
			handler.HandleNewMessage(msg, clients)
		}
	}
}
