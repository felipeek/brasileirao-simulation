package main

import (
	"errors"
	"math"
	"math/rand"
)

type Fixture struct {
	homeTeam      string
	awayTeam      string
	homeTeamScore int
	awayTeamScore int
	played        bool
}

const (
	LOG_ADJUST_FACTOR        = 0.5 // the smaller, the more 'balanced' the results
	HOME_BONUS_FACTOR        = 2.0
	RECENT_FORM_BONUS_FACTOR = 1.0
)

func (f *Fixture) Play(homeTeamLastPlayedFixtures []Fixture, awayTeamLastPlayedFixtures []Fixture) error {
	homeTeam := TeamsGetWithName(f.homeTeam)
	awayTeam := TeamsGetWithName(f.awayTeam)

	// Additional strength given to the home team (home factor)
	homeStadiumStrength := HOME_BONUS_FACTOR * (homeTeam.HomeFactor / 10)

	// Consider home team recent form
	homeTeamRawFormBonus, err := calculateFormBonus(f.homeTeam, homeTeamLastPlayedFixtures)
	if err != nil {
		return err
	}

	// Consider away team recent form
	awayTeamRawFormBonus, err := calculateFormBonus(f.awayTeam, awayTeamLastPlayedFixtures)
	if err != nil {
		return err
	}

	homeTeamFormBonus := (RECENT_FORM_BONUS_FACTOR / 2.0) + (RECENT_FORM_BONUS_FACTOR * homeTeamRawFormBonus)
	awayTeamFormBonus := (RECENT_FORM_BONUS_FACTOR / 2.0) + (RECENT_FORM_BONUS_FACTOR * awayTeamRawFormBonus)

	// Home team strength factors
	homeAttackStrength := 1.5*homeTeam.Attack + homeTeam.Midfield
	homeDefenseStrength := 1.5*homeTeam.Defense + homeTeam.Midfield

	// Away team strength factors
	awayDefenseStrength := 1.5*awayTeam.Defense + awayTeam.Midfield
	awayAttackStrength := 1.5*awayTeam.Attack + awayTeam.Midfield

	// Calculate home/away strength without bonuses
	homeRawStrength := homeAttackStrength / (1 + awayDefenseStrength/homeAttackStrength)
	awayRawStrength := awayAttackStrength / (1 + homeDefenseStrength/awayAttackStrength)

	// Final non-attentuated strength of each team for this match
	homeStrength := homeStadiumStrength * homeTeamFormBonus * homeRawStrength
	awayStrength := awayTeamFormBonus * awayRawStrength

	// Attenuate strengths by employing a log-based function
	homeLambda := attenuateStrength(homeStrength)
	awayLambda := attenuateStrength(awayStrength)

	// Generate final scores based on a poisson distribution
	f.homeTeamScore = poissonKnuth(homeLambda)
	f.awayTeamScore = poissonKnuth(awayLambda)

	//fmt.Printf("%s: %f \n", f.homeTeam, homeTeamFormBonus)
	//fmt.Printf("%s: %f \n\n", f.awayTeam, awayTeamFormBonus)
	//fmt.Printf("%s: %f -> %f\n", f.homeTeam, homeLambda, homeLambda)
	//fmt.Printf("%s: %f -> %f\n\n", f.awayTeam, awayLambda, awayLambda)

	f.played = true
	return nil
}

// https://www.johndcook.com/blog/2010/06/14/generating-poisson-random-values/
func poissonKnuth(lambda float64) int {
	if lambda <= 0 {
		return 0
	}

	L := math.Exp(-lambda)
	k := int(0)
	p := 1.0

	for p > L {
		k++
		u := rand.Float64()
		p *= u
	}

	return k - 1
}

func attenuateStrength(x float64) float64 {
	if x <= 0 {
		return 0
	}
	return math.Log(1 + LOG_ADJUST_FACTOR*x)
}

func calculateFormBonus(teamName string, teamAlreadyPlayedFixtures []Fixture) (float64, error) {
	// Last matches are analyzed and summed to this bonus. The maximum bonus is 1
	formBonus := 0.0

	if len(teamAlreadyPlayedFixtures) <= 0 {
		return formBonus, nil
	}

	// last match contributes 35%
	fixtureContribution, err := getFixtureFormBonusContribution(teamAlreadyPlayedFixtures[0], teamName, 0.35)
	if err != nil {
		return 0.0, err
	}
	formBonus += fixtureContribution

	if len(teamAlreadyPlayedFixtures) <= 1 {
		return formBonus, nil
	}

	// 2nd match contributes 20%
	fixtureContribution, err = getFixtureFormBonusContribution(teamAlreadyPlayedFixtures[1], teamName, 0.20)
	if err != nil {
		return 0.0, err
	}
	formBonus += fixtureContribution

	if len(teamAlreadyPlayedFixtures) <= 2 {
		return formBonus, nil
	}

	// 3rd match contributes 15%
	fixtureContribution, err = getFixtureFormBonusContribution(teamAlreadyPlayedFixtures[2], teamName, 0.15)
	if err != nil {
		return 0.0, err
	}
	formBonus += fixtureContribution

	if len(teamAlreadyPlayedFixtures) <= 3 {
		return formBonus, nil
	}

	// 4th match contributes 15%
	fixtureContribution, err = getFixtureFormBonusContribution(teamAlreadyPlayedFixtures[3], teamName, 0.15)
	if err != nil {
		return 0.0, err
	}
	formBonus += fixtureContribution

	if len(teamAlreadyPlayedFixtures) <= 4 {
		return formBonus, nil
	}

	// 5th match contributes 15%
	fixtureContribution, err = getFixtureFormBonusContribution(teamAlreadyPlayedFixtures[4], teamName, 0.15)
	if err != nil {
		return 0.0, err
	}
	formBonus += fixtureContribution

	return formBonus, nil
}

func getFixtureFormBonusContribution(fixture Fixture, teamName string, contributionPctg float64) (float64, error) {
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

func getTeamScoredAndConcededGoalsInFixture(teamName string, fixture Fixture) (int, int, error) {
	if fixture.homeTeam == teamName {
		return fixture.homeTeamScore, fixture.awayTeamScore, nil
	} else if fixture.awayTeam == teamName {
		return fixture.awayTeamScore, fixture.homeTeamScore, nil
	}
	return -1, -1, errors.New("Fixture does not have target team")
}
