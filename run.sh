#!/bin/bash

# Check if the gpt-key file exists
if [ -f gpt-key ]; then
    # Read the content of the gpt-key file
    GPT_KEY=$(<gpt-key)
    # Run the Go program with the --gpt-key parameter
    go run main.go --gptApiKey "$GPT_KEY"
else
    # Run the Go program without the --gpt-key parameter
    go run main.go
fi
