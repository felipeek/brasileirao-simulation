package simulation

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/felipeek/brasileirao-simulation/internal/gpt"
	"github.com/felipeek/brasileirao-simulation/internal/util"
)

func Simulate(nonInteractive bool, gptApiKey string, enableTerminalColors bool) {
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
		err = PlayAllFixturesNonInteractive(&schedule, enableTerminalColors)
	} else {
		err = PlayAllFixturesIteractive(&schedule, gptApiKey, enableTerminalColors)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: [%s]\n", err.Error())
		os.Exit(1)
	}
}

func PlayAllFixturesNonInteractive(s *Schedule, enableTerminalColors bool) error {
	err := s.PlayAllFixtures()
	if err != nil {
		return err
	}
	s.Print(enableTerminalColors)

	standings := GenerateStandings(s)
	err = standings.Print(enableTerminalColors)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("##################################################################\n")
	fmt.Printf("The champion: [%s]!\n", standings.TeamStatistics[0].Name)
	fmt.Printf("##################################################################\n")

	return nil
}

func PlayAllFixturesIteractive(s *Schedule, gptApiKey string, enableTerminalColors bool) error {
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
		err = standings.Print(enableTerminalColors)
		if err != nil {
			return err
		}

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

			// TODO: refactor this logic of sharing the existing attributes between simulation and gpt, this is currently messy
			// also add functions to increment/decrement these attributes as part of Team, and clamp there
			if attributeType == gpt.MORALE_ATTRIBUTE.Name {
				randomTeam.DynamicAttributes.Morale += diff
				randomTeam.DynamicAttributes.Morale = util.UtilClamp(randomTeam.DynamicAttributes.Morale, 0, 10)
			} else if attributeType == gpt.PHYSICAL_CONDITION_ATTRIBUTE.Name {
				randomTeam.DynamicAttributes.PhysicalCondition += diff
				randomTeam.DynamicAttributes.PhysicalCondition = util.UtilClamp(randomTeam.DynamicAttributes.PhysicalCondition, 0, 10)
			}
		}
	}

	return nil
}
