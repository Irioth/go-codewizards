package main

import (
	"flag"
	. "github.com/Irioth/go-codewizards"
	"github.com/Irioth/go-codewizards/runner"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 3 {
		args = []string{"127.0.0.1", "31001", "0000000000000000"}

	}
	r := runner.New(args[0]+":"+args[1], args[2], NewBasicStrategy)
	if err := r.Run(); err != nil {
		panic(err)
	}
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
