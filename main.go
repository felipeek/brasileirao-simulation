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
	gptApiKey := flag.String("gpt-api-key", "", "GPT API Key")
	disableTerminalColors := flag.Bool("disable-terminal-colors", false, "Disable colors in the terminal output")

	flag.Parse()

	simulation.Simulate(*nonInteractive, *gptApiKey, !*disableTerminalColors)
}
