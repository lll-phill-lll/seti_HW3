package game

type Figure interface {
	makeMove(x, y int) error
	isNoFig() bool
}
