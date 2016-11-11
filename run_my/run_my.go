package main

import (
	"flag"
	"github.com/Irioth/go-codewizards/runner"
	"strategy"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 3 {
		args = []string{"127.0.0.1", "31001", "0000000000000000"}

	}
	r := runner.New(args[0]+":"+args[1], args[2], strategy.New)
	if err := r.Run(); err != nil {
		panic(err)
	}
}
