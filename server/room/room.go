package room

import "chess/server/game"

type Room struct {
	WhitePlayer Player
	BlackPlayer Player
	Game game.Game
	Has2Players bool
}
