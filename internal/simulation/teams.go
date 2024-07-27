package simulation

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/felipeek/brasileirao-simulation/internal/gpt"
	"github.com/felipeek/brasileirao-simulation/internal/util"
)

type Team struct {
	Name              string
	Attack            float64 // 0-10
	Midfield          float64 // 0-10
	Defense           float64 // 0-10
	HomeFactor        float64 // 0-10
	DynamicAttributes TeamDynamicAttributes
}

type TeamDynamicAttributes struct {
	LastFixtures      []*Fixture
	Morale            float64 // 0-10
	PhysicalCondition float64 // 0-10
}

const (
	TEAMS_PATH = "teams/"

	// The bigger the value, the bigger the potential morale update values
	MORALE_UPDATE_STDDEV = 0.2
	// The bigger the value, the bigger the potential physical condition update values
	PHYSICAL_CONDITION_UPDATE_STDDEV = 0.2
)

var teams map[string]*Team = make(map[string]*Team)

func TeamsLoad() error {
	files, err := os.ReadDir(TEAMS_PATH)
	if err != nil {
		return err
	}

	for _, dirEntry := range files {
		filePath := TEAMS_PATH + dirEntry.Name()
		raw, err := util.UtilReadFile(filePath)

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

		team.DynamicAttributes.LastFixtures = make([]*Fixture, 0)
		team.DynamicAttributes.Morale = 5
		team.DynamicAttributes.PhysicalCondition = 5

		teams[team.Name] = &team
	}

	return nil
}

func TeamsGetWithName(name string) *Team {
	return teams[name]
}

func TeamsGet() map[string]*Team {
	return teams
}

func TeamsGetAllNames() []string {
	names := make([]string, 0, len(teams))
	for name := range teams {
		names = append(names, name)
	}
	return names
}

func (t *Team) GenerateGptBasedRandomEvent(gptApiKey string, roundNum int) (string, float64, string, error) {
	valueDiff := SimUtilRandomValueFromNormalDistribution(0.0, 4.0)

	attributeType, msg, err := gpt.GptRetrieveMessage(gptApiKey, t.Name, valueDiff, roundNum)
	return attributeType, valueDiff, msg, err
}

func (t *Team) updateDynamicAttributes(playedFixture *Fixture) error {
	// Not very performant, but shouldn't matter...
	slices.Reverse(t.DynamicAttributes.LastFixtures)
	t.DynamicAttributes.LastFixtures = append(t.DynamicAttributes.LastFixtures, playedFixture)
	slices.Reverse(t.DynamicAttributes.LastFixtures)

	t.DynamicAttributes.PhysicalCondition = t.DynamicAttributes.PhysicalCondition + SimUtilRandomValueFromNormalDistribution(0.0, PHYSICAL_CONDITION_UPDATE_STDDEV)
	t.DynamicAttributes.PhysicalCondition = util.UtilClamp(t.DynamicAttributes.PhysicalCondition, 0, 10)

	goalDiff := 0
	if playedFixture.homeTeam == t.Name {
		goalDiff = playedFixture.homeTeamScore - playedFixture.awayTeamScore
	} else if playedFixture.awayTeam == t.Name {
		goalDiff = playedFixture.awayTeamScore - playedFixture.homeTeamScore
	} else {
		return errors.New("Team did not play received match")
	}

	// this is proportional to the stddev because, if the stddev is increased, we want to also shift the mean proportionally,
	// to ensure that the match results will continue having a meaningful impact on the morale update.
	normalMean := float64(goalDiff) * MORALE_UPDATE_STDDEV
	t.DynamicAttributes.Morale = t.DynamicAttributes.Morale + SimUtilRandomValueFromNormalDistribution(normalMean, MORALE_UPDATE_STDDEV)
	t.DynamicAttributes.Morale = util.UtilClamp(t.DynamicAttributes.Morale, 0, 10)
	return nil
}
