package game

import (
	"github.com/gdamore/tcell"
	"time"
)

const (
	FieldSymbol  = ' '
	SnakeSymbol  = '▇'
	BorderSymbol = '▇'
	AppleSymbol  = '●'
)

type Game struct {
	Screen        *tcell.Screen
	Field         *Field
	Score         *Score
	Width, Height int
	Ended         bool
	Styles
	Configuration
}
type Styles struct {
	DefaultStyle, SnakeStyle, AppleStyle tcell.Style
}
type Configuration struct {
	StartDelay, MinDelay, DelayChange      time.Duration
	MaxApples, NewAppleChance, SnakeLength int
}
type Score struct {
	X, Y  int
	Value int
	Game  *Game
}
type Cell struct {
	X, Y int
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
type Field struct {
	Screen        *tcell.Screen
	Game          Game
	X, Y          int
	Width, Height int
	Style         tcell.Style
	Cells         []Cell
	Snake         *Snake
	Apples        []Apple
}
