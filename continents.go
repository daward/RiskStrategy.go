package main

type Continent struct {
	territories *TerritorySet
	value       int
}

type ContinentSet struct {
	continents []Continent
}

func (c *Continent) score(p *Path) int {
	if p.Territories.containsSet(c.territories) {
		return c.value
	}
	return 0
}

func (cs *ContinentSet) score(p *Path) int {
	retVal := 0
	for _, continent := range cs.continents {
		retVal += continent.score(p)
	}

	return retVal
}

func NewContinentSet(risk *RiskBoard) *ContinentSet {
	retVal := ContinentSet{
		continents: []Continent{},
	}

	asia := Continent{
		value:       7,
		territories: risk.index([]string{"Siam", "India", "Middle East", "Afghanistan", "China", "Mongolia", "Japan", "Ural", "Siberia", "Irkutsk", "Yakutsk", "Kamchatka"}),
	}

	australia := Continent{
		value:       2,
		territories: risk.index([]string{"New Guinea", "Indonesia", "Western Australia", "Eastern Australia"}),
	}

	southAmerica := Continent{
		value: 2,
		territories: risk.index([]string{
			"Venezuela",
			"Peru",
			"Brazil",
			"Argentina",
		}),
	}

	northAmerica := Continent{
		value: 5,
		territories: risk.index([]string{
			"Central America",
			"Eastern United States",
			"Western United States",
			"Quebec",
			"Ontario",
			"Alberta",
			"Northwest Territory",
			"Alaska",
			"Greenland",
		}),
	}

	europe := Continent{
		value: 5,
		territories: risk.index([]string{
			"Iceland",
			"Great Britain",
			"Scandinavia",
			"Northern Europe",
			"Western Europe",
			"Southern Europe",
			"Ukraine",
		}),
	}

	africa := Continent{
		value: 3,
		territories: risk.index([]string{
			"North Africa",
			"Egypt",
			"East Africa",
			"Congo",
			"South Africa",
			"Madagascar",
		}),
	}

	retVal.continents = []Continent{asia, australia, southAmerica, northAmerica, europe, africa}
	return &retVal
}
