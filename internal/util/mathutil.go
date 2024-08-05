package util

import (
	"math"
	"math/rand"
)

// Box-Muller transform.
func RandomValueFromNormalDistribution(center, stddev float64) float64 {
	u1 := rand.Float64()
	u2 := rand.Float64()
	z0 := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2.0*math.Pi*u2)
	return center + z0*stddev
}

// https://www.johndcook.com/blog/2010/06/14/generating-poisson-random-values/
func PoissonKnuth(lambda float64) int {
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

func AttenuateStrength(strength float64) float64 {
	// the bigger this number, the more the curve will be 'flatten', making it less steep
	// this means that the strength will have less influence and the output will be more balanced
	flattenFactor := 3.0

	// the adjuster serves as a simple linear operation to adjust the output average
	// i.e., assuming strength == 1 (which is a good pick since this is independent from the flattenFactor),
	// then the result will be strength * adjuster
	adjuster := 0.666
	return math.Pow(strength, 1.0/flattenFactor) * adjuster
}

// Given a contribution C in [0,10] and an impact I on the overall simulation, finds a multiplier M
// For example, if I=0.1487, then M=0.5 for C=0 and M=2.0 for C=10,
// which means that the max contribution will multiply by 2 (and minimum contribution divide by half) the strength
// Another example, if I=0.019, then M=0.9102 for C=0 and M=1.0987
// The following relation will be always true: max(M) = 1 / min(M)
// To test MAX multiplier in calculator: (1+I)^(10−5)
// To test MIN multiplier in calculator: (1+I)^(0−5)
func GetMultiplierFromContributionFactor(contribution, impact float64) float64 {
	return math.Pow(1+impact, contribution-5)
}
