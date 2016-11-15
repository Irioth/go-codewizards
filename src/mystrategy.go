package main

import (
	. "codewizards"
)

type MyStrategy struct{}

func New() Strategy {
	return &MyStrategy{}
}

func (s *MyStrategy) Move(me *Wizard, world *World, game *Game, move *Move) {
	// put your code here
}
