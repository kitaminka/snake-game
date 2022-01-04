package game

import (
	"github.com/gdamore/tcell"
	"time"
)

func NewSnake(f *Field, x, y int, length int, delay time.Duration, style tcell.Style) Snake {
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
