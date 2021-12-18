package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const FieldSymbol = '-'
const SnakeSymbol = '0'
const AppleSymbol = '0'

const DefaultColor = tcell.ColorWhite
const SnakeColor = tcell.ColorGreenYellow
const AppleColor = tcell.ColorRed

const StartDelay time.Duration = 100
const MinDelay time.Duration = 50
const DelayChange time.Duration = 10
const MaxApples = 3
const NewAppleChance = 5
const SnakeLength = 5

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

func main() {
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

	f := newField(&s, width/2-50, 3, 100, 12, nil, nil, defStyle)

	snake := newSnake(&f, f.x+f.width/2, f.y+f.height/2, SnakeLength, StartDelay, snakeStyle)
	score := Score{width/2 - 50, 1, 0, &f}

	f.snake = &snake
	f.score = &score

	newApple(&f, appleStyle)

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
					newApple(f, f.apples[i].style)
				}
			}
		}

		f.snake.MoveSnake(grow)
		f.snake.DrawSnake()
		f.score.updateScore()

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

func (score *Score) updateScore() {
	f := score.field

	score.value = len(f.snake.tail) - SnakeLength + 1

	drawText(*f.screen, score.x, score.y, score.x+len(strconv.Itoa(score.value))+7, score.y, f.style, "Score: "+strconv.Itoa(score.value))
}

func newField(s *tcell.Screen, x, y, width, height int, snake *Snake, score *Score, style tcell.Style) Field {
	var f Field
	f.x = x
	f.y = y
	f.width = width
	f.height = height
	f.style = style
	f.screen = s
	f.snake = snake
	f.score = score
	f.cells = make([]Cell, width*height)

	var dx, dy int

	for i := range f.cells {
		f.cells[i] = Cell{x + dx, y + dy}
		dx++
		if dx >= width {
			dx = 0
			dy++
		}
	}

	return f
}

func (c Cell) DrawCell(s tcell.Screen, symbol rune, style tcell.Style) {
	s.SetContent(c.x, c.y, symbol, nil, style)
}

func (f Field) DrawField() {
	for _, c := range f.cells {
		c.DrawCell(*f.screen, FieldSymbol, f.style)
	}
}

func newSnake(f *Field, x, y int, length int, delay time.Duration, style tcell.Style) Snake {
	var snake Snake
	snake.field = f
	snake.head = Cell{x, y}
	snake.tail = make([]Cell, length-1)
	snake.style = style
	snake.direction.y = -1
	snake.direction.x = 0
	snake.delay = delay

	for i := range snake.tail {
		snake.tail[i] = Cell{x, y + i + 1}
	}

	return snake
}

func (snake Snake) DrawSnake() {
	snake.head.DrawCell(*snake.field.screen, SnakeSymbol, snake.style)

	for _, cell := range snake.tail {
		cell.DrawCell(*snake.field.screen, SnakeSymbol, snake.style)
	}
}

func (snake Snake) CheckSnake() bool {
	for _, cell := range snake.tail {
		if snake.head.x == cell.x && snake.head.y == cell.y {
			return true
		}
	}

	return false
}

func (snake *Snake) MoveSnake(grow bool) {
	var bufferCell Cell

	for i := range snake.tail {
		if i == 0 {
			bufferCell, snake.tail[i] = snake.tail[i], snake.head
		} else {
			snake.tail[i], bufferCell = bufferCell, snake.tail[i]
		}
		if i == len(snake.tail)-1 && grow {
			snake.tail = append(snake.tail, bufferCell)
		}
	}

	snake.head.x += snake.direction.x
	snake.head.y += snake.direction.y

	snake.BorderTeleportation()
}

func (snake *Snake) BorderTeleportation() {
	if snake.head.x < snake.field.x {
		snake.head.x = snake.field.x + snake.field.width - 1
	} else if snake.head.x >= snake.field.x+snake.field.width {
		snake.head.x = snake.field.x
	} else if snake.head.y < snake.field.y {
		snake.head.y = snake.field.y + snake.field.height - 1
	} else if snake.head.y >= snake.field.y+snake.field.height {
		snake.head.y = snake.field.y
	}
}

func newApple(f *Field, style tcell.Style) Apple {
	var apple Apple
	apple.style = style
	apple.field = f
	apple.MoveApple()
	f.apples = append(f.apples, apple)

	return apple
}

func (apple *Apple) UpdateApple() bool {
	if apple.CheckApple() {
		apple.MoveApple()
		return true
	}

	return false
}

func (apple *Apple) MoveApple() {
	touchSnake := true
	snake := apple.field.snake

	for touchSnake {
		apple.cell.x = rand.Intn(apple.field.width) + apple.field.x
		apple.cell.y = rand.Intn(apple.field.height) + apple.field.y

		if snake.head.x != apple.cell.x && snake.head.y != apple.cell.y {
			touchSnake = false
		}
		for _, cell := range snake.tail {
			if cell.x != apple.cell.x && cell.y != apple.cell.y {
				touchSnake = false
			}
		}
	}
}

func (apple Apple) CheckApple() bool {
	if apple.cell.x == apple.field.snake.head.x && apple.cell.y == apple.field.snake.head.y {
		return true
	}

	for _, cell := range apple.field.snake.tail {
		if apple.cell.x == cell.x && apple.cell.y == cell.y {
			return true
		}
	}

	return false
}

func (apple Apple) DrawApple() {
	apple.cell.DrawCell(*apple.field.screen, AppleSymbol, apple.style)
}
