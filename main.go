package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/pkg/profile"
)

type ConquestScore struct {
	Country string  `json:"country"`
	Score   float64 `json:"score"`
}

type PositionResult struct {
	InitialTerritories []string        `json:"startingPosition"`
	Score              float64         `json:"finalScore"`
	ConqueringPath     []ConquestScore `json:"path"`
}

func (this *PositionResult) evaluate(depth int, risk *RiskBoard, continents *ContinentSet) {

	// this will be the set of paths we're currently processing
	paths := make([]*Path, 1)
	// this will be the set
	nextPaths := make([]*Path, 0)

	risk.clearPaths()

	p := Path{
		Map:               risk,
		continents:        continents,
		Borders:           &TerritorySet{data: 0},
		BorderTerritories: &TerritorySet{data: 0},
		TotalScore:        0,
		Conquests:         this.InitialTerritories,
		IndexedConquests:  risk.index(this.InitialTerritories),
	}
	p.detectBorders()
	p.setTotalScore()
	paths[0] = &p

	// since this is ultimately a breadth first search, we'll be using our share of queues
	pop := func() *Path {
		firstElement := paths[0]
		paths = append(paths[:0], paths[1:]...)
		return firstElement
	}

	min := func(i, j int) int {
		if i < j {
			return i
		}
		return j
	}

	var bestPath *Path

	for {
		if len(paths) == 0 {
			break
		}

		// get the first path in the queue
		currentPath := pop()

		// make sure its worth looking at
		if !currentPath.isComplete() && !currentPath.isRedundant {
			// then put all the next paths in the queue
			nextPaths = append(nextPaths, currentPath.expand()...)
		}

		// we've gone through all the paths for the current distance from our origin
		// but we probably have more paths we've been saving away to look at
		if len(paths) == 0 && len(nextPaths) != 0 {

			// in order to not look at lesser paths, lets sort what we have
			sort.Slice(nextPaths, func(i, j int) bool {
				if nextPaths[i].TotalScore > nextPaths[j].TotalScore {
					return true
				}
				return false
			})

			// grab the best path for record keeping
			bestPath = nextPaths[0]

			paths = make([]*Path, min(depth, len(nextPaths)))
			// then make next paths our new paths
			copy(paths, nextPaths)
			nextPaths = make([]*Path, 0)
		}
	}
	this.Score = bestPath.TotalScore / float64(42-len(this.InitialTerritories))

	this.ConqueringPath = make([]ConquestScore, 42)
	this.ConqueringPath[0] = ConquestScore{
		Country: bestPath.Conquests[0],
		Score:   0,
	}

	for i, val := range bestPath.Conquests[1:] {
		p = *p.conquer(risk.countryIndex[val])
		this.ConqueringPath[i+1] = ConquestScore{
			Country: val,
			Score:   p.TotalScore,
		}
	}
}

func main() {
	defer profile.Start(profile.ProfilePath(".")).Stop()

	// build the risk board, indexing the countries as we go
	risk := riskboard()

	// then initialize the continents so we can determine bonuses later
	continents := NewContinentSet(risk)

	// set up an array to put results into
	results := make([]*PositionResult, len(risk.countryLookup))

	// iterate through all the countries, evaluating the best path
	// to conquer for each
	i := 0
	for _, country := range risk.countryLookup {
		results[i] = &PositionResult{
			InitialTerritories: []string{country},
		}
		results[i].evaluate(1000, risk, continents)
		i++
	}

	// it would be useful to sort the final outcomes
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	jsonStr, _ := json.Marshal(results)
	err := ioutil.WriteFile("results.json", jsonStr, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}
