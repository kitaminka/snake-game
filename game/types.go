package game

import (
	"github.com/gdamore/tcell"
	"time"
)

type Cell struct {
	x, y int
}
type Snake struct {
	head      Cell
	tail      []Cell
	field     *Field
	style     tcell.Style
	direction struct {
		x, y int
	}
	delay time.Duration
}
type Apple struct {
	cell  Cell
	style tcell.Style
	field *Field
}
type Score struct {
	x, y  int
	value int
	field *Field
}
type Field struct {
	x, y          int
	width, height int
	style         tcell.Style
	cells         []Cell
	screen        *tcell.Screen
	snake         *Snake
	score         *Score
	apples        []Apple
}
