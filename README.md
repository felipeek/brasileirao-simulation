# brasileirao-simulation

Small application that simulates a season of Brasileirao (Brazilian Football Championship).

## Running

```bash
$ go run main.go -help
  -disable-terminal-colors
    	Disable colors in the terminal output
  -gpt-api-key string
    	GPT API Key
  -non-interactive
    	Run in non-interactive mode
```

To run, simply:

```bash
$ go run main.go
```

If your terminal does not support custom font styles, or if the font styles do not integrate well with your terminal colors, disable coloring via `-disable-terminal-colors`.

If you have an OpenAI API Key, you can enable integration with GPT via `-git-api-key <your-api-key>`.
When GPT integration is enabled, random events are generated after each around affecting the team statuses.

Use `-non-interactive` to simulate the whole tournament at one go.