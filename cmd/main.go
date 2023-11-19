package main

import (
	"fmt"
	"log"
	"miniIRC/cmd/handler"
	"miniIRC/cmd/server"
	"net"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println(color.GreenString("Hello, MiniIRC!"))
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ln, err := net.Listen("tcp", "127.0.0.1:"+os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Server Crashed: %s\n", err)
	}
	log.Println("TCP server started listening to port :", os.Getenv("PORT"))

	messages := make(chan handler.Message)
	go server.Server(messages)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to read from %s", conn.LocalAddr().String())
		}
		handler.UserConnected(conn, messages)
		go handler.GetMessages(conn, messages)
	}
}
