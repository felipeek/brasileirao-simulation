#!/bin/bash
mkdir -p bin
go build -o ./bin/bsim -v ./src/main.go ./src/teams.go ./src/util.go ./src/tournament.go ./src/fixture.go ./src/standings.go ./src/simmath.go
