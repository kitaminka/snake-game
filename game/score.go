package game

import "strconv"

func NewScore(x, y int, game *Game) Score {
	var score Score
	score.X = x
	score.Y = y
	score.Game = game

	return score
}

func (score *Score) UpdateScore() {
	field := score.Game.Field

	score.Value = len(field.Snake.Tail) - score.Game.SnakeLength + 1

	DrawText(*field.Screen, score.X, score.Y, score.X+len(strconv.Itoa(score.Value))+7, score.Y, field.Style, "Score: "+strconv.Itoa(score.Value))
}
