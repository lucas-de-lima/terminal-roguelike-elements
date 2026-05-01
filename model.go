package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type GameState int

const (
	StateMenuClass GameState = iota
	StateMenuElement
	StateMap
	StateCombat
	StateLevelUp
	StateGameOver
)

type model struct {
	state         GameState
	classCursor   int
	elementCursor int
	grid          [MapH][MapW]rune
	playerX       int
	playerY       int
	player        *Character
	enemy         *Character
	log           string
	levelOptions  []Skill
}

func initialModel() model {
	return model{
		state:         StateMenuClass,
		classCursor:   0,
		elementCursor: 0,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}

		switch m.state {
		case StateMenuClass:
			switch msg.String() {
			case "up", "w":
				m.classCursor--
				if m.classCursor < 0 {
					m.classCursor = len(AvailableClasses) - 1
				}
			case "down", "s":
				m.classCursor++
				if m.classCursor >= len(AvailableClasses) {
					m.classCursor = 0
				}
			case "enter", " ":
				// TODO: Go to StateMenuElement
				return startGame(m.classCursor, 0), nil // Placeholder for elementIndex
			}
		case StateMenuElement:
			// TODO: Implement element selection
		case StateMap:
			return updateMap(m, msg)
		case StateCombat:
			return updateCombat(m, msg)
		case StateLevelUp:
			return updateLevelUp(m, msg)
		case StateGameOver:
			if msg.String() == "r" {
				return initialModel(), nil
			}
		}
	}
	return m, nil
}
