package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	nonInteractive := flag.Bool("non-interactive", false, "Run in non-interactive mode")

	flag.Parse()

	err := TeamsLoad()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load teams: %v\n", err)
	}

	teams := TeamsGet()

	schedule, err := GenerateSchedule(teams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to generate fixtures: %v\n", err)
	}

	if *nonInteractive {
		err = PlayAllFixturesNonInteractive(&schedule)
	} else {
		err = PlayAllFixturesIteractive(&schedule)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: [%s]\n", err.Error())
		os.Exit(1)
	}
}

func PlayAllFixturesNonInteractive(s *Schedule) error {
	err := s.PlayAllFixtures()
	if err != nil {
		return err
	}
	s.Print()

	standings := GenerateStandings(s)
	standings.Print()

	fmt.Println()
	fmt.Printf("##################################################################\n")
	fmt.Printf("The champion: [%s]!\n", standings.TeamStatistics[0].Name)
	fmt.Printf("##################################################################\n")

	return nil
}

func PlayAllFixturesIteractive(s *Schedule) error {
	fmt.Println("Press [ENTER] to play the next round.")

	for !s.finished {
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
		err := s.PlayNextRoundFixtures()
		if err != nil {
			return err
		}
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

	return nil
}
