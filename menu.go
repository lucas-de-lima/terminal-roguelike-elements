package main

import (
	"fmt"
)

func startGame(classIndex, elementIndex int) model {
	m := model{state: StateMap, floor: 1}
	c := AvailableClasses[classIndex]
	e := AvailableElements[elementIndex]

	m.player = &Character{
		Name:    c.Name,
		Symbol:  c.Symbol,
		Element: e.Name,
		Stats: Stats{
			MaxHP:          c.HP,
			HP:             c.HP,
			Str:            c.Str,
			Int:            c.Int,
			CritChance:     c.CritChance,
			CritMultiplier: c.CritMultiplier,
			Level:          1,
			XP:             0,
			NextXP:         50,
		},
		Skills: []Skill{&GolpeRapido{BaseSkill{name: "Ataque Rápido", level: 1, element: "Qualquer"}}},
	}

	generateFloor(&m)
	m.log = fmt.Sprintf("A descida começa no Bioma: %s", m.currentBiome.Name)
	return m
}
