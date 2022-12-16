package main

type Path struct {
	// friendly borders are all the territories you own
	// that are adjacent to territories you do not own
	FriendlyBorders *TerritorySet

	// enemy borders are the territories you don't own that are next to
	// territories that you do own
	EnemyBorders *TerritorySet

	// territories are all the countries you occupy
	Territories *TerritorySet

	// the accumulated score across all conquests
	TotalScore float64

	// the risk board this path is on
	Map *RiskBoard

	// the continents on the board
	continents *ContinentSet

	// if at some point we find out this path is not any better than an identical path
	// we will call it redundant and stop
	isRedundant bool

	// the total number of conquests that have taken place
	distance int

	// the path that this path is built on
	Parent *Path

	Territory uint64
}

func (p *Path) detectBorders() {
	// go through all the territories
	p.Territories.walk(func(territory uint64) {
		p.Map.board.Visit(territory, func(neighbor uint64) (skip bool) {
			// if this particular node is not a territory already,
			// then this territory is a border country and the visited node is a border
			if !p.isTerritory(uint64(neighbor)) {
				// we have a neighbor we haven't conquered yet (its not a territory), that means its a border land
				p.EnemyBorders.add(uint64(neighbor))
				// the territory in question has a non-territory neighbor, therefore
				// it is a border territory
				p.FriendlyBorders.add(territory)
			}
			return false
		})
	})
}

func (p *Path) isTerritory(territory uint64) bool {
	return p.Territories.contains(territory)
}

func (p *Path) detectNewBorders() {

	// assume our new conquest ends up being protected by its neighbors
	protectedConquest := true
	// only the new node has the ability to change anything, so focus there
	p.Map.Visit(p.Territory, func(neighbor uint64) (skip bool) {

		neighborId := uint64(neighbor)
		// for all our neighbors that are friendly borders, we have a little more work to do
		if p.FriendlyBorders.contains(neighborId) {
			//fmt.Println(p.Conquest + " -> " + p.Map.countryLookup[neighborId] + " (friendly)")
			// assume to start that the neighbor of our territory is now protected by our new conquest
			protectedBorder := true
			// start visiting our neighbor's neighbors to find out if it actually is protected
			p.Map.Visit(neighbor, func(n2 uint64) (skip bool) {

				// if the neighbor has a neighbor that is not one of our territories, then
				// our new conquest did not protect it and we can bail out
				if !p.isTerritory(uint64(n2)) {
					//fmt.Println(p.Conquest + " -> " + p.Map.countryLookup[neighborId] + " -> " + p.Map.countryLookup[uint64(n2)] + " (enemy)")
					protectedBorder = false
					// bail out of the visitation loop
					return true
				}
				//fmt.Println(p.Conquest + " -> " + p.Map.countryLookup[neighborId] + " -> " + p.Map.countryLookup[uint64(n2)] + " (friendly)")
				return false
			})

			// if our neighbor is protected, we can take it out of the friendly borders
			if protectedBorder {
				p.FriendlyBorders.remove(neighborId)
			}

		} else { // but if our neighbor is not a friendly border, then we know the new conquest is not protected and we have a new enemy!
			//fmt.Println(p.Conquest + " -> " + p.Map.countryLookup[neighborId] + " (enemy)")
			p.EnemyBorders.add(neighborId)
			protectedConquest = false
		}
		return false
	})

	// this isn't a border, its completely surrounded by friendly territories
	if protectedConquest {
		// if it is protected, take it out
		p.FriendlyBorders.remove(p.Territory)
	}
}

func (p *Path) score() float64 {
	// the reinforcements you get is the number territories you have divided by 3
	// (rounded down) plus any bonuses for continents you control

	if p.isComplete() {
		return 0
	}

	reinforcements := p.Territories.size()/3 + p.continents.score(p)

	// this is the key heuristic proposed. The strength of your position is measured
	// by the reinforcements you get, divided by the territories you have to protect
	return float64(reinforcements) / float64(p.FriendlyBorders.size())
}

// your total score will always go up, it is an accumulation of all your previous scores
func (p *Path) setTotalScore() {
	p.TotalScore = p.TotalScore + p.score()
}

// at some point there is nowhere else to go!
func (p *Path) isComplete() bool {
	return p.FriendlyBorders.size() == 0
}

// make a new path for each border country
func (p *Path) expand() []*Path {
	retVal := make([]*Path, p.EnemyBorders.size())
	i := 0
	p.EnemyBorders.walk(func(territory uint64) {
		retVal[i] = p.conquer(territory)
		i++
	})

	return retVal
}

// recursively walk up the parents to put together the exact
// path you took to get here
func (p *Path) conquests() []string {
	current := p
	retVal := make([]string, current.distance+1)
	for {
		retVal[current.distance] = p.Map.countryLookup[current.Territory]
		if current.Parent == nil {
			break
		}
		current = current.Parent
		if current == current.Parent {
			panic("Don't be your own grandpa")
		}
	}

	return retVal
}

// add a new territory to your domain
func (p *Path) conquer(territory uint64) *Path {

	// make sure its valid
	if p.isTerritory(territory) {
		panic("You can only conquer lands adjacent to your own, you may not conquer something you already have")
	}

	// build out a new copy of each territory set
	newFriendlyBorders := TerritorySet{data: p.FriendlyBorders.data}
	newFriendlyBorders.add(territory)

	newEnemyBorders := TerritorySet{data: p.EnemyBorders.data}
	newEnemyBorders.remove(territory)

	newTerritories := TerritorySet{data: p.Territories.data}
	newTerritories.add(territory)

	// build the next path with those territories
	nextPath := Path{
		FriendlyBorders: &newFriendlyBorders,
		EnemyBorders:    &newEnemyBorders,
		TotalScore:      p.TotalScore,
		Map:             p.Map,
		continents:      p.continents,
		isRedundant:     false,
		Territory:       territory,
		Territories:     &newTerritories,
		Parent:          p,
		distance:        p.distance + 1,
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

	if p.Territories != this.Territories {
		return false
	}

	if this.TotalScore > p.TotalScore {
		p.isRedundant = true
	} else {
		this.isRedundant = true
	}

	return true
}
