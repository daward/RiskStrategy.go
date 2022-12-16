package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/pkg/profile"
)

type ConquestScore struct {
	Country  string  `json:"country"`
	Score    float64 `json:"score"`
	Distance int     `json:"distance"`
}

type PositionResult struct {
	InitialTerritories []string        `json:"startingPosition"`
	Score              float64         `json:"finalScore"`
	ConqueringPath     []ConquestScore `json:"path"`
}

func (this *PositionResult) evaluate(depth int, risk *RiskBoard, continents *ContinentSet) {

	// this will be the set of paths we're currently processing
	paths := make([]*Path, 1)
	// this will be the set of paths to examine next
	nextPaths := make([]*Path, 0)

	// the risk board does accumulate data from previous runs, reset the state
	risk.clearPaths()

	// build the initial territory holdings
	p := risk.buildTerritoryPath(continents, this.InitialTerritories)
	paths[0] = p

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
		for _, currentPath := range paths {

			// make sure its worth looking at
			if !currentPath.isComplete() && !currentPath.isRedundant {
				// then put all the next paths in the queue
				nextPaths = append(nextPaths, currentPath.expand()...)
			}
		}
		paths = make([]*Path, min(depth, len(nextPaths)))

		// we've gone through all the paths for the current distance from our origin
		// but we probably have more paths we've been saving away to look at
		if len(nextPaths) != 0 {

			// in order to not look at lesser paths, lets sort what we have
			sort.Slice(nextPaths, func(i, j int) bool {
				if nextPaths[i].TotalScore > nextPaths[j].TotalScore {
					return true
				}
				return false
			})

			// grab the best path for record keeping
			bestPath = nextPaths[0]

			// then make next paths our new paths
			copy(paths, nextPaths)
			nextPaths = make([]*Path, 0)
		}
	}
	this.computeTotalScore(bestPath, continents)

}

func (this *PositionResult) computeTotalScore(bestPath *Path, continents *ContinentSet) {
	this.Score = bestPath.TotalScore / float64(42-len(this.InitialTerritories))
	conquests := bestPath.conquests()
	this.ConqueringPath = make([]ConquestScore, 42)
	this.ConqueringPath[0] = ConquestScore{
		Country: conquests[0],
		Score:   0,
	}

	bestPath.Map.clearPaths()

	p := bestPath.Map.startPath(continents, conquests[0])

	for i, val := range conquests[1:] {
		p = p.conquer(bestPath.Map.countryIndex[val])
		this.ConqueringPath[i+1] = ConquestScore{
			Country:  val,
			Score:    p.TotalScore,
			Distance: p.distance,
		}
	}
}

func checkAllTerritories(risk *RiskBoard, continents *ContinentSet) []*PositionResult {
	i := 0
	results := make([]*PositionResult, len(risk.countryLookup))

	// iterate through all the countries, evaluating the best path
	// to conquer for each
	for _, country := range risk.countryLookup {
		results[i] = &PositionResult{
			InitialTerritories: []string{country},
		}
		results[i].evaluate(10000, risk, continents)
		i++
	}
	return results
}

func checkSpecificTerritorySet(territories []string, risk *RiskBoard, continents *ContinentSet) []*PositionResult {

	results := make([]*PositionResult, 1)
	results[0] = &PositionResult{
		InitialTerritories: territories,
	}
	results[0].evaluate(1000, risk, continents)
	return results
}

func main() {
	defer profile.Start(profile.ProfilePath(".")).Stop()

	// build the risk board, indexing the countries as we go
	risk := riskboard()

	// then initialize the continents so we can determine bonuses later
	continents := NewContinentSet(risk)

	//results := checkSpecificTerritorySet([]string{"Brazil"}, risk, continents)

	results := checkAllTerritories(risk, continents)

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
