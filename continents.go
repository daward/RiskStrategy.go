package main

type Continent struct {
	territories []uint64
	value       int
}

type ContinentSet struct {
	continents []Continent
}

func (c *Continent) score(p *Path) int {
	for _, territory := range c.territories {
		if !p.isTerritory(uint64(territory)) {
			return 0
		}
	}
	return c.value
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
		value: 7,
		territories: []uint64{
			risk.countryIndex["Siam"],
			risk.countryIndex["India"],
			risk.countryIndex["Middle East"],
			risk.countryIndex["Afghanistan"],
			risk.countryIndex["China"],
			risk.countryIndex["Mongolia"],
			risk.countryIndex["Japan"],
			risk.countryIndex["Ural"],
			risk.countryIndex["Siberia"],
			risk.countryIndex["Irkutsk"],
			risk.countryIndex["Yakutsk"],
			risk.countryIndex["Kamchatka"],
		},
	}

	australia := Continent{
		value: 2,
		territories: []uint64{
			risk.countryIndex["New Guinea"],
			risk.countryIndex["Indonesia"],
			risk.countryIndex["Western Australia"],
			risk.countryIndex["Eastern Australia"],
		},
	}

	southAmerica := Continent{
		value: 2,
		territories: []uint64{
			risk.countryIndex["Venezuela"],
			risk.countryIndex["Peru"],
			risk.countryIndex["Brazil"],
			risk.countryIndex["Argentina"],
		},
	}

	northAmerica := Continent{
		value: 5,
		territories: []uint64{
			risk.countryIndex["Central America"],
			risk.countryIndex["Eastern United States"],
			risk.countryIndex["Western United States"],
			risk.countryIndex["Quebec"],
			risk.countryIndex["Ontario"],
			risk.countryIndex["Alberta"],
			risk.countryIndex["Northwest Territory"],
			risk.countryIndex["Alaska"],
			risk.countryIndex["Greenland"],
		},
	}

	europe := Continent{
		value: 5,
		territories: []uint64{
			risk.countryIndex["Iceland"],
			risk.countryIndex["Great Britain"],
			risk.countryIndex["Scandinavia"],
			risk.countryIndex["Northern Europe"],
			risk.countryIndex["Western Europe"],
			risk.countryIndex["Southern Europe"],
			risk.countryIndex["Ukraine"],
		},
	}

	africa := Continent{
		value: 3,
		territories: []uint64{
			risk.countryIndex["North Africa"],
			risk.countryIndex["Egypt"],
			risk.countryIndex["East Africa"],
			risk.countryIndex["Congo"],
			risk.countryIndex["South Africa"],
			risk.countryIndex["Madagascar"],
		},
	}

	retVal.continents = append(retVal.continents, asia, australia, southAmerica, northAmerica, europe, africa)
	return &retVal
}
