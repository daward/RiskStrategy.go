package main

type Path struct {
	BorderTerritories *TerritorySet
	Borders           *TerritorySet
	TotalScore        float64
	Map               *RiskBoard
	continents        *ContinentSet
	isRedundant       bool
	Conquests         []string
	IndexedConquests  *TerritorySet
}

func (p *Path) detectBorders() {
	// go through all the territories
	p.IndexedConquests.walk(func(territory uint64) {
		p.Map.board.Visit(int(territory), func(neighbor int, c int64) (skip bool) {
			// if this particular node is not a territory already,
			// then this territory is a border country and the visited node is a border
			if !p.isTerritory(uint64(neighbor)) {
				// we have a neighbor we haven't conquered yet (its not a territory), that means its a border land
				p.Borders.add(uint64(neighbor))
				// the territory in question has a non-territory neighbor, therefore
				// it is a border territory
				p.BorderTerritories.add(territory)
			}
			return false
		})
	})
}

func (p *Path) isTerritory(territory uint64) bool {
	return p.IndexedConquests.contains(territory)
}

func (p *Path) detectNewBorders() int {

	count := p.BorderTerritories.size()
	p.BorderTerritories.walk(func(territory uint64) {
		protected := true
		// visit all the neighbors of the territories
		p.Map.board.Visit(int(territory), func(neighbor int, c int64) (skip bool) {
			// if this particular node is not a territory already,
			// then this territory is a border country and the visited node is a border
			if !p.IndexedConquests.contains(uint64(neighbor)) {
				// we have a neighbor we haven't conquered yet (its not a territory), that means its a border land
				p.Borders.add(uint64(neighbor))
				protected = false
			}
			return false
		})
		// if this new border territory is protected, then it is no longer a border territory
		if protected {
			p.BorderTerritories.remove(territory)
		}
	})
	return count
}

func (p *Path) score() float64 {
	// the reinforcements you get is the number territories you have divided by 3
	// (rounded down) plus any bonuses for continents you control

	if p.isComplete() {
		return 0
	}

	reinforcements := p.IndexedConquests.size()/3 + p.continents.score(p)

	// this is the key heuristic proposed. The strength of your position is measured
	// by the reinforcements you get, divided by the territories you have to protect
	return float64(reinforcements) / float64(p.BorderTerritories.size())
}

// your total score will always go up, it is an accumulation of all your previous scores
func (p *Path) setTotalScore() {
	p.TotalScore = p.TotalScore + p.score()
}

// at some point there is nowhere else to go!
func (p *Path) isComplete() bool {
	return p.BorderTerritories.size() == 0
}

// make a new path for each border country
func (p *Path) expand() []*Path {
	retVal := make([]*Path, p.Borders.size())
	i := 0
	p.Borders.walk(func(territory uint64) {
		retVal[i] = p.conquer(territory)
		i++
	})

	return retVal
}

func (p *Path) latestConquest() string {
	return p.Conquests[p.IndexedConquests.size()-1]
}

// add a new territory to your domain
func (p *Path) conquer(territory uint64) *Path {

	// make sure its valid
	if p.isTerritory(territory) {
		panic("You can only conquer lands adjacent to your own, you may not conquer something you already have")
	}

	// then copy the border territories over
	newBorderTerritories := TerritorySet{data: p.BorderTerritories.data}
	newBorderTerritories.add(territory)

	newTerritories := TerritorySet{data: p.IndexedConquests.data}
	newTerritories.add(territory)

	newConquests := make([]string, p.IndexedConquests.size()+1)
	copy(newConquests, p.Conquests)
	newConquests[p.IndexedConquests.size()] = p.Map.countryLookup[territory]

	// build the next path with those territories
	nextPath := Path{
		BorderTerritories: &newBorderTerritories,
		Borders:           &TerritorySet{data: 0},
		TotalScore:        p.TotalScore,
		Map:               p.Map,
		continents:        p.continents,
		isRedundant:       false,
		Conquests:         newConquests,
		IndexedConquests:  &newTerritories,
	}

	// record it as a path to this territory
	p.Map.includePath(&nextPath, territory)

	if !nextPath.isRedundant {

		// build out the borders
		nextPath.detectNewBorders()

		// and add to our running total score
		nextPath.setTotalScore()
	}

	return &nextPath
}

func (this *Path) markRedundant(p *Path) bool {

	if p.IndexedConquests != this.IndexedConquests {
		return false
	}

	if this.TotalScore > p.TotalScore {
		p.isRedundant = true
	} else {
		this.isRedundant = true
	}

	return true
}
