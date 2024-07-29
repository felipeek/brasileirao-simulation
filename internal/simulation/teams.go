package simulation

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
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
	TEAM_DYNAMIC_ATTRIBUTE_MORALE_NAME             = "MORALE"
	TEAM_DYNAMIC_ATTRIBUTE_PHYSICAL_CONDITION_NAME = "PHYSICAL_CONDITION"
)

type AttributeType struct {
	Name        string
	Description string
}

const (
	TEAMS_PATH = "teams/"

	// The bigger the value, the bigger the potential morale update values
	MORALE_UPDATE_STDDEV = 0.2
	// The bigger the value, the bigger the potential physical condition update values
	PHYSICAL_CONDITION_UPDATE_STDDEV = 0.3
)

var teams map[string]*Team = make(map[string]*Team)

func teamsLoad() error {
	files, err := os.ReadDir(TEAMS_PATH)
	if err != nil {
		return err
	}

	for _, dirEntry := range files {
		filePath := TEAMS_PATH + dirEntry.Name()
		raw, err := util.ReadFile(filePath)

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

func teamsGetWithName(name string) *Team {
	return teams[name]
}

func teamsGet() map[string]*Team {
	return teams
}

func teamsGetAllNames() []string {
	names := make([]string, 0, len(teams))
	for name := range teams {
		names = append(names, name)
	}
	return names
}

func teamsGetDynamicAttributeMetadata() []AttributeType {
	dynamicAttributesMetadata := make([]AttributeType, 0)

	// NOTE: The order here should be aligned with the order of the TEAM_DYNAMIC_ATTRIBUTE_* constants
	dynamicAttributesMetadata = append(dynamicAttributesMetadata, AttributeType{
		Name:        TEAM_DYNAMIC_ATTRIBUTE_MORALE_NAME,
		Description: "The morale of the squad, ranging from 0 to 10.",
	})

	dynamicAttributesMetadata = append(dynamicAttributesMetadata, AttributeType{
		Name:        TEAM_DYNAMIC_ATTRIBUTE_PHYSICAL_CONDITION_NAME,
		Description: "The physical condition of the squad, ranging from 0 to 10.",
	})

	return dynamicAttributesMetadata
}

func (t *Team) generateGptBasedRandomEvent(gptApiKey string) (string, error) {
	dynamicAttributesMetadatas := teamsGetDynamicAttributeMetadata()
	randomPos := util.RandomInt(len(dynamicAttributesMetadatas))
	attributeType := dynamicAttributesMetadatas[randomPos]
	valueDiff := util.RandomValueFromNormalDistribution(0.0, 4.0)

	err := t.changeDynamicAttribute(attributeType, valueDiff)
	if err != nil {
		return "", err
	}

	msg, err := gpt.GptRetrieveMessage(gptApiKey, t.Name, attributeType.Name, attributeType.Description, valueDiff)

	signal := '+'
	if valueDiff < 0 {
		signal = '-'
	}
	fullMsg := msg + fmt.Sprintf("\n\t- Effect: %s's %s: %c%.2f", t.Name, attributeType.Name, signal, math.Abs(valueDiff))
	return fullMsg, err
}

func (t *Team) changeDynamicAttribute(attributeType AttributeType, valueDiff float64) error {
	if attributeType.Name == TEAM_DYNAMIC_ATTRIBUTE_MORALE_NAME {
		t.changeMorale(valueDiff)
	} else if attributeType.Name == TEAM_DYNAMIC_ATTRIBUTE_PHYSICAL_CONDITION_NAME {
		t.changePhysicalCondition(valueDiff)
	} else {
		return fmt.Errorf("unknown dynamic attribute [%s]", attributeType.Name)
	}
	return nil
}

func (t *Team) changeMorale(moraleDiff float64) {
	t.DynamicAttributes.Morale += moraleDiff
	t.DynamicAttributes.Morale = util.Clamp(t.DynamicAttributes.Morale, 0, 10)
}

func (t *Team) changePhysicalCondition(physicalCondDiff float64) {
	t.DynamicAttributes.PhysicalCondition += physicalCondDiff
	t.DynamicAttributes.PhysicalCondition = util.Clamp(t.DynamicAttributes.PhysicalCondition, 0, 10)
}

func (t *Team) updateDynamicAttributes(playedFixture *Fixture) error {
	// Not very performant, but shouldn't matter...
	slices.Reverse(t.DynamicAttributes.LastFixtures)
	t.DynamicAttributes.LastFixtures = append(t.DynamicAttributes.LastFixtures, playedFixture)
	slices.Reverse(t.DynamicAttributes.LastFixtures)

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
	moraleNormalMean := float64(goalDiff) * MORALE_UPDATE_STDDEV

	t.changeMorale(util.RandomValueFromNormalDistribution(moraleNormalMean, MORALE_UPDATE_STDDEV))
	t.changePhysicalCondition(util.RandomValueFromNormalDistribution(0.0, PHYSICAL_CONDITION_UPDATE_STDDEV))
	return nil
}
