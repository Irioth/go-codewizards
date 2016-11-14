package main

import (
	"codewizards/runner"
	"mystrategy"
)

func main() {
	runner.Start(mystrategy.New)
}
