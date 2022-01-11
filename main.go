package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/kitaminka/snake-game/game"
	"log"
	"math/rand"
	"os"
	"strconv"
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
	FieldWidth                   = 100
	FieldHeight                  = 12
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
	configuration := game.Configuration{
		StartDelay:     StartDelay,
		MinDelay:       MinDelay,
		DelayChange:    DelayChange,
		MaxApples:      MaxApples,
		NewAppleChance: NewAppleChance,
		SnakeLength:    SnakeLength,
		FieldWidth:     FieldWidth,
		FieldHeight:    FieldHeight,
	}

	width, height := s.Size()

	g := game.NewGame(&s, nil, nil, width, height, styles, configuration)

	g.InitGame()
	g.StartGame()

	s.Fini()
	fmt.Print("Game Over!\nGame result: " + strconv.Itoa(g.Score.Value))
	os.Exit(0)
}
