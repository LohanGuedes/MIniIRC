package server

import (
	"github.com/fatih/color"
	"log"
	"miniIRC/cmd/handler"
	"net"
	"time"
)

// An array/map of the clients...

type Client struct {
	UserName    string
	Connection  net.Conn
	LastMessage time.Time
	StrikeCount int
	BannedAt    time.Time
}

const welcome = `
============================================
|  Welcome to MiniIRC TPC chat! - lguedes  |
|  Version 0.0.1                           |
|  Ecole 42 Rules! <3 Join it! 42.fr       |
|  Special thanks to saiago's <3           |
============================================
`

func Server(messages chan handler.Message) {
	clients := map[string]*Client{}

	banned := false
	for {
		msg := <-messages
		switch msg.Type {
		case handler.ClientConnected:
			if !banned {
				log.Println("Client Successfully Connect")
				_, err := msg.Connection.Write([]byte(color.GreenString(welcome)))
				if err != nil {
					log.Println("Could not Write message to", msg.Connection.LocalAddr().String(), err)
					return
				}
				clients[msg.Connection.LocalAddr().String()] = &Client{
					Connection:  msg.Connection,
					LastMessage: time.Now(),
				}
			}
		case handler.NewMessage:
		}
		// for _, client := range clients {
		// 	fmt.Println(client.Connection.LocalAddr().String())
		// 	_, err := client.Connection.Write([]byte(msg.Content))
		//
		// 	if err != nil {
		//
		// 		log.Println("Failed spreading the message")
		// 	}
		// 	msg.Connection.Write([]byte(msg.Content))
		// }
	}

}
