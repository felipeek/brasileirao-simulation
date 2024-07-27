package main

import (
	"fmt"
	"math/rand"
	"slices"
)

type Round struct {
	fixtures []*Fixture
}

type Schedule struct {
	currentRoundIdx int
	nextRoundIdx    int
	finished        bool
	rounds          []*Round
}

func GenerateSchedule(teams map[string]*Team) (Schedule, error) {
	if len(teams)%2 != 0 {
		return Schedule{}, fmt.Errorf("Number of teams must be pair")
	}

	schedule := Schedule{}
	roundRobinTeams := []string{}

	schedule.currentRoundIdx = -1
	schedule.nextRoundIdx = 0
	schedule.finished = false

	// Create a slice containing all available teams
	// It will be used to construct the schedule
	for teamName, _ := range teams {
		roundRobinTeams = append(roundRobinTeams, teamName)
	}

	// Randomize array (simple algorithm, not very good randomization)
	for i := 0; i < 2*len(roundRobinTeams); i++ {
		r1 := rand.Int() % len(roundRobinTeams)
		r2 := rand.Int() % len(roundRobinTeams)

		cached := roundRobinTeams[r1]
		roundRobinTeams[r1] = roundRobinTeams[r2]
		roundRobinTeams[r2] = cached
	}

	// Build the schedule by performing round-robins over the list of teams, while keeping the first one static,
	// and defining the matches as first x last, 2nd first x 2nd last, etc
	// This guarantees that all teams play only a single game per round, and that all games only appear once in the schedule
	for i := 0; i < len(roundRobinTeams)-1; i++ {
		round := Round{}

		for j := 0; j < len(roundRobinTeams)/2; j++ {
			fixture := Fixture{roundRobinTeams[j], roundRobinTeams[len(roundRobinTeams)-1-j], -1, -1, false}
			round.fixtures = append(round.fixtures, &fixture)
		}

		schedule.rounds = append(schedule.rounds, &round)

		cached := roundRobinTeams[len(roundRobinTeams)-1]
		current := ""
		for j := 0; j < len(roundRobinTeams)-1; j++ {
			current = cached
			cached = roundRobinTeams[j+1]
			roundRobinTeams[j+1] = current
		}
	}

	homeAwayCountMap := make(map[string]int64)

	for teamName, _ := range teams {
		homeAwayCountMap[teamName] = 0
	}

	// Here we make an adjustment to rebalance the away/home teams. The goal is to:
	// 1. Make sure that all teams have a balanced number of home vs away games
	// 2. Make sure that teams don't have many games in a row in which they are always either the home team or away team
	// For that, we use a simply algorithm that stores a counter indicating how many times each team appeared as home and away
	// And iterate over the rounds trying to balance everything out, switching the away/home team of all fixtures when necessary
	for _, round := range schedule.rounds {
		for _, fixture := range round.fixtures {
			firstTeamCount := homeAwayCountMap[fixture.homeTeam]
			secondTeamCount := homeAwayCountMap[fixture.awayTeam]

			if UtilInt64Abs(firstTeamCount) > UtilInt64Abs(secondTeamCount) {
				if firstTeamCount > 0 {
					cached := fixture.homeTeam
					fixture.homeTeam = fixture.awayTeam
					fixture.awayTeam = cached
				}
			} else {
				if secondTeamCount < 0 {
					cached := fixture.homeTeam
					fixture.homeTeam = fixture.awayTeam
					fixture.awayTeam = cached
				}
			}

			homeAwayCountMap[fixture.homeTeam] = homeAwayCountMap[fixture.homeTeam] + 1
			homeAwayCountMap[fixture.awayTeam] = homeAwayCountMap[fixture.awayTeam] - 1
		}
	}

	// Now that everything is ready, we can simply duplicate the current schedule, to create the other 'half-season'
	numHalfRounds := len(schedule.rounds)
	for i := 0; i < numHalfRounds; i++ {
		existingRound := schedule.rounds[i]
		counterpartRound := Round{}

		for _, fixture := range existingRound.fixtures {
			counterpartFixture := Fixture{fixture.awayTeam, fixture.homeTeam, -1, -1, false}
			counterpartRound.fixtures = append(counterpartRound.fixtures, &counterpartFixture)
		}

		schedule.rounds = append(schedule.rounds, &counterpartRound)
	}

	//for teamName, _ := range teams {
	//	fmt.Printf("%s: %d\n", teamName, homeAwayCountMap[teamName])
	//}

	return schedule, nil
}

func (r *Round) PlayFixturesOfRound(lastRounds []*Round) error {
	for _, fixture := range r.fixtures {
		err := fixture.Play()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Schedule) PlayAllFixtures() error {
	for _, round := range s.rounds {
		for _, fixture := range round.fixtures {
			if !fixture.played {
				err := fixture.Play()
				if err != nil {
					return err
				}
			}
		}
	}

	s.currentRoundIdx = len(s.rounds) - 1
	s.nextRoundIdx = -1
	s.finished = true
	return nil
}

func (s *Schedule) PlayNextRoundFixtures() error {
	if s.finished {
		return nil
	}

	round := s.rounds[s.nextRoundIdx]
	err := round.PlayFixturesOfRound(s.getPlayedRounds())
	if err != nil {
		return err
	}

	s.currentRoundIdx += 1
	s.nextRoundIdx += 1
	if s.nextRoundIdx == len(s.rounds) {
		s.nextRoundIdx = -1
		s.finished = true
	}
	return nil
}

func (r *Round) Print() {
	for _, fixture := range r.fixtures {
		fmt.Printf("\t%s %d x %d %s\n", fixture.homeTeam, fixture.homeTeamScore, fixture.awayTeamScore, fixture.awayTeam)
	}
}

func (s *Schedule) Print() {
	for i, round := range s.rounds {
		fmt.Printf("Round [%d]\n", i+1)
		round.Print()
	}
}

func (s *Schedule) PrintLastPlayedRound() {
	if s.currentRoundIdx >= 0 {
		fmt.Printf("Round [%d]\n", s.currentRoundIdx+1)
		round := s.rounds[s.currentRoundIdx]
		round.Print()
	}
}

func getTeamLastFixtures(teamName string, lastRounds []*Round) []Fixture {
	lastFixtures := []Fixture{}
	for _, r := range lastRounds {
		for _, f := range r.fixtures {
			if f.homeTeam == teamName || f.awayTeam == teamName {
				lastFixtures = append(lastFixtures, *f)
				break
			}
		}
	}

	return lastFixtures
}

func (s *Schedule) getPlayedRounds() []*Round {
	if s.currentRoundIdx < 0 {
		return []*Round{}
	}

	// Create a new slice to avoid modifying the original
	lastPlayedRounds := make([]*Round, s.currentRoundIdx+1)
	copy(lastPlayedRounds, s.rounds[0:s.currentRoundIdx+1])

	// Reverse the new slice
	slices.Reverse(lastPlayedRounds)

	return lastPlayedRounds
}
