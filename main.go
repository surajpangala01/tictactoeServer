package main

import (
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/gorilla/websocket"
	"net/http"
	"server/room"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	playersQueue := queue.New()
	http.HandleFunc("/", testing)
	http.HandleFunc("/ws", handleWebSocket(playersQueue))

	fmt.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Error starting server", err)
	}
}

func testing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to my Go HTTP server!")

}

func handleWebSocket(playersQueue *queue.Queue) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Upgrade the HTTP connection to a WebSocket connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error upgrading to WebSocket:", err)
			return
		}
		//defer conn.Close()

		fmt.Printf("WebSocket connection established!")

		player := room.Player{
			Conn: conn,
		}
		err = conn.WriteMessage(websocket.TextMessage, []byte("Hello"))
		if err != nil {
			fmt.Println(err)
		}
		playersQueue.Enqueue(player)

		if playersQueue.Len() >= 2 {
			player1, ok := playersQueue.Dequeue().(room.Player)
			if !ok {
				fmt.Println("error in type assertion from Queue for player1")
			}
			player1.Mark = "x"
			player2, ok := playersQueue.Dequeue().(room.Player)
			if !ok {
				fmt.Println("error in type assertion from Queue for player1")
			}
			player2.Mark = "o"
			time.Sleep(1 * time.Second)
			room.CreateRoom(&player1, &player2)

		} else {
			fmt.Println("No one in queue")
		}
	}
}
