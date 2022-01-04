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

	drawText(s, width/2-5, 1, width/2+5, 1, snakeStyle, "Snake")
	drawText(s, width/2+1, 1, width/2+5, 1, defStyle, "Game")

	rand.Seed(time.Now().UTC().UnixNano())

	f := NewField(&s, width/2-50, 3, 100, 12, nil, nil, defStyle)

	f.DrawBorder(BorderSymbol, defStyle)

	snake := NewSnake(&f, f.x+f.width/2, f.y+f.height/2, SnakeLength, StartDelay, snakeStyle)
	score := Score{width/2 - 51, 1, 0, &f}

	f.snake = &snake
	f.score = &score

	NewApple(&f, appleStyle)

	gameOver := make(chan bool)

	go animationCycle(&f, gameOver)
	go gameCycle(&f, gameOver)

	<-gameOver

	s.Fini()
	fmt.Print("Game Over!\nGame result: " + strconv.Itoa(score.value))
	os.Exit(0)
}

func gameCycle(f *Field, gameOver chan bool) {
	for {
		s := *f.screen

		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC {
				gameOver <- true
			} else if ev.Key() == tcell.KeyDown && (f.snake.direction.x != 0 && f.snake.direction.y != 1) {
				f.snake.direction.y = 1
				f.snake.direction.x = 0
			} else if ev.Key() == tcell.KeyUp && (f.snake.direction.x != 0 && f.snake.direction.y != -1) {
				f.snake.direction.y = -1
				f.snake.direction.x = 0
			} else if ev.Key() == tcell.KeyRight && (f.snake.direction.x != 1 && f.snake.direction.y != 0) {
				f.snake.direction.y = 0
				f.snake.direction.x = 1
			} else if ev.Key() == tcell.KeyLeft && (f.snake.direction.x != -1 && f.snake.direction.y != 0) {
				f.snake.direction.y = 0
				f.snake.direction.x = -1
			}
		}
	}
}

func animationCycle(f *Field, gameOver chan bool) {
	for {
		f.DrawField()

		grow := false

		for i := range f.apples {
			if f.apples[i].UpdateApple() {
				grow = true
				if f.snake.delay > MinDelay {
					f.snake.delay -= DelayChange
				}
				if len(f.apples) < MaxApples && rand.Intn(NewAppleChance) == 1 {
					NewApple(f, f.apples[i].style)
				}
			}
		}

		f.snake.MoveSnake(grow)
		f.snake.DrawSnake()
		f.score.UpdateScore()

		for _, apple := range f.apples {
			apple.DrawApple()
		}

		if f.snake.CheckSnake() {
			gameOver <- true
		}

		(*f.screen).Show()

		<-time.After(time.Millisecond * f.snake.delay)
	}
}

func (score *Score) UpdateScore() {
	f := score.field

	score.value = len(f.snake.tail) - SnakeLength + 1

	drawText(*f.screen, score.x, score.y, score.x+len(strconv.Itoa(score.value))+7, score.y, f.style, "Score: "+strconv.Itoa(score.value))
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
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
