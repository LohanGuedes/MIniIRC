package handler

import (
	"log"
	"net"
	"os"
	"strconv"
)

type MessageType int

const (
	ClientDisconnected MessageType = iota + 1
	PrivateMessage
	NewMessage
	QuitMessage
	ClientConnected
)

type Message struct {
	Type       MessageType
	Content    string
	Connection net.Conn
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
	log.Println("User Connected:", c.LocalAddr().String())
	m <- Message{
		Connection: c,
		Type:       ClientConnected,
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
		conn.Write([]byte(content))
		message <- Message{
			Connection: conn,
			Content:    content,
			Type:       NewMessage,
		}
	}
}
