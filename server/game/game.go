package game

const (
	White int = 0
	Black int = 1
)

type game struct {
	Id    int
	Order int
	Field [][]Figure
}

func (g *game) MakeMove(x1, y1, x2, y2 int) error {
	err := g.Field[x1][y1].makeMove(x2, y2)
	if err != nil {
		return err
	}

	g.Field[x2][y2] = g.Field[x1][y1]
	g.Field[x1][y1] = newNoFig(x1, y1)

	return nil
}

func newGame(id int) Game {
	return &game{
		Id:    id,
		Order: 0,
		Field: [][]Figure{
			{newRook(0, 0, White), newKnight(0, 1, White), newBishop(0, 2, White), newQueen(0, 3, White), newKing(0, 4, White), newBishop(0, 5, White), newKnight(0, 6, White), newRook(0, 7, White)},
			{newPawn(1, 0, White), newPawn(1, 1, White), newPawn(1, 2, White), newPawn(1, 3, White), newPawn(1, 4, White), newPawn(1, 5, White), newPawn(1, 6, White), newPawn(1, 7, White)},
			{newNoFig(2, 0), newNoFig(2, 1), newNoFig(2, 2), newNoFig(2, 3), newNoFig(2, 4), newNoFig(2, 5), newNoFig(2, 6), newNoFig(2, 7)},
			{newNoFig(3, 0), newNoFig(3, 1), newNoFig(3, 2), newNoFig(3, 3), newNoFig(3, 4), newNoFig(3, 5), newNoFig(3, 6), newNoFig(3, 7)},
			{newNoFig(4, 0), newNoFig(4, 1), newNoFig(4, 2), newNoFig(4, 3), newNoFig(4, 4), newNoFig(4, 5), newNoFig(4, 6), newNoFig(4, 7)},
			{newNoFig(5, 0), newNoFig(5, 1), newNoFig(5, 2), newNoFig(5, 3), newNoFig(5, 4), newNoFig(5, 5), newNoFig(5, 6), newNoFig(5, 7)},
			{newPawn(6, 0, Black), newPawn(6, 1, Black), newPawn(6, 2, Black), newPawn(6, 3, Black), newPawn(6, 4, Black), newPawn(6, 5, Black), newPawn(6, 6, Black), newPawn(1, 7, White)},
			{newRook(7, 0, Black), newKnight(7, 1, Black), newBishop(7, 2, Black), newKing(7, 4, Black), newQueen(7, 3, Black), newBishop(7, 5, Black), newKnight(7, 6, Black), newRook(0, 7, White)},
		},
	}
}
