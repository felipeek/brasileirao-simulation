package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	err := TeamsLoad()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load teams: %v\n", err)
	}

	teams := TeamsGet()

	schedule, err := GenerateSchedule(teams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to generate fixtures: %v\n", err)
	}

	schedule.PlayAllFixtures()
	schedule.Print()

	standings := GenerateStandings(&schedule)
	standings.Print()
}
