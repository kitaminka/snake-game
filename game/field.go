package game

import "github.com/gdamore/tcell"

func NewField(s *tcell.Screen, x, y, width, height int, snake *Snake, score *Score, style tcell.Style) Field {
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

func (f Field) DrawBorder(symbol rune, style tcell.Style) {
	for i := f.x - 1; i < f.x+f.width+1; i++ {
		(*f.screen).SetContent(i, f.y-1, symbol, nil, style)
	}
	for i := f.x - 1; i < f.x+f.width+1; i++ {
		(*f.screen).SetContent(i, f.y+f.height, symbol, nil, style)
	}
	for i := f.y; i < f.y+f.height; i++ {
		(*f.screen).SetContent(f.x-1, i, symbol, nil, style)
	}
	for i := f.y; i < f.y+f.height; i++ {
		(*f.screen).SetContent(f.x+f.width, i, symbol, nil, style)
	}
}
