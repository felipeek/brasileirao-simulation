package simulation

import (
	"fmt"
	"math/rand"
	"sort"
)

type TeamStatistic struct {
	Name         string
	Matches      int
	Points       int
	Won          int
	Drawn        int
	Lost         int
	GoalsFor     int
	GoalsAgainst int
	GoalsDiff    int
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
			if !fixture.played {
				continue
			}
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
		pointsDiff := teamStatistics[i].Points - teamStatistics[j].Points
		if pointsDiff > 0 {
			return true
		} else if pointsDiff < 0 {
			return false
		}

		goalsDiff := teamStatistics[i].GoalsDiff - teamStatistics[j].GoalsDiff
		if goalsDiff > 0 {
			return true
		} else if goalsDiff < 0 {
			return false
		}

		goalsForDiff := teamStatistics[i].GoalsFor - teamStatistics[j].GoalsFor
		if goalsForDiff > 0 {
			return true
		} else if goalsForDiff < 0 {
			return false
		}

		iScore, jScore := summedH2HResults(s, teamStatistics[i].Name, teamStatistics[j].Name)
		h2hDiff := iScore - jScore
		if h2hDiff > 0 {
			return true
		} else if h2hDiff < 0 {
			return false
		}

		return rand.Int()%2 == 0
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

func summedH2HResults(s *Schedule, team1Name string, team2Name string) (int, int) {
	team1SummedScore := 0
	team2SummedScore := 0

	for _, r := range s.rounds {
		for _, f := range r.fixtures {
			if f.played && f.homeTeam == team1Name && f.awayTeam == team2Name {
				team1SummedScore += f.homeTeamScore
				team2SummedScore += f.awayTeamScore
			} else if f.played && f.homeTeam == team2Name && f.awayTeam == team1Name {
				team2SummedScore += f.homeTeamScore
				team1SummedScore += f.awayTeamScore
			}
		}
	}

	return team1SummedScore, team2SummedScore
}
