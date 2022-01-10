package game

import "time"

const (
	FieldSymbol  = ' '
	SnakeSymbol  = '▇'
	BorderSymbol = '▇'
	AppleSymbol  = '●'
)

const (
	StartDelay     time.Duration = 100
	MinDelay       time.Duration = 50
	DelayChange    time.Duration = 10
	MaxApples                    = 3
	NewAppleChance               = 5
	SnakeLength                  = 5
)
