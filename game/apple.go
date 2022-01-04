package game

import (
	"github.com/gdamore/tcell"
	"math/rand"
)

func NewApple(f *Field, style tcell.Style) Apple {
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
