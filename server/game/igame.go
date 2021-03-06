package game

type Game interface {
	MakeMove(x1, y1, x2, y2 int) error
	MarshalJSON() ([]byte, error)
	GetID() int
	GetOrder() int
}
