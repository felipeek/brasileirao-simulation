package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Team struct {
	Name     string
	Attack   int64
	Midfield int64
	Defense  int64
}

const (
	TEAMS_PATH = "teams/"
)

var teams map[string]Team = make(map[string]Team)

func TeamsLoad() error {
	files, err := os.ReadDir(TEAMS_PATH)
	if err != nil {
		return err
	}

	for _, dirEntry := range files {
		filePath := TEAMS_PATH + dirEntry.Name()
		raw, err := UtilReadFile(filePath)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open team [%s]: %v\n", filePath, err)
			return err
		}

		var team Team
		err = json.Unmarshal(raw, &team)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse team [%s]: %v\n", filePath, err)
			return err
		}

		teams[team.Name] = team
	}

	return nil
}

func TeamsGetWithName(name string) Team {
	return teams[name]
}

func TeamsGet() map[string]Team {
	return teams
}
