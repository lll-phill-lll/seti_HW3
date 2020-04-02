package room

import (
	"chess/server/game"
	"fmt"
	"net"
)

type Player struct {
	Conn     net.Conn
	Games    []game.Game
	Login    string
	Password string
	IsAdmin  bool
}

func (p Player) String() string {
	return fmt.Sprintln("login:", p.Login, "isAdmin", p.IsAdmin)
}
