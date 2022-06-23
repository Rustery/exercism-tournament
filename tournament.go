package tournament

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
)

type TeamResult struct {
	mp int
	w  int
	d  int
	l  int
	p  int
}

var teams []string
var results map[string]*TeamResult

func getTeamResult(name string) *TeamResult {
	r, ok := results[name]
	if !ok {
		results[name] = &TeamResult{}
		teams = append(teams, name)
		r = results[name]
	}
	return r
}

func (team *TeamResult) win() {
	team.mp++
	team.w++
	team.p += 3
}

func (team *TeamResult) loss() {
	team.mp++
	team.l++
}

func (team *TeamResult) draw() {
	team.mp++
	team.d++
	team.p += 1
}

func printResults(writer io.Writer) error {
	sort.SliceStable(teams, func(i, j int) bool {
		if results[teams[i]].p == results[teams[j]].p {
			return teams[i] < teams[j]
		}
		return results[teams[i]].p > results[teams[j]].p
	})
	fmt.Fprintf(writer, "%-31s| %2s | %2s | %2s | %2s | %2s\n", "Team", "MP", "W", "D", "L", "P")
	for _, team := range teams {
		result := results[team]
		_, err := fmt.Fprintf(writer, "%-31s| %2d | %2d | %2d | %2d | %2d\n", team, result.mp, result.w, result.d, result.l, result.p)
		if err != nil {
			return err
		}
	}
	return nil
}

func Tally(reader io.Reader, writer io.Writer) error {
	results = make(map[string]*TeamResult)
	teams = make([]string, 0)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {

		var teamAName, teamBName, outcome string

		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		data := strings.Split(line, ";")
		if len(data) < 3 {
			return errors.New("bad params")
		}

		teamAName, teamBName, outcome = data[0], data[1], data[2]
		teamA := getTeamResult(teamAName)
		teamB := getTeamResult(teamBName)

		switch outcome {
		case "win":
			teamA.win()
			teamB.loss()
		case "draw":
			teamA.draw()
			teamB.draw()
		case "loss":
			teamB.win()
			teamA.loss()
		default:
			return errors.New("bad params")
		}
	}

	err := printResults(writer)
	return err
}
