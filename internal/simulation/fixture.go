package simulation

import (
	"errors"

	"github.com/felipeek/brasileirao-simulation/internal/util"
)

type Fixture struct {
	homeTeam      string
	awayTeam      string
	homeTeamScore int
	awayTeamScore int
	played        bool
}

const (
	HOME_BONUS_FACTOR                      = 2.0
	RECENT_FORM_CONTRIBUTION_IMPACT        = 0.08
	MORALE_CONTRIBUTION_IMPACT             = 0.05
	PHYSICAL_CONDITION_CONTRIBUTION_IMPACT = 0.068
)

var recentFormMatchContributions = [5]float64{0.35, 0.20, 0.15, 0.15, 0.15}

func (f *Fixture) play() error {
	homeTeam := teamsGetWithName(f.homeTeam)
	awayTeam := teamsGetWithName(f.awayTeam)

	// Additional strength given to the home team (home factor)
	homeStadiumStrength := HOME_BONUS_FACTOR * (homeTeam.HomeFactor / 10)

	// Calculate home team recent form contribution
	homeTeamRawFormContribution, err := calculateFormContribution(f.homeTeam, homeTeam.DynamicAttributes.LastFixtures)
	if err != nil {
		return err
	}
	homeTeamFormContribution := util.GetMultiplierFromContributionFactor(homeTeamRawFormContribution, RECENT_FORM_CONTRIBUTION_IMPACT)

	// Calculate away team recent form contribution
	awayTeamRawFormContribution, err := calculateFormContribution(f.awayTeam, awayTeam.DynamicAttributes.LastFixtures)
	if err != nil {
		return err
	}
	awayTeamFormContribution := util.GetMultiplierFromContributionFactor(awayTeamRawFormContribution, RECENT_FORM_CONTRIBUTION_IMPACT)

	// Caclulate morale contribution
	homeTeamMoraleContribution := util.GetMultiplierFromContributionFactor(homeTeam.DynamicAttributes.Morale, MORALE_CONTRIBUTION_IMPACT)
	awayTeamMoraleContribution := util.GetMultiplierFromContributionFactor(awayTeam.DynamicAttributes.Morale, MORALE_CONTRIBUTION_IMPACT)

	// Caclulate physical condition contribution
	homeTeamPhysicalConditionContribution := util.GetMultiplierFromContributionFactor(homeTeam.DynamicAttributes.PhysicalCondition, PHYSICAL_CONDITION_CONTRIBUTION_IMPACT)
	awayTeamPhysicalConditionContribution := util.GetMultiplierFromContributionFactor(awayTeam.DynamicAttributes.PhysicalCondition, PHYSICAL_CONDITION_CONTRIBUTION_IMPACT)

	// Home team strength factors
	homeAttackStrength := 1.5*homeTeam.Attack + homeTeam.Midfield
	homeDefenseStrength := 1.5*homeTeam.Defense + homeTeam.Midfield

	// Away team strength factors
	awayDefenseStrength := 1.5*awayTeam.Defense + awayTeam.Midfield
	awayAttackStrength := 1.5*awayTeam.Attack + awayTeam.Midfield

	// Calculate home/away strength without other contributions
	homeRawStrength := homeAttackStrength / (1 + awayDefenseStrength/homeAttackStrength)
	awayRawStrength := awayAttackStrength / (1 + homeDefenseStrength/awayAttackStrength)

	// Final non-attentuated strength of each team for this match
	homeStrength := homeStadiumStrength * homeTeamFormContribution * homeTeamMoraleContribution * homeTeamPhysicalConditionContribution * homeRawStrength
	awayStrength := awayTeamFormContribution * awayTeamMoraleContribution * awayTeamPhysicalConditionContribution * awayRawStrength

	// Attenuate strengths by employing a log-based function
	homeLambda := util.AttenuateStrength(homeStrength)
	awayLambda := util.AttenuateStrength(awayStrength)

	// Generate final scores based on a poisson distribution
	f.homeTeamScore = util.PoissonKnuth(homeLambda)
	f.awayTeamScore = util.PoissonKnuth(awayLambda)

	//fmt.Printf("%s: %f -> %f\n", f.homeTeam, homeStrength, homeLambda)
	//fmt.Printf("%s: %f -> %f\n\n", f.awayTeam, awayStrength, awayLambda)

	f.played = true

	err = homeTeam.updateDynamicAttributes(f)
	if err != nil {
		return err
	}

	err = awayTeam.updateDynamicAttributes(f)
	if err != nil {
		return err
	}

	return nil
}

// Return a contribution based on recent form in the interval 0-10
func calculateFormContribution(teamName string, teamAlreadyPlayedFixtures []*Fixture) (float64, error) {
	// Last matches are analyzed and summed to this contribution.
	formContribution := 0.0

	for i, contribution := range recentFormMatchContributions {
		if len(teamAlreadyPlayedFixtures) <= i {
			// if the match was not played, we add half, simulating a 'neutral' contribution
			formContribution += contribution / 2.0
		} else {
			fixtureContribution, err := getFixtureFormContribution(teamAlreadyPlayedFixtures[i], teamName, contribution)
			if err != nil {
				return 0.0, err
			}
			formContribution += fixtureContribution
		}
	}

	return formContribution * 10, nil
}

func getFixtureFormContribution(fixture *Fixture, teamName string, contributionPctg float64) (float64, error) {
	scored, conceded, err := getTeamScoredAndConcededGoalsInFixture(teamName, fixture)
	if err != nil {
		return 0.0, err
	}

	if scored > conceded {
		return contributionPctg, nil
	} else if scored == conceded {
		return contributionPctg / 2.0, nil
	}

	return 0.0, nil
}

func getTeamScoredAndConcededGoalsInFixture(teamName string, fixture *Fixture) (int, int, error) {
	if fixture.homeTeam == teamName {
		return fixture.homeTeamScore, fixture.awayTeamScore, nil
	} else if fixture.awayTeam == teamName {
		return fixture.awayTeamScore, fixture.homeTeamScore, nil
	}
	return -1, -1, errors.New("Fixture does not have target team")
}
