package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/felipeek/brasileirao-simulation/internal/simulation"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	nonInteractive := flag.Bool("non-interactive", false, "Run in non-interactive mode")
	gptApiKey := flag.String("gptApiKey", "", "GPT API Key")

	flag.Parse()

	simulation.Simulate(*nonInteractive, *gptApiKey)
}
