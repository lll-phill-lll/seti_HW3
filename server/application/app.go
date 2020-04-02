package application

import (
	"chess/server/room"
	"log"
	"os"
	"sync"
)

var MyLog *log.Logger

func enableLog() {
	f, err := os.OpenFile("log.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	MyLog = log.New(f, "CHESS", log.LstdFlags)
}

type App struct {
	mu    sync.Mutex
	Rand  *Rand
	Serv  *Serv
	Users []room.Player
	Rooms []room.Room
}

func (a *App) Start(host string) {
	enableLog()
	a.Serv.StartServe(host)
}
