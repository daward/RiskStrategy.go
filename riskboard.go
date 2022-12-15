package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/yourbasic/graph"
)

type RiskBoard struct {
	board         graph.Iterator
	countryIndex  map[string]uint64
	countryLookup map[uint64]string
	paths         map[PathCoordinate]*Path // territory plus length of path
}

type PathCoordinate struct {
	territory uint64
	pathId    uint64
}

func riskboard() *RiskBoard {

	jsonFile, _ := os.Open("./mapping.json")
	// if we os.Open returns an error then handle it
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var arr [][]string
	json.Unmarshal([]byte(byteValue), &arr)

	m := make(map[string]uint64)
	m2 := make(map[uint64]string)

	g := graph.New(len(arr))
	for i := 0; i < len(arr); i++ {
		v1 := arr[i]
		if countryId1, ok := m[v1[0]]; !ok {
			countryId1 = uint64(len(m))
			m[v1[0]] = countryId1
			m2[countryId1] = v1[0]
		}

		if countryId2, ok := m[v1[1]]; !ok {
			countryId2 = uint64(len(m))
			m[v1[1]] = countryId2
			m2[countryId2] = v1[1]
		}

		g.AddBoth(int(m[v1[0]]), int(m[v1[1]]))
	}

	var retVal RiskBoard
	retVal.board = g
	retVal.countryIndex = m
	retVal.countryLookup = m2
	retVal.paths = make(map[PathCoordinate]*Path)

	return &retVal
}

func (this *RiskBoard) startPath(continents *ContinentSet, initialTerritory string) *Path {
	friendlyBorders := TerritorySet{data: 0}
	friendlyBorders.add(this.countryIndex[initialTerritory])

	p := &Path{
		Map:             this,
		continents:      continents,
		EnemyBorders:    &TerritorySet{data: 0},
		FriendlyBorders: &friendlyBorders,
		TotalScore:      0,
		Conquest:        initialTerritory,
		Territories:     this.index([]string{initialTerritory}),
		distance:        0,
		Parent:          nil,
	}
	p.detectBorders()

	return p
}

func (this *RiskBoard) buildTerritoryPath(continents *ContinentSet, InitialTerritories []string) *Path {

	p := this.startPath(continents, InitialTerritories[0])
	for _, country := range InitialTerritories[1:] {
		nextPath := p.conquer(this.countryIndex[country])
		p = nextPath
		p.TotalScore = 0
	}
	p.setTotalScore()

	return p
}

func (r *RiskBoard) index(countries []string) *TerritorySet {
	data := uint64(0)
	setBit := func(n uint64, pos uint64) uint64 {
		n |= (1 << pos)
		return n
	}
	for _, str := range countries {
		data = setBit(data, r.countryIndex[str])
	}
	return &TerritorySet{data}
}

func (r *RiskBoard) clearPaths() {
	r.paths = make(map[PathCoordinate]*Path)
}

func (r *RiskBoard) includePath(path *Path, territory uint64) {
	pathId := path.Territories.data
	coord := PathCoordinate{territory, pathId}

	// get the paths at this node
	matchingPath, exists := r.paths[coord]

	if exists {
		if matchingPath.TotalScore >= path.TotalScore {
			path.isRedundant = true
		} else {
			matchingPath.isRedundant = true
			r.paths[coord] = path
		}
	}
}
