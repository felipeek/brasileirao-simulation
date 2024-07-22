package main

import (
	"fmt"
	"sort"
)

type TeamStatistic struct {
	Name         string
	Matches      int64
	Points       int64
	Won          int64
	Drawn        int64
	Lost         int64
	GoalsFor     int64
	GoalsAgainst int64
	GoalsDiff    int64
}

type Standings struct {
	TeamStatistics []*TeamStatistic
}

func GenerateStandings(s *Schedule) Standings {
	teams := TeamsGet()
	standingsMap := make(map[string]*TeamStatistic)

	for _, team := range teams {
		teamStatistic := TeamStatistic{}
		teamStatistic.Name = team.Name
		standingsMap[team.Name] = &teamStatistic
	}

	for _, round := range s.rounds {
		for _, fixture := range round.fixtures {
			homeTeamStatistics := standingsMap[fixture.homeTeam]
			awayTeamStatistics := standingsMap[fixture.awayTeam]

			if fixture.homeTeamScore > fixture.awayTeamScore {
				homeTeamStatistics.Points += 3
				homeTeamStatistics.Won += 1
				awayTeamStatistics.Lost += 1
			} else if fixture.awayTeamScore > fixture.homeTeamScore {
				awayTeamStatistics.Points += 3
				awayTeamStatistics.Won += 1
				homeTeamStatistics.Lost += 1
			} else {
				homeTeamStatistics.Points += 1
				awayTeamStatistics.Points += 1
				homeTeamStatistics.Drawn += 1
				awayTeamStatistics.Drawn += 1
			}

			homeTeamStatistics.Matches += 1
			awayTeamStatistics.Matches += 1
			homeTeamStatistics.GoalsFor += fixture.homeTeamScore
			homeTeamStatistics.GoalsAgainst += fixture.awayTeamScore
			awayTeamStatistics.GoalsFor += fixture.awayTeamScore
			awayTeamStatistics.GoalsAgainst += fixture.homeTeamScore
			homeTeamStatistics.GoalsDiff += (fixture.homeTeamScore - fixture.awayTeamScore)
			awayTeamStatistics.GoalsDiff += (fixture.awayTeamScore - fixture.homeTeamScore)
		}
	}

	teamStatistics := []*TeamStatistic{}
	for _, teamStatistic := range standingsMap {
		teamStatistics = append(teamStatistics, teamStatistic)
	}

	sort.Slice(teamStatistics, func(i, j int) bool {
		// todo if points == points, consider wins/goals
		return teamStatistics[i].Points > teamStatistics[j].Points
	})

	standings := Standings{}
	standings.TeamStatistics = teamStatistics
	return standings
}

func (s *Standings) Print() {
	headerFormat := "%-6s %-20s %-8s %-6s %-6s %-6s %-6s %-9s %-12s %-9s\n"
	dataFormat := "%-6d %-20s %-8d %-6d %-6d %-6d %-6d %-9d %-12d %-9d\n"
	fmt.Printf(headerFormat, "Rank", "Team", "Matches", "Points", "Won", "Drawn", "Lost", "GoalsFor", "GoalsAgainst", "GoalsDiff")

	for i, team := range s.TeamStatistics {
		fmt.Printf(dataFormat, i+1, team.Name, team.Matches, team.Points, team.Won, team.Drawn, team.Lost, team.GoalsFor, team.GoalsAgainst, team.GoalsDiff)
	}
}
