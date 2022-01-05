package game

import (
	"github.com/gdamore/tcell"
	"math/rand"
)

func NewApple(f *Field, style tcell.Style) Apple {
	var apple Apple
	apple.Style = style
	apple.Field = f
	apple.MoveApple()
	f.Apples = append(f.Apples, apple)

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
	snake := apple.Field.Snake

	for touchSnake {
		apple.Cell.x = rand.Intn(apple.Field.Width) + apple.Field.X
		apple.Cell.y = rand.Intn(apple.Field.Height) + apple.Field.Y

		if snake.Head.x != apple.Cell.x && snake.Head.y != apple.Cell.y {
			touchSnake = false
		}
		for _, cell := range snake.Tail {
			if cell.x != apple.Cell.x && cell.y != apple.Cell.y {
				touchSnake = false
			}
		}
	}
}

func (apple Apple) CheckApple() bool {
	if apple.Cell.x == apple.Field.Snake.Head.x && apple.Cell.y == apple.Field.Snake.Head.y {
		return true
	}

	for _, cell := range apple.Field.Snake.Tail {
		if apple.Cell.x == cell.x && apple.Cell.y == cell.y {
			return true
		}
	}

	return false
}

func (apple Apple) DrawApple() {
	apple.Cell.DrawCell(*apple.Field.Screen, AppleSymbol, apple.Style)
}
