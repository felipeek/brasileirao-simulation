package simulation

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/felipeek/brasileirao-simulation/internal/util"
)

func Simulate(nonInteractive bool, gptApiKey string) {
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

	if nonInteractive {
		err = PlayAllFixturesNonInteractive(&schedule)
	} else {
		err = PlayAllFixturesIteractive(&schedule, gptApiKey)
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

func PlayAllFixturesIteractive(s *Schedule, gptApiKey string) error {
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
			return nil
		}

		if gptApiKey != "" {
			reader.ReadString('\n')

			teamsNames := TeamsGetAllNames()
			teamName := util.UtilRandomChoiceStr(teamsNames...).(string)
			randomTeam := TeamsGetWithName(teamName)
			attributeType, diff, eventStr, err := randomTeam.GenerateGptBasedRandomEvent(gptApiKey, s.currentRoundIdx+1)
			if err != nil {
				return err
			}
			fmt.Printf("Round [%d] Event:\n", s.currentRoundIdx+1)
			fmt.Printf("\t- %s\n", eventStr)

			signal := '+'
			if diff < 0 {
				signal = '-'
			}
			fmt.Printf("\t- Effect: %s's %s: %c%.2f\n\n", teamName, attributeType, signal, math.Abs(diff))
		}
	}

	return nil
}
