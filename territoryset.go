package main

import (
	"fmt"
	"math/bits"
)

type TerritorySet struct {
	data uint64
}

func (ts *TerritorySet) contains(territory uint64) bool {
	return ts.data&(1<<uint64(territory)) != 0
}

func (ts *TerritorySet) containsSet(territories *TerritorySet) bool {
	return ts.data&territories.data == territories.data
}

func (ts *TerritorySet) add(territory uint64) {
	ts.data = ts.data | (1 << territory)
}

func (ts *TerritorySet) remove(territory uint64) {
	ts.data = ts.data & ^(1 << territory)
}

func (ts *TerritorySet) walk(callback func(territory uint64)) {
	for territory := 0; territory < 42; territory++ {
		// Check if the current bit is set.
		if ts.data&(1<<territory) != 0 {
			callback(uint64(territory))
		}
	}
}

func (ts *TerritorySet) size() int {
	return bits.OnesCount64(ts.data)
}

func (ts *TerritorySet) print(risk *RiskBoard) {
	str := ""
	ts.walk(func(territory uint64) {
		str = str + ", " + risk.countryLookup[territory]
	})
	fmt.Println(str)
}
