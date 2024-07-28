package simulation

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/bit101/go-ansi"
	"github.com/felipeek/brasileirao-simulation/internal/util"
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
	TeamStatistics         []*TeamStatistic
	PreviousTeamStatistics []*TeamStatistic
}

func GenerateStandings(s *Schedule) Standings {
	standings := Standings{}
	standings.TeamStatistics = generateTeamStatisticsUntilRound(s, s.currentRoundIdx)
	if s.currentRoundIdx > 0 {
		standings.PreviousTeamStatistics = generateTeamStatisticsUntilRound(s, s.currentRoundIdx-1)
	} else {
		standings.PreviousTeamStatistics = nil
	}
	return standings
}

func generateTeamStatisticsUntilRound(s *Schedule, roundIdx int) []*TeamStatistic {
	standingsMap := fillStandingsMapUntilRound(s, roundIdx)

	teamStatistics := []*TeamStatistic{}
	for _, teamStatistic := range standingsMap {
		teamStatistics = append(teamStatistics, teamStatistic)
	}

	sort.Slice(teamStatistics, func(i, j int) bool {
		return tieBreak(*teamStatistics[i], *teamStatistics[j], s)
	})

	return teamStatistics
}

func fillStandingsMapUntilRound(s *Schedule, roundIdx int) map[string]*TeamStatistic {
	teams := TeamsGet()
	standingsMap := make(map[string]*TeamStatistic)

	for _, team := range teams {
		teamStatistic := TeamStatistic{}
		teamStatistic.Name = team.Name
		standingsMap[team.Name] = &teamStatistic
	}

	for i := 0; i <= roundIdx; i++ {
		round := s.rounds[i]

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

	return standingsMap
}

func tieBreak(t1, t2 TeamStatistic, s *Schedule) bool {
	pointsDiff := t1.Points - t2.Points
	if pointsDiff > 0 {
		return true
	} else if pointsDiff < 0 {
		return false
	}

	goalsDiff := t1.GoalsDiff - t2.GoalsDiff
	if goalsDiff > 0 {
		return true
	} else if goalsDiff < 0 {
		return false
	}

	goalsForDiff := t1.GoalsFor - t2.GoalsFor
	if goalsForDiff > 0 {
		return true
	} else if goalsForDiff < 0 {
		return false
	}

	iScore, jScore := summedH2HResults(s, t1.Name, t2.Name)
	h2hDiff := iScore - jScore
	if h2hDiff > 0 {
		return true
	} else if h2hDiff < 0 {
		return false
	}

	return rand.Int()%2 == 0
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

func (s *Standings) Print(enableTerminalColors bool) error {
	headerFormat := "%-6s %-20s %-8s %-6s %-6s %-6s %-6s %-9s %-12s %-9s %-12s %-6s\n"
	fmt.Printf(headerFormat, "Rank", "Team", "Matches", "Points", "Won", "Drawn", "Lost", "GoalsFor", "GoalsAgainst", "GoalsDiff", "RecentForm", "Change")

	for i, teamStatistics := range s.TeamStatistics {
		team := TeamsGetWithName(teamStatistics.Name)

		printStandingsRank(enableTerminalColors, i+1)
		fmt.Printf(" ")
		printStandingsTeamName(enableTerminalColors, i+1, teamStatistics.Name)
		fmt.Printf(" ")
		printStandingsMatches(enableTerminalColors, teamStatistics.Matches)
		fmt.Printf(" ")
		printStandingsPoints(enableTerminalColors, teamStatistics.Points)
		fmt.Printf(" ")
		printStandingsWon(enableTerminalColors, teamStatistics.Won)
		fmt.Printf(" ")
		printStandingsDrawn(enableTerminalColors, teamStatistics.Drawn)
		fmt.Printf(" ")
		printStandingsLost(enableTerminalColors, teamStatistics.Lost)
		fmt.Printf(" ")
		printStandingsGoalsFor(enableTerminalColors, teamStatistics.GoalsFor)
		fmt.Printf(" ")
		printStandingsGoalsAgainst(enableTerminalColors, teamStatistics.GoalsAgainst)
		fmt.Printf(" ")
		printStandingsGoalsDiff(enableTerminalColors, teamStatistics.GoalsDiff)
		fmt.Printf(" ")
		printStandingsRecentForm(enableTerminalColors, teamStatistics.Name, team.DynamicAttributes.LastFixtures)
		fmt.Printf(" ")
		teamPositionChange, err := getTeamPositionChange(teamStatistics.Name, s.TeamStatistics, s.PreviousTeamStatistics)
		if err != nil {
			return err
		}
		printStandingsChanges(enableTerminalColors, teamPositionChange)
		fmt.Println()
	}

	return nil
}

func printStandingsRank(enableTerminalColors bool, rank int) {
	format := "%-6d"
	if !enableTerminalColors {
		fmt.Printf(format, rank)
	} else {
		ansi.Printf(getRankPrintColor(rank), format, rank)
	}
}

func printStandingsTeamName(enableTerminalColors bool, rank int, teamName string) {
	format := "%-20s"
	if !enableTerminalColors {
		fmt.Printf(format, teamName)
	} else {
		ansi.Printf(getRankPrintColor(rank), format, teamName)
	}
}

func printStandingsMatches(enableTerminalColors bool, matches int) {
	format := "%-8d"
	if !enableTerminalColors {
		fmt.Printf(format, matches)
	} else {
		ansi.Printf(ansi.White, format, matches)
	}
}

func printStandingsPoints(enableTerminalColors bool, points int) {
	format := "%-6d"
	if !enableTerminalColors {
		fmt.Printf(format, points)
	} else {
		ansi.Printf(ansi.BoldWhite, format, points)
	}
}

func printStandingsWon(enableTerminalColors bool, won int) {
	format := "%-6d"
	if !enableTerminalColors {
		fmt.Printf(format, won)
	} else {
		ansi.Printf(ansi.White, format, won)
	}
}

func printStandingsDrawn(enableTerminalColors bool, drawn int) {
	format := "%-6d"
	if !enableTerminalColors {
		fmt.Printf(format, drawn)
	} else {
		ansi.Printf(ansi.White, format, drawn)
	}
}

func printStandingsLost(enableTerminalColors bool, lost int) {
	format := "%-6d"
	if !enableTerminalColors {
		fmt.Printf(format, lost)
	} else {
		ansi.Printf(ansi.White, format, lost)
	}
}

func printStandingsGoalsFor(enableTerminalColors bool, goalsFor int) {
	format := "%-9d"
	if !enableTerminalColors {
		fmt.Printf(format, goalsFor)
	} else {
		ansi.Printf(ansi.White, format, goalsFor)
	}
}

func printStandingsGoalsAgainst(enableTerminalColors bool, goalsAgainst int) {
	format := "%-12d"
	if !enableTerminalColors {
		fmt.Printf(format, goalsAgainst)
	} else {
		ansi.Printf(ansi.White, format, goalsAgainst)
	}
}

func printStandingsGoalsDiff(enableTerminalColors bool, goalsDiff int) {
	format := "%-9d"
	if !enableTerminalColors {
		fmt.Printf(format, goalsDiff)
	} else {
		ansi.Printf(ansi.White, format, goalsDiff)
	}
}

func printStandingsRecentForm(enableTerminalColors bool, teamName string, recentMatches []*Fixture) {
	matchChar := "● "
	noMatchChar := "─ "
	if !enableTerminalColors {
		for i := 0; i < 5; i++ {
			fmt.Print(matchChar)
		}
	} else {
		for i := 4; i >= 0; i-- {
			if i >= len(recentMatches) {
				ansi.Printf(ansi.BoldWhite, noMatchChar)
				continue
			}
			fixture := recentMatches[i]
			goalDiff := fixture.homeTeamScore - fixture.awayTeamScore
			if teamName == fixture.awayTeam {
				goalDiff = -goalDiff
			}
			if goalDiff > 0 {
				ansi.Printf(ansi.BoldGreen, matchChar)
			} else if goalDiff < 0 {
				ansi.Printf(ansi.BoldRed, matchChar)
			} else {
				ansi.Printf(ansi.BoldWhite, matchChar)
			}
		}
	}
	fmt.Print("  ")
}

func getTeamPositionChange(teamName string, currentStatistics []*TeamStatistic, previousStatistics []*TeamStatistic) (int, error) {
	if previousStatistics == nil {
		return 0, nil
	}

	currentPosition := -1
	previousPosition := -1

	for i, teamStatistic := range currentStatistics {
		if teamStatistic.Name == teamName {
			currentPosition = i
			break
		}
	}

	for i, teamStatistic := range previousStatistics {
		if teamStatistic.Name == teamName {
			previousPosition = i
			break
		}
	}

	if currentPosition == -1 || previousPosition == -1 {
		return 0, fmt.Errorf("Team [%s] was not found in teamStatistics", teamName)
	}

	return previousPosition - currentPosition, nil
}

func printStandingsChanges(enableTerminalColors bool, teamPositionChange int) {
	format := "%c %-4d"
	upChar := '↑'
	downChar := '↓'
	noChangeChar := '─'
	if !enableTerminalColors {
		if teamPositionChange > 0 {
			fmt.Printf(format, upChar, teamPositionChange)
		} else if teamPositionChange < 0 {
			fmt.Printf(format, downChar, util.UtilIntAbs(teamPositionChange))
		} else {
			fmt.Printf(format, noChangeChar, teamPositionChange)
		}
	} else {
		if teamPositionChange > 0 {
			ansi.Printf(ansi.BoldGreen, format, upChar, teamPositionChange)
		} else if teamPositionChange < 0 {
			ansi.Printf(ansi.BoldRed, format, downChar, util.UtilIntAbs(teamPositionChange))
		} else {
			ansi.Printf(ansi.BoldWhite, format, noChangeChar, teamPositionChange)
		}
	}
}

func getRankPrintColor(rank int) ansi.AnsiColor {
	if rank <= 4 {
		return ansi.BoldCyan
	} else if rank <= 6 {
		return ansi.BoldBlue
	} else if rank <= 12 {
		return ansi.BoldYellow
	} else if rank <= 16 {
		return ansi.BoldWhite
	} else {
		return ansi.BoldRed
	}
}
