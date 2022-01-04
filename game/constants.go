package game

import (
	"github.com/gdamore/tcell"
	"time"
)

const FieldSymbol = ' '
const SnakeSymbol = '▇'
const BorderSymbol = '▇'
const AppleSymbol = '●'

const DefaultColor = tcell.ColorWhite
const SnakeColor = tcell.ColorGreenYellow
const AppleColor = tcell.ColorRed

const StartDelay time.Duration = 100
const MinDelay time.Duration = 50
const DelayChange time.Duration = 10
const MaxApples = 3
const NewAppleChance = 5
const SnakeLength = 5
