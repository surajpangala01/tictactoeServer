package room

import (
	// "encoding/json"
	// "encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type playload struct {
	Oper int         `json:"op"`
	Data interface{} `json:"data"`
}

type Player struct {
	Conn *websocket.Conn
	Mark string
}

var board [3][3]string

func initBoard() {
	for i := 1; i < 3; i++ {
		for j := 1; j < 3; j++ {
			board[i][j] = " "
		}
	}
}
func CreateRoom(player1 *Player, player2 *Player) {
	createdRoom := playload{
		Oper: 1,
		Data: "x",
	}
	// playload.Data =
	sendData(player1, &createdRoom)
	createdRoom.Data = "o"
	sendData(player2, &createdRoom)
	initBoard()
	go ListenToMsg(player1)
	go ListenToMsg(player2)
}

func ListenToMsg(player *Player) {
	var data playload
	err := player.Conn.ReadJSON(&data)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	// var data playload
	// err = json.Unmarshal(dataBytes, &data)
	// if err != nil {
	// 	fmt.Println("Error decoding JSON data")

	// }

	fmt.Println("Recieved Message: %v\n", data)

}

func sendData(player *Player, playload interface{}) bool {

	err := player.Conn.WriteJSON(playload)
	if err != nil {
		fmt.Println("Some error in sending data", playload, err)
		return false
	}
	return true
}
