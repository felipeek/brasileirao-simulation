package main

import (
	"bufio"
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

	PlayAllFixturesIteractive(&schedule)
}

func PlayAllFixturesNonInteractive(s *Schedule) {
	s.PlayAllFixtures()
	s.Print()

	standings := GenerateStandings(s)
	standings.Print()
}

func PlayAllFixturesIteractive(s *Schedule) {
	fmt.Println("Press [ENTER] to play the next round.")

	for !s.finished {
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
		s.PlayNextRoundFixtures()
		s.PrintLastPlayedRound()

		standings := GenerateStandings(s)
		standings.Print()

		if s.finished {
			fmt.Println()
			fmt.Printf("##################################################################\n")
			fmt.Printf("The champion: [%s]!\n", standings.TeamStatistics[0].Name)
			fmt.Printf("##################################################################\n")
		}
	}

}
