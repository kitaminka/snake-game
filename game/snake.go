package game

import (
	"github.com/gdamore/tcell"
	"time"
)

func NewSnake(f *Field, x, y int, length int, delay time.Duration, style tcell.Style) Snake {
	var snake Snake
	snake.Field = f
	snake.Head = Cell{x, y}
	snake.Tail = make([]Cell, length-1)
	snake.Style = style
	snake.Direction.Y = -1
	snake.Direction.X = 0
	snake.Delay = delay

	for i := range snake.Tail {
		snake.Tail[i] = Cell{x, y + i + 1}
	}

	return snake
}

func (snake Snake) DrawSnake() {
	snake.Head.DrawCell(*snake.Field.Screen, SnakeSymbol, snake.Style)

	for _, cell := range snake.Tail {
		cell.DrawCell(*snake.Field.Screen, SnakeSymbol, snake.Style)
	}
}

func (snake Snake) CheckSnake() bool {
	for _, cell := range snake.Tail {
		if snake.Head.x == cell.x && snake.Head.y == cell.y {
			return true
		}
	}

	return false
}

func (snake *Snake) MoveSnake(grow bool) {
	var bufferCell Cell

	for i := range snake.Tail {
		if i == 0 {
			bufferCell, snake.Tail[i] = snake.Tail[i], snake.Head
		} else {
			snake.Tail[i], bufferCell = bufferCell, snake.Tail[i]
		}
		if i == len(snake.Tail)-1 && grow {
			snake.Tail = append(snake.Tail, bufferCell)
		}
	}

	snake.Head.x += snake.Direction.X
	snake.Head.y += snake.Direction.Y

	snake.BorderTeleportation()
}

func (snake *Snake) BorderTeleportation() {
	if snake.Head.x < snake.Field.X {
		snake.Head.x = snake.Field.X + snake.Field.Width - 1
	} else if snake.Head.x >= snake.Field.X+snake.Field.Width {
		snake.Head.x = snake.Field.X
	} else if snake.Head.y < snake.Field.Y {
		snake.Head.y = snake.Field.Y + snake.Field.Height - 1
	} else if snake.Head.y >= snake.Field.Y+snake.Field.Height {
		snake.Head.y = snake.Field.Y
	}
}
