package game

import (
	"fmt"
	"github.com/gdamore/tcell"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func NewGame(s *tcell.Screen, f *Field, score *Score) Game {
	var game Game
	game.Screen = s
	game.Field = f
	game.Score = score
	game.Ended = false

	return game
}

func StartGame() {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	width, _ := s.Size()

	defStyle := tcell.StyleDefault.Foreground(DefaultColor)
	snakeStyle := tcell.StyleDefault.Foreground(SnakeColor)
	appleStyle := tcell.StyleDefault.Foreground(AppleColor)

	DrawText(s, width/2-5, 1, width/2+5, 1, snakeStyle, "Snake")
	DrawText(s, width/2+1, 1, width/2+5, 1, defStyle, "Game")

	rand.Seed(time.Now().UTC().UnixNano())

	f := NewField(&s, width/2-50, 3, 100, 12, nil, nil, defStyle)

	f.DrawBorder(BorderSymbol, defStyle)

	snake := NewSnake(&f, f.X+f.Width/2, f.Y+f.Height/2, SnakeLength, StartDelay, snakeStyle)
	score := Score{width/2 - 51, 1, 0, &f}

	f.Snake = &snake
	f.Score = &score

	NewApple(&f, appleStyle)

	gameOver := make(chan bool)

	go AnimationCycle(&f, gameOver)
	go EventCycle(&f, gameOver)

	<-gameOver

	s.Fini()
	fmt.Print("Game Over!\nGame result: " + strconv.Itoa(score.Value))
	os.Exit(0)
}

func EventCycle(f *Field, gameOver chan bool) {
	for {
		s := *f.Screen

		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC {
				gameOver <- true
			} else if ev.Key() == tcell.KeyDown && (f.Snake.Direction.X != 0 && f.Snake.Direction.Y != 1) {
				f.Snake.Direction.Y = 1
				f.Snake.Direction.X = 0
			} else if ev.Key() == tcell.KeyUp && (f.Snake.Direction.X != 0 && f.Snake.Direction.Y != -1) {
				f.Snake.Direction.Y = -1
				f.Snake.Direction.X = 0
			} else if ev.Key() == tcell.KeyRight && (f.Snake.Direction.X != 1 && f.Snake.Direction.Y != 0) {
				f.Snake.Direction.Y = 0
				f.Snake.Direction.X = 1
			} else if ev.Key() == tcell.KeyLeft && (f.Snake.Direction.X != -1 && f.Snake.Direction.Y != 0) {
				f.Snake.Direction.Y = 0
				f.Snake.Direction.X = -1
			}
		}
	}
}

func AnimationCycle(f *Field, gameOver chan bool) {
	for {
		f.DrawField()

		grow := false

		for i := range f.Apples {
			if f.Apples[i].UpdateApple() {
				grow = true
				if f.Snake.Delay > MinDelay {
					f.Snake.Delay -= DelayChange
				}
				if len(f.Apples) < MaxApples && rand.Intn(NewAppleChance) == 1 {
					NewApple(f, f.Apples[i].Style)
				}
			}
		}

		f.Snake.MoveSnake(grow)
		f.Snake.DrawSnake()
		f.Score.UpdateScore()

		for _, apple := range f.Apples {
			apple.DrawApple()
		}

		if f.Snake.CheckSnake() {
			gameOver <- true
		}

		(*f.Screen).Show()

		<-time.After(time.Millisecond * f.Snake.Delay)
	}
}

func (score *Score) UpdateScore() {
	f := score.Field

	score.Value = len(f.Snake.Tail) - SnakeLength + 1

	DrawText(*f.Screen, score.X, score.Y, score.X+len(strconv.Itoa(score.Value))+7, score.Y, f.Style, "Score: "+strconv.Itoa(score.Value))
}

func DrawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}
