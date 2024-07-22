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
	LOG_ADJUST_FACTOR = 10 // the bigger, the more 'balanced' the results
	BONUS_HOME        = 1.3
)

func (f *Fixture) Play() {
	homeTeam := TeamsGetWithName(f.homeTeam)
	awayTeam := TeamsGetWithName(f.awayTeam)

	homeLambda := BONUS_HOME * (1.5*float64(homeTeam.Attack) - float64(awayTeam.Defense) + 0.5*(float64(homeTeam.Midfield)-float64(awayTeam.Midfield)))
	awayLambda := 1.5*float64(awayTeam.Attack) - float64(homeTeam.Defense) + 0.5*(float64(awayTeam.Midfield)-float64(homeTeam.Midfield))

	f.homeTeamScore = MaxInt64(0, poissonKnuth(CustomLog(homeLambda)))
	f.awayTeamScore = MaxInt64(0, poissonKnuth(CustomLog(awayLambda)))

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

func CustomLog(x float64) float64 {
	if x <= 0 {
		// Logarithm not defined for non-positive values
		return math.NaN() // Returns "not a number"
	}
	// Calculate 0.5 * log_base(sqrt(5)) of x + 1
	return 0.5*math.Log(x)/math.Log(LOG_ADJUST_FACTOR) + 1
}
