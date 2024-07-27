package simulation

import (
	"math"
	"math/rand"
)

// Box-Muller transform.
func SimUtilRandomValueFromNormalDistribution(center, stddev float64) float64 {
	u1 := rand.Float64()
	u2 := rand.Float64()
	z0 := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2.0*math.Pi*u2)
	return center + z0*stddev
}

// https://www.johndcook.com/blog/2010/06/14/generating-poisson-random-values/
func SimUtilPoissonKnuth(lambda float64) int {
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

func SimUtilAttenuateStrength(x float64) float64 {
	if x <= 0 {
		return 0
	}
	return math.Log(1 + LOG_ADJUST_FACTOR*x)
}

// Given a contribution C in [0,10] and an impact I on the overall simulation, finds a multiplier M
// For example, if I=0.1487, then M=0.5 for C=0 and M=2.0 for C=10,
// which means that the max contribution will multiply by 2 (and minimum contribution divide by half) the strength
// Another example, if I=0.019, then M=0.9102 for C=0 and M=1.0987
// The following relation will be always true: max(M) = 1 / min(M)
// To test MAX multiplier in calculator: (1+I)^(10−5)
// To test MIN multiplier in calculator: (1+I)^(0−5)
func SimUtilGetMultiplierFromContributionFactor(contribution, impact float64) float64 {
	return math.Pow(1+impact, contribution-5)
}
