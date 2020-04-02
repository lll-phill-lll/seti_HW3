package application

import (
	"chess/server/room"
	"sync"
)

type App struct {
	mu sync.Mutex
	Rand *Rand
	Serv *Serv
	Users []room.Player
	Rooms []room.Room
}

func (a *App) Start(host string) {
	a.Serv.StartServe(host)
}


