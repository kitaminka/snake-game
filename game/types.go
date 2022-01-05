package game

import (
	"github.com/gdamore/tcell"
	"time"
)

type Cell struct {
	x, y int
}
type Snake struct {
	Head      Cell
	Tail      []Cell
	Field     *Field
	Style     tcell.Style
	Direction struct {
		X, Y int
	}
	Delay time.Duration
}
type Apple struct {
	Cell  Cell
	Style tcell.Style
	Field *Field
}
type Score struct {
	X, Y  int
	Value int
	Field *Field
}
type Field struct {
	X, Y          int
	Width, Height int
	Style         tcell.Style
	Cells         []Cell
	Screen        *tcell.Screen
	Snake         *Snake
	Score         *Score
	Apples        []Apple
}
