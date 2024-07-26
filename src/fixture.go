package main

import (
	"math"
	"math/rand"
)

type Fixture struct {
	homeTeam      string
	awayTeam      string
	homeTeamScore int64
	awayTeamScore int64
	played        bool
}

const (
	LOG_ADJUST_FACTOR = 0.5 // the smaller, the more 'balanced' the results
	BONUS_HOME        = 2.0
)

func (f *Fixture) Play() {
	homeTeam := TeamsGetWithName(f.homeTeam)
	awayTeam := TeamsGetWithName(f.awayTeam)

	// Additional strength given to the home team (home factor)
	homeStadiumStrength := BONUS_HOME * (float64(homeTeam.HomeFactor) / 10)

	// Home team strength factors
	homeAttackStrength := 1.5*homeTeam.Attack + homeTeam.Midfield
	homeDefenseStrength := 1.5*homeTeam.Defense + homeTeam.Midfield

	// Away team strength factors
	awayDefenseStrength := 1.5*awayTeam.Defense + awayTeam.Midfield
	awayAttackStrength := 1.5*awayTeam.Attack + awayTeam.Midfield

	// Final non-attentuated strength of each team for this match
	homeStrength := homeStadiumStrength * (homeAttackStrength / (1 + awayDefenseStrength/homeAttackStrength))
	awayStrength := awayAttackStrength / (1 + homeDefenseStrength/awayAttackStrength)

	// Attenuate strengths by employing a log-based function
	homeLambda := attenuateStrength(homeStrength)
	awayLambda := attenuateStrength(awayStrength)

	// Generate final scores based on a poisson distribution
	f.homeTeamScore = poissonKnuth(homeLambda)
	f.awayTeamScore = poissonKnuth(awayLambda)

	//fmt.Printf("%s: %f -> %f\n", f.homeTeam, homeLambda, homeLambda)
	//fmt.Printf("%s: %f -> %f\n\n", f.awayTeam, awayLambda, awayLambda)

	f.played = true
}

// https://www.johndcook.com/blog/2010/06/14/generating-poisson-random-values/
func poissonKnuth(lambda float64) int64 {
	L := math.Exp(-lambda)
	k := int64(0)
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
