package game

import "github.com/gdamore/tcell"

func NewField(s *tcell.Screen, x, y, width, height int, snake *Snake, style tcell.Style) Field {
	var field Field
	field.X = x
	field.Y = y
	field.Width = width
	field.Height = height
	field.Style = style
	field.Screen = s
	field.Snake = snake
	field.Cells = make([]Cell, width*height)

	var dx, dy int

	for i := range field.Cells {
		field.Cells[i] = Cell{x + dx, y + dy}
		dx++
		if dx >= width {
			dx = 0
			dy++
		}
	}

	return field
}

func (c Cell) DrawCell(s tcell.Screen, symbol rune, style tcell.Style) {
	s.SetContent(c.X, c.Y, symbol, nil, style)
}

func (field Field) DrawField() {
	for _, c := range field.Cells {
		c.DrawCell(*field.Screen, FieldSymbol, field.Style)
	}
}

func (field Field) DrawBorder(symbol rune, style tcell.Style) {
	for i := field.X - 1; i < field.X+field.Width+1; i++ {
		(*field.Screen).SetContent(i, field.Y-1, symbol, nil, style)
	}
	for i := field.X - 1; i < field.X+field.Width+1; i++ {
		(*field.Screen).SetContent(i, field.Y+field.Height, symbol, nil, style)
	}
	for i := field.Y; i < field.Y+field.Height; i++ {
		(*field.Screen).SetContent(field.X-1, i, symbol, nil, style)
	}
	for i := field.Y; i < field.Y+field.Height; i++ {
		(*field.Screen).SetContent(field.X+field.Width, i, symbol, nil, style)
	}
}
