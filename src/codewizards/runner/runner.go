package runner

import (
	"flag"
	. "codewizards"
)

type Runner struct {
	addr    string
	token   string
	factory StrategyFactory
}

type StrategyFactory func() Strategy

func Start(factory StrategyFactory) {
	flag.Parse()
	args := flag.Args()
	if len(args) != 3 {
		args = []string{"127.0.0.1", "31001", "0000000000000000"}

	}
	r := New(args[0]+":"+args[1], args[2], factory)
	if err := r.Run(); err != nil {
		panic(err)
	}
}

func New(addr, token string, factory StrategyFactory) *Runner {
	return &Runner{addr, token, factory}
}

func (r *Runner) Run() error {
	client, err := NewClient(r.addr)
	if err != nil {
		return err
	}
	defer client.Close()

	client.WriteToken(r.token)
	client.WriteProtocolVersion()
	teamSize := client.ReadTeamSize()
	game := client.ReadGameContext()

	strategies := make([]Strategy, teamSize)
	for i := range strategies {
		strategies[i] = r.factory()
	}

	playerContext := client.ReadPlayerContext()
	for playerContext != nil {
		playerWizards := playerContext.Wizards
		if playerWizards == nil || len(playerWizards) != teamSize {
			break
		}

		moves := make([]*Move, teamSize)

		for i, wizard := range playerWizards {
			moves[i] = NewMove()
			strategies[i].Move(wizard, playerContext.World, game, moves[i])
		}
		client.WriteMoves(moves)

		playerContext = client.ReadPlayerContext()
	}

	return nil
}
