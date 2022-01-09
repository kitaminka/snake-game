package game

import (
	"github.com/gdamore/tcell"
	"math/rand"
)

func NewApple(field *Field, style tcell.Style) Apple {
	var apple Apple
	apple.Style = style
	apple.Field = field
	apple.MoveApple()
	field.Apples = append(field.Apples, apple)

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
		apple.Cell.X = rand.Intn(apple.Field.Width) + apple.Field.X
		apple.Cell.Y = rand.Intn(apple.Field.Height) + apple.Field.Y

		if snake.Head.X != apple.Cell.X && snake.Head.Y != apple.Cell.Y {
			touchSnake = false
		}
		for _, cell := range snake.Tail {
			if cell.X != apple.Cell.X && cell.Y != apple.Cell.Y {
				touchSnake = false
			}
		}
	}
}

func (apple Apple) CheckApple() bool {
	if apple.Cell.X == apple.Field.Snake.Head.X && apple.Cell.Y == apple.Field.Snake.Head.Y {
		return true
	}

	for _, cell := range apple.Field.Snake.Tail {
		if apple.Cell.X == cell.X && apple.Cell.Y == cell.Y {
			return true
		}
	}

	return false
}

func (apple Apple) DrawApple() {
	apple.Cell.DrawCell(*apple.Field.Screen, AppleSymbol, apple.Style)
}
