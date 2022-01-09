package main

import (
	"github.com/gdamore/tcell"
	"log"
)

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	width, _ := s.Size()
}
