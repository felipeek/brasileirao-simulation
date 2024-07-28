# brasileirao-simulation

Small application that simulates a season of Brasileirao (Brazilian Football Championship).

## Running

```bash
$ go run main.go --help
  -enable-terminal-colors
    	Enable colors in the terminal output (default true)
  -gpt-api-key string
    	GPT API Key
  -non-interactive
    	Run in non-interactive mode
```

To run, simply:

```bash
$ go run main.go
```

If your terminal does not support custom font styles, or if the font styles do not integrate well with your terminal colors, disable coloring via `-enable-terminal-colors false`.

If you have an OpenAPI API Key, you can enable integration with GPT via `-git-api-key <your-api-key>`.
When GPT integration is enabled, random events are generated after each around affecting the team statuses.

Use `-non-interactive true` to simulate the whole tournament at one go.