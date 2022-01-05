package game

import "github.com/gdamore/tcell"

func NewField(s *tcell.Screen, x, y, width, height int, snake *Snake, score *Score, style tcell.Style) Field {
	var f Field
	f.X = x
	f.Y = y
	f.Width = width
	f.Height = height
	f.Style = style
	f.Screen = s
	f.Snake = snake
	f.Score = score
	f.Cells = make([]Cell, width*height)

	var dx, dy int

	for i := range f.Cells {
		f.Cells[i] = Cell{x + dx, y + dy}
		dx++
		if dx >= width {
			dx = 0
			dy++
		}
	}

	return f
}

func (c Cell) DrawCell(s tcell.Screen, symbol rune, style tcell.Style) {
	s.SetContent(c.X, c.Y, symbol, nil, style)
}

func (f Field) DrawField() {
	for _, c := range f.Cells {
		c.DrawCell(*f.Screen, FieldSymbol, f.Style)
	}
}

func (f Field) DrawBorder(symbol rune, style tcell.Style) {
	for i := f.X - 1; i < f.X+f.Width+1; i++ {
		(*f.Screen).SetContent(i, f.Y-1, symbol, nil, style)
	}
	for i := f.X - 1; i < f.X+f.Width+1; i++ {
		(*f.Screen).SetContent(i, f.Y+f.Height, symbol, nil, style)
	}
	for i := f.Y; i < f.Y+f.Height; i++ {
		(*f.Screen).SetContent(f.X-1, i, symbol, nil, style)
	}
	for i := f.Y; i < f.Y+f.Height; i++ {
		(*f.Screen).SetContent(f.X+f.Width, i, symbol, nil, style)
	}
}
