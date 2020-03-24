package game

import (
	"errors"
	"fmt"
)

func newPawn(x, y, col int) Figure {
	return &pawn{currentPosX: x, currentPosY: y, color: col}
}

func newRook(x, y, col int) Figure {
	return &rook{currentPosX: x, currentPosY: y, color: col}
}

func newKnight(x, y, col int) Figure {
	return &knight{currentPosX: x, currentPosY: y, color: col}
}

func newBishop(x, y, col int) Figure {
	return &bishop{currentPosX: x, currentPosY: y, color: col}
}

func newKing(x, y, col int) Figure {
	return &king{currentPosX: x, currentPosY: y, color: col}
}

func newQueen(x, y, col int) Figure {
	return &queen{currentPosX: x, currentPosY: y, color: col}
}

func newNoFig(x, y int) Figure {
	return &noFig{currentPosX: x, currentPosY: y}
}

type pawn struct {
	currentPosX int
	currentPosY int
	color       int
}

func (p *pawn) isNoFig() bool {
	return false
}

func (p *pawn) makeMove(x, y int) error {
	// TODO: add move validation
	p.currentPosX = x
	p.currentPosY = y

	return nil
}

type rook struct {
	currentPosX int
	currentPosY int
	color       int
}

func (r *rook) isNoFig() bool {
	return false
}

func (r *rook) makeMove(x, y int) error {
	// TODO: add move validation
	r.currentPosX = x
	r.currentPosY = y

	return nil
}

type knight struct {
	currentPosX int
	currentPosY int
	color       int
}

func (kn *knight) isNoFig() bool {
	return false
}

func (kn *knight) makeMove(x, y int) error {
	// TODO: add move validation
	kn.currentPosX = x
	kn.currentPosY = y

	return nil
}

type bishop struct {
	currentPosX int
	currentPosY int
	color       int
}

func (b *bishop) isNoFig() bool {
	return false
}

func (b *bishop) makeMove(x, y int) error {
	// TODO: add move validation
	b.currentPosX = x
	b.currentPosY = y

	return nil
}

type king struct {
	currentPosX int
	currentPosY int
	color       int
}

func (kg *king) isNoFig() bool {
	return false
}

func (kg *king) makeMove(x, y int) error {
	// TODO: add move validation
	kg.currentPosX = x
	kg.currentPosY = y

	return nil
}

type queen struct {
	currentPosX int
	currentPosY int
	color       int
}

func (q *queen) isNoFig() bool {
	return false
}

func (q *queen) makeMove(x, y int) error {
	// TODO: add move validation
	q.currentPosX = x
	q.currentPosY = y

	return nil
}

type noFig struct {
	currentPosX int
	currentPosY int
}

func (nf *noFig) isNoFig() bool {
	return true
}

func (nf *noFig) makeMove(x, y int) error {
	return errors.New(fmt.Sprintf("no figure at (%d, %d)", x, y))
}
