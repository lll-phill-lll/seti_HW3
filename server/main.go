package main

import (
	"chess/server/application"
	"chess/server/room"
)

func getHost() string {
	return "127.0.0.1:8080"
}
func main() {
	r := &application.Rand{}
	player1 := room.Player{Password:"lala", Login:"lolo"}
	player2 := room.Player{Password:"1", Login:"2"}
	serv := &application.Serv{Rand:r, Users: []room.Player{player1, player2}, Rooms: []room.Room{}}
	app := application.App{Rand:r, Serv: serv, Users: []room.Player{player1, player2}, Rooms: []room.Room{}}

	host := getHost()

	app.Start(host)

}


