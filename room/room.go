package room

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// Payload represents the structure of messages exchanged between players.
type Payload struct {
	Oper int         `json:"op"`
	Data interface{} `json:"data"`
}

// Player represents a participant in the game.
type Player struct {
	Conn *websocket.Conn
	Mark string
}

var (
	players map[string]*Player
	board   [3][3]string
)

// initBoard initializes the game board to an empty state.
func initBoard() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			board[i][j] = " "
		}
	}
}

// sendData sends a payload to the player's WebSocket connection.
func (player *Player) sendData(data Payload) bool {
	err := player.Conn.WriteJSON(data)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return false
	}
	return true
}

// CreateRoom sets up a new game room with two players and initializes the game state.
func CreateRoom(player1, player2 *Player) {
	players = make(map[string]*Player)
	players["x"] = player1
	players["o"] = player2

	// Notify both players of their assigned marks
	player1.sendData(Payload{Oper: 1, Data: "x"})
	player2.sendData(Payload{Oper: 1, Data: "o"})

	// Initialize the game board
	initBoard()

	// Start listening for messages from both players
	go listenToMessages(player1)
	go listenToMessages(player2)
}

// listenToMessages listens for messages from a player and processes them.
func listenToMessages(player *Player) {
	for {
		var data Payload
		err := player.Conn.ReadJSON(&data)
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}

		fmt.Printf("Received Message: %+v\n", data)
		parseMessage(data, player.Mark)
	}
}

// sendData sends a payload to the specified player's WebSocket connection.
func sendData(player *Player, payload Payload) bool {
	err := player.Conn.WriteJSON(payload)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return false
	}
	return true
}

// parseMessage handles and processes incoming messages based on their operation type.
func parseMessage(data Payload, mark string) {
	switch data.Oper {
	case 2: // Handle move operation
		moveData, ok := data.Data.(map[string]interface{})
		if !ok {
			fmt.Println("Invalid data format for move")
			return
		}

		row, rowOk := moveData["row"].(float64) // JSON numbers are float64
		col, colOk := moveData["col"].(float64)

		if !rowOk || !colOk {
			fmt.Println("Invalid move coordinates")
			return
		}

		// Check if the cell is already filled
		if board[int(row)][int(col)] != " " {
			fmt.Printf("Cell already filled at [%d, %d]\n", int(row), int(col))
			return
		}

		// Update the board and send the move to the opponent
		board[int(row)][int(col)] = mark
		opponentMark := "o"
		if mark == "o" {
			opponentMark = "x"
		}
		if !players[opponentMark].sendData(data) {
			fmt.Println("data not sent")
		}
	}
}
