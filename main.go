package main

import (
	"github.com/gdamore/tcell"
	"github.com/kitaminka/snake-game/game"
	"log"
	"math/rand"
	"time"
)

const (
	DefaultColor = tcell.ColorWhite
	SnakeColor   = tcell.ColorGreenYellow
	AppleColor   = tcell.ColorRed
)

const (
	StartDelay     time.Duration = 100
	MinDelay       time.Duration = 50
	DelayChange    time.Duration = 10
	MaxApples                    = 3
	NewAppleChance               = 5
	SnakeLength                  = 5
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	styles := game.Styles{
		DefaultStyle: tcell.StyleDefault.Foreground(DefaultColor),
		SnakeStyle:   tcell.StyleDefault.Foreground(SnakeColor),
		AppleStyle:   tcell.StyleDefault.Foreground(AppleColor),
	}
	width, height := s.Size()

	g := game.NewGame(&s, nil, nil, width, height, styles)
	f := game.NewField(&s, width/2-50, 3, 100, 12, nil, styles.DefaultStyle)
	snake := game.NewSnake(&f, f.X+f.Width/2, f.Y+f.Height/2, SnakeLength, StartDelay, styles.SnakeStyle)
	score := game.Score{X: width/2 - 51, Y: 1, Game: &g}

	f.Snake = &snake
	g.Score = &score
	g.Field = &f

	game.NewApple(&f, styles.AppleStyle)
	g.StartGame()
}
