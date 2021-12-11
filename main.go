package main

import (
	"github.com/gdamore/tcell"
	"log"
	"math/rand"
	"os"
	"time"
)

const FieldSymbol = '-'
const SnakeSymbol = '0'
const AppleSymbol = '0'

const DefaultColor = tcell.ColorWhite
const SnakeColor = tcell.ColorDarkGreen
const AppleColor = tcell.ColorRed

const StartDelay time.Duration = 100
const MinDelay time.Duration = 30
const DelayChange time.Duration = 10
const MaxApples = 3
const NewAppleChance = 5

type Cell struct {
	x int
	y int
}
type Snake struct {
	head      Cell
	tail      []Cell
	field     *Field
	style     tcell.Style
	direction struct {
		x int
		y int
	}
	delay time.Duration
}
type Apple struct {
	cell  Cell
	style tcell.Style
	field *Field
}
type Field struct {
	x      int
	y      int
	width  int
	height int
	style  tcell.Style
	cells  []Cell
	screen *tcell.Screen
	snake  *Snake
	apples []Apple
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

	drawText(s, width/2-5, 1, width/2+5, 1, defStyle, "Snake Game")

	rand.Seed(time.Now().UTC().UnixNano())

	f := newField(&s, width/2-50, 3, 100, 12, nil, defStyle)
	snake := newSnake(&f, f.x+f.width/2, f.y+f.height/2, 5, StartDelay, snakeStyle)
	f.snake = &snake
	newApple(&f, rand.Intn(f.width)+f.x, rand.Intn(f.height)+f.y, appleStyle)

	gameOver := make(chan bool)

	go animationCycle(&f, gameOver)
	go gameCycle(&f, gameOver)

	<-gameOver

	s.Fini()
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
			switch ev.Key() {
			case tcell.KeyCtrlC:
				gameOver <- true
			case tcell.KeyDown:
				f.snake.direction.y = 1
				f.snake.direction.x = 0
			case tcell.KeyUp:
				f.snake.direction.y = -1
				f.snake.direction.x = 0
			case tcell.KeyRight:
				f.snake.direction.y = 0
				f.snake.direction.x = 1
			case tcell.KeyLeft:
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
					newApple(f, rand.Intn(f.width)+f.x, rand.Intn(f.height)+f.y, f.apples[i].style)
				}
			}
		}

		f.snake.MoveSnake(grow)
		f.snake.DrawSnake()

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

func newField(s *tcell.Screen, x int, y int, width int, height int, snake *Snake, style tcell.Style) Field {
	var f Field
	f.x = x
	f.y = y
	f.width = width
	f.height = height
	f.style = style
	f.screen = s
	f.snake = snake
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

func newSnake(f *Field, x int, y int, length int, delay time.Duration, style tcell.Style) Snake {
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

func newApple(f *Field, x int, y int, style tcell.Style) Apple {
	var apple Apple
	apple.cell = Cell{x, y}
	apple.style = style
	apple.field = f
	f.apples = append(f.apples, apple)

	return apple
}

func (apple *Apple) UpdateApple() bool {
	if apple.CheckApple() {
		apple.cell.x = rand.Intn(apple.field.width) + apple.field.x
		apple.cell.y = rand.Intn(apple.field.height) + apple.field.y
		return true
	}

	return false
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
