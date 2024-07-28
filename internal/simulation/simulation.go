package simulation

import (
	"bufio"
	"fmt"
	"os"

	"github.com/felipeek/brasileirao-simulation/internal/util"
)

func Simulate(nonInteractive bool, gptApiKey string, enableTerminalColors bool) {
	err := teamsLoad()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load teams: %v\n", err)
	}

	teams := teamsGet()

	schedule, err := generateSchedule(teams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to generate fixtures: %v\n", err)
	}

	if nonInteractive {
		err = playAllFixturesNonInteractive(&schedule, enableTerminalColors)
	} else {
		err = playAllFixturesIteractive(&schedule, gptApiKey, enableTerminalColors)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: [%s]\n", err.Error())
		os.Exit(1)
	}
}

func playAllFixturesNonInteractive(s *Schedule, enableTerminalColors bool) error {
	err := s.playAllFixtures()
	if err != nil {
		return err
	}
	s.print(enableTerminalColors)

	standings := standingsGenerate(s)
	err = standings.print(enableTerminalColors)
	if err != nil {
		return err
	}

	printChampionMessage(standings.TeamStatistics[0].Name)

	return nil
}

func playAllFixturesIteractive(s *Schedule, gptApiKey string, enableTerminalColors bool) error {
	fmt.Println("Press [ENTER] to play the next round.")

	for !s.finished {
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
		err := s.playNextRoundFixtures()
		if err != nil {
			return err
		}
		s.printLastPlayedRound(enableTerminalColors)

		standings := standingsGenerate(s)
		err = standings.print(enableTerminalColors)
		if err != nil {
			return err
		}

		if s.finished {
			printChampionMessage(standings.TeamStatistics[0].Name)
			return nil
		}

		if gptApiKey != "" {
			reader.ReadString('\n')

			teamsNames := teamsGetAllNames()
			teamName := util.UtilRandomChoiceStr(teamsNames...).(string)
			randomTeam := teamsGetWithName(teamName)
			eventStr, err := randomTeam.generateGptBasedRandomEvent(gptApiKey)
			if err != nil {
				return err
			}
			fmt.Printf("Round [%d] Event:\n", s.currentRoundIdx+1)
			fmt.Printf("\t- %s\n", eventStr)
			fmt.Printf("\n\n")
		}
	}

	return nil
}

func printChampionMessage(championName string) {
	fmt.Println()
	fmt.Printf("##################################################################\n")
	fmt.Printf("The champion: [%s]!\n", championName)
	fmt.Printf("##################################################################\n")
}
