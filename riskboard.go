package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type RiskBoard struct {
	board         *Graph
	countryIndex  map[string]uint64
	countryLookup []string
	paths         map[PathCoordinate]*Path // territory plus length of path
}

type PathCoordinate struct {
	territory uint64
	pathId    uint64
}

type Node struct {
	links uint64
}

type Graph struct {
	nodes []Node
}

func (this *Graph) AddBoth(v1, v2 uint64) {
	this.nodes[v1].links = this.nodes[v1].links | (1 << v2)
	this.nodes[v2].links = this.nodes[v2].links | (1 << v1)
}

func (this *Graph) Visit(node uint64, callback func(neighbor uint64) bool) {
	currentNode := this.nodes[node]
	for link := uint64(0); link < 42; link++ {
		// Check if the current bit is set.
		if currentNode.links&(1<<link) != 0 {
			if callback(link) {
				break
			}
		}
	}
}

func initGraph() *Graph {
	g := Graph{
		nodes: make([]Node, 42),
	}
	for _, node := range g.nodes {
		node.links = 0
	}

	return &g
}

func riskboard() *RiskBoard {

	g := initGraph()
	jsonFile, _ := os.Open("./mapping.json")
	// if we os.Open returns an error then handle it
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var arr [][]string
	json.Unmarshal([]byte(byteValue), &arr)

	m := make(map[string]uint64)
	m2 := make([]string, 42)

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

		g.AddBoth(m[v1[0]], m[v1[1]])
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
		Territory:       this.countryIndex[initialTerritory],
		Territories:     this.index([]string{initialTerritory}),
		distance:        0,
		Parent:          nil,
	}
	p.detectBorders()

	return p
}

func (this *RiskBoard) Visit(node uint64, callback func(neighbor uint64) bool) {
	this.board.Visit(node, callback)
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
