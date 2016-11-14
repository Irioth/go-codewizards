package main

import (
	. "codewizards"
	"codewizards/runner"
)

func main() {
	runner.Start(NewBasicStrategy)
}

type BasicStrategy struct{}

func NewBasicStrategy() Strategy {
	return &BasicStrategy{}
}

func (s *BasicStrategy) Move(me *Wizard, world *World, game *Game, move *Move) {
	move.Speed = game.WizardForwardSpeed
	move.StrafeSpeed = game.WizardStrafeSpeed
	move.Turn = game.WizardMaxTurnAngle
	move.Action = Action_MagicMissle
}
