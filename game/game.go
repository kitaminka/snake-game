package game

import (
	"github.com/gdamore/tcell"
	"math/rand"
	"strconv"
	"time"
)

func NewGame(s *tcell.Screen, f *Field, score *Score, width, height int, styles Styles) Game {
	var game Game
	game.Screen = s
	game.Field = f
	game.Score = score
	game.Ended = false
	game.Width = width
	game.Height = height
	game.Styles = styles

	return game
}

func (game Game) StartGame() {
	DrawText(*game.Screen, game.Width/2-5, 1, game.Width/2+5, 1, game.SnakeStyle, "Snake")
	DrawText(*game.Screen, game.Width/2+1, 1, game.Width/2+5, 1, game.DefaultStyle, "Game")
	game.Field.DrawBorder(BorderSymbol, game.DefaultStyle)

	gameOver := make(chan bool)

	go AnimationCycle(&game, gameOver)
	go EventCycle(&game, gameOver)

	<-gameOver
}

func EventCycle(game *Game, gameOver chan bool) {
	for {
		s := *game.Screen
		snake := game.Field.Snake

		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC {
				gameOver <- true
			} else if ev.Key() == tcell.KeyDown && (snake.Direction.X != 0 && snake.Direction.Y != 1) {
				snake.Direction.Y = 1
				snake.Direction.X = 0
			} else if ev.Key() == tcell.KeyUp && (snake.Direction.X != 0 && snake.Direction.Y != -1) {
				snake.Direction.Y = -1
				snake.Direction.X = 0
			} else if ev.Key() == tcell.KeyRight && (snake.Direction.X != 1 && snake.Direction.Y != 0) {
				snake.Direction.Y = 0
				snake.Direction.X = 1
			} else if ev.Key() == tcell.KeyLeft && (snake.Direction.X != -1 && snake.Direction.Y != 0) {
				snake.Direction.Y = 0
				snake.Direction.X = -1
			}
		}
	}
}

func AnimationCycle(game *Game, gameOver chan bool) {
	field := game.Field

	for {
		field.DrawField()

		grow := false

		for i := range field.Apples {
			if field.Apples[i].UpdateApple() {
				grow = true
				if field.Snake.Delay > MinDelay {
					field.Snake.Delay -= DelayChange
				}
				if len(field.Apples) < MaxApples && rand.Intn(NewAppleChance) == 1 {
					NewApple(field, field.Apples[i].Style)
				}
			}
		}

		field.Snake.MoveSnake(grow)
		field.Snake.DrawSnake()
		game.Score.UpdateScore()

		for _, apple := range field.Apples {
			apple.DrawApple()
		}

		if field.Snake.CheckSnake() {
			gameOver <- true
		}

		(*field.Screen).Show()

		<-time.After(time.Millisecond * field.Snake.Delay)
	}
}

func (score *Score) UpdateScore() {
	field := score.Game.Field

	score.Value = len(field.Snake.Tail) - SnakeLength + 1

	DrawText(*field.Screen, score.X, score.Y, score.X+len(strconv.Itoa(score.Value))+7, score.Y, field.Style, "Score: "+strconv.Itoa(score.Value))
}

func DrawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
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
