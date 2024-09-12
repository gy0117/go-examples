package main

import (
	"github.com/gorilla/websocket"
	"log"
	"os"
)

func main() {
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		for {
			err = c.WriteMessage(websocket.TextMessage, []byte("heartbeat-000"))
			if err != nil {
				log.Println("aaaaa----WriteMessage error:", err)
				return
			}
		}
	}()

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("xxxxx--Read error:", err)
				return
			}
			log.Printf("xxxx-Received: %s", message)
		}
	}()

	for _, msg := range os.Args[1:] {
		err = c.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println("Write error:", err)
			return
		}
	}
}
