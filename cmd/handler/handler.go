package handler

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
)

type MessageType int

const (
	ClientDisconnected MessageType = iota + 1
	PrivateMessage
	NewMessage
	ClientConnected
	QuitMessage
)

const welcome = `
============================================
|  Welcome to MiniIRC TPC chat! - lguedes  |
|  Version 0.0.1                           |
|  Ecole 42 Rules! <3 Join it! 42.fr       |
|  Special thanks to saiago's <3           |
============================================

`

// TODO: make this a file more golang idiomatic

type Message struct {
	Type       MessageType
	Content    string
	Connection net.Conn
}

type Client struct {
	UserName    string
	Connection  net.Conn
	LastMessage time.Time
	StrikeCount int
	BannedAt    time.Time
}

func UserDisconnected(c net.Conn, m chan Message, err error) {
	log.Println("1 Failed to read from connection", c.LocalAddr().String(), err)
	c.Close()
	m <- Message{
		Content:    "User Disconnected from the Server due to an Error",
		Connection: c,
		Type:       ClientDisconnected,
	}
}

func UserConnected(c net.Conn, m chan Message) {
	log.Println("User Connected:", c.RemoteAddr().String())
	m <- Message{
		Connection: c,
		Type:       ClientConnected,
	}
}

func HandleNewConnection(msg Message) (*Client, error) {
	n, err := msg.Connection.Write([]byte(color.GreenString(welcome)))
	if err != nil || n < len(color.GreenString(welcome)) {
		return &Client{}, errors.New(fmt.Sprintln("Could not Write message to", msg.Connection.RemoteAddr().String(), err))
	}
	log.Println("Client Successfully Connected")
	return &Client{
		Connection:  msg.Connection,
		LastMessage: time.Now(),
	}, nil
}

func HandleNewMessage(msg Message, clients map[string]*Client) {
	authorAddr := msg.Connection.RemoteAddr()
	// author := clients[authorAddr.String()]
	log.Printf("New message from %v %s\n", msg.Connection.RemoteAddr().String(), msg.Content)
	for _, client := range clients {
		fmt.Println(client.UserName)
		if client.Connection.RemoteAddr().String() != authorAddr.String() {
			_, err := client.Connection.Write([]byte(client.UserName + ": " + msg.Content))
			if err != nil {
				log.Println("Failed spreading the message")
			}
		}
	}
}

func GetMessages(conn net.Conn, message chan Message) {
	bufferSize, err := strconv.Atoi(os.Getenv("MESSAGE_MAX_LEN"))
	if err != nil {
		log.Println("Failed to get MESSAGE_MAX_LEN from .env Falling back to 512bytes", err)
		bufferSize = 512
	}
	messageBuffer := make([]byte, bufferSize)
	for {
		n, err := conn.Read(messageBuffer)

		if err != nil {
			UserDisconnected(conn, message, err)
			return
		}

		// TODO Create a split in order to check for user commands
		content := string(messageBuffer[0:n])
		message <- Message{
			Connection: conn,
			Content:    content,
			Type:       NewMessage,
		}
	}
}
