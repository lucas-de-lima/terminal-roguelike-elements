package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type GameState int

const (
	StateMainMenu GameState = iota
	StateMap
	StateCombat
	StateLevelUp
	StateMilestone
	StateShop
	StateGameOver
)

type model struct {
	state            GameState
	menuCursor       int
	menuStep         int   // 0 para Classe, 1 para Elemento
	chosenClass      int   // Salva a classe escolhida na primeira etapa
	floor            int   // Controla a profundidade (Andar atual)
	currentBiome     Biome // Guarda o bioma do andar
	grid             [MapH][MapW]rune
	enemies          []*Character
	totalEnemies     int // Guarda o total inicial para a Bússola
	playerX          int
	playerY          int
	player           *Character
	enemy            *Character
	log              string
	levelOptions     []Skill
	milestoneOptions []StatReward
}

func initialModel() model {
	return model{
		state:      StateMainMenu,
		menuCursor: 0,
		menuStep:   0,
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
		case StateMainMenu:
			// Define o tamanho da lista baseada no passo atual
			maxCursor := len(AvailableClasses)
			if m.menuStep == 1 {
				maxCursor = len(AvailableElements)
			}

			switch msg.String() {
			case "up", "w":
				m.menuCursor--
				if m.menuCursor < 0 {
					m.menuCursor = maxCursor - 1
				}
				return m, nil
			case "down", "s":
				m.menuCursor++
				if m.menuCursor >= maxCursor {
					m.menuCursor = 0
				}
				return m, nil
			case "enter", " ":
				if m.menuStep == 0 {
					// Passo 1 Concluído: Salvou a classe, vai pro elemento
					m.chosenClass = m.menuCursor
					m.menuStep = 1
					m.menuCursor = 0 // Reseta o cursor pro novo menu
					return m, nil
				} else {
					// Passo 2 Concluído: Inicia o jogo!
					return startGame(m.chosenClass, m.menuCursor), nil
				}
			case "esc":
				if m.menuStep == 1 {
					// Permite voltar para a seleção de classe
					m.menuStep = 0
					m.menuCursor = m.chosenClass
					return m, nil
				}
			}
		case StateMap:
			return updateMap(m, msg)
		case StateCombat:
			return updateCombat(m, msg)
		case StateLevelUp:
			return updateLevelUp(m, msg)
		case StateMilestone:
			return updateMilestone(m, msg)
		case StateShop:
			return updateShop(m, msg)
		case StateGameOver:
			if msg.String() == "r" {
				return initialModel(), nil
			}
		}
	}
	return m, nil
}
