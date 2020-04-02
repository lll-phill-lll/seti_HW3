package room

import (
	"chess/server/game"
	"net"
)

type Player struct {
	Conn net.Conn
	Games []game.Game
	Login string
	Password string
}
