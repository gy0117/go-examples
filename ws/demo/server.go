package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func echo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	//go func() {
	//	err2 := conn.WriteMessage(websocket.TextMessage, []byte("0xAAC"))
	//	if err != nil {
	//		log.Println("Write error:", err2)
	//	}
	//}()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("kkk---Read error:", err)
			break
		}

		log.Printf("kkk--Received: %s", message)

		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", echo)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
