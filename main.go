package main

import (
	"github.com/gdamore/tcell"
	"log"
)

type Cell struct {
	x int
	y int
}
type Snake struct {
	head  Cell
	tail  []Cell
	field Field
	style tcell.Style
}
type Field struct {
	x      int
	y      int
	style  tcell.Style
	cells  []Cell
	screen tcell.Screen
}

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	width, height := s.Size()

	defStyle := tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorWhite)
	snakeStyle := tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDarkGreen)
	//appleStyle := tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorRed)

	drawText(s, width/2-5, 1, width/2+5, 1, defStyle, "Snake Game")
	f := newField(s, width/2-50, 3, 100, 12, defStyle)
	f.DrawField()

	snake := newSnake(f, width/2, height/2, 3, snakeStyle)
	snake.DrawSnake()

	s.Show()
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

func newField(s tcell.Screen, x int, y int, width int, height int, style tcell.Style) Field {
	var f Field
	f.x = x
	f.y = y
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
		c.DrawCell(f.screen, '-', f.style)
	}
}

func newSnake(f Field, x int, y int, length int, style tcell.Style) Snake {
	var snake Snake
	snake.field = f
	snake.head = Cell{x, y}
	snake.tail = make([]Cell, length-1)
	snake.style = style
	for i := range snake.tail {
		snake.tail[i] = Cell{x, y + i + 1}
	}

	return snake
}

func (snake Snake) DrawSnake() {
	snake.head.DrawCell(snake.field.screen, '0', snake.style)
	for _, cell := range snake.tail {
		cell.DrawCell(snake.field.screen, '0', snake.style)
	}
}
