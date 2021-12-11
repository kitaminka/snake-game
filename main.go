package main

import (
	"github.com/gdamore/tcell"
	"log"
	"os"
	"time"
)

const FieldSymbol = '-'
const SnakeSymbol = '0'

const DefaultColor = tcell.ColorWhite
const SnakeColor = tcell.ColorDarkGreen

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
type Field struct {
	x      int
	y      int
	width  int
	height int
	style  tcell.Style
	cells  []Cell
	screen *tcell.Screen
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
	//appleStyle := tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorRed)

	drawText(s, width/2-5, 1, width/2+5, 1, defStyle, "Snake Game")

	f := newField(&s, width/2-50, 3, 100, 12, defStyle)
	snake := newSnake(f, f.x+f.width/2, f.y+f.height/2, 5, snakeStyle)

	ch := make(chan bool)

	go animationCycle(&snake)
	go gameCycle(&snake, ch)

	<-ch

	s.Fini()
	os.Exit(0)
}

func gameCycle(snake *Snake, ch chan bool) {
	for {
		s := *snake.field.screen

		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyCtrlC:
				ch <- true
			case tcell.KeyDown:
				snake.direction.y = 1
				snake.direction.x = 0
			case tcell.KeyUp:
				snake.direction.y = -1
				snake.direction.x = 0
			case tcell.KeyRight:
				snake.direction.y = 0
				snake.direction.x = 1
			case tcell.KeyLeft:
				snake.direction.y = 0
				snake.direction.x = -1
			}
		}
	}
}
func animationCycle(snake *Snake) {
	for {
		snake.field.DrawField()
		snake.DrawSnake()
		snake.MoveSnake()
		(*snake.field.screen).Show()
		<-time.After(time.Millisecond * snake.delay)
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

func newField(s *tcell.Screen, x int, y int, width int, height int, style tcell.Style) Field {
	var f Field
	f.x = x
	f.y = y
	f.width = width
	f.height = height
	f.style = style
	f.screen = s
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

func newSnake(f Field, x int, y int, length int, style tcell.Style) Snake {
	var snake Snake
	snake.field = &f
	snake.head = Cell{x, y}
	snake.tail = make([]Cell, length-1)
	snake.style = style
	snake.direction.y = -1
	snake.direction.x = 0
	snake.delay = 100

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

func (snake *Snake) MoveSnake() {
	var bufferCell Cell
	for i := range snake.tail {
		if i == 0 {
			bufferCell, snake.tail[i] = snake.tail[i], snake.head
		} else {
			snake.tail[i], bufferCell = bufferCell, snake.tail[i]
		}
	}
	snake.head.x += snake.direction.x
	snake.head.y += snake.direction.y
}
