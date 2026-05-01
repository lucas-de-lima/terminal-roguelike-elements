package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Stats struct {
	MaxHP, HP  float64
	Str, Int   float64
	Level      int
	XP, NextXP float64
}

type ClassDef struct {
	Name    string
	Desc    string
	Symbol  string
	HP      float64
	Str     float64
	Int     float64
	Element string
}

var AvailableClasses = []ClassDef{
	{"Guerreiro", "Combate corpo-a-corpo pesado. Alta sobrevivência.", EmojiWarrior, 150, 12, 4, "Earth"},
	{"Mago Elementalista", "Domínio sobre os elementos. Alto dano mágico.", EmojiMage, 80, 4, 15, "Fire"},
	{"Ladino", "Ataques furtivos e críticos. Status equilibrados.", EmojiRogue, 100, 10, 8, "Wind"},
}

type Character struct {
	Name    string
	Symbol  string
	Stats   Stats
	Skills  []Skill
	Element string
}

// Helper para gerenciar Level Up (Agora usado pelo Baú e pelo Combate)
func checkLevelUp(m model) model {
	if m.player.Stats.XP >= m.player.Stats.NextXP {
		m.state = StateLevelUp
		m.player.Stats.XP -= m.player.Stats.NextXP
		m.player.Stats.NextXP *= 1.3
		m.player.Stats.Level++
		m.player.Stats.MaxHP *= 1.1
		m.player.Stats.HP = m.player.Stats.MaxHP // Cura total no level up
		m.levelOptions = []Skill{getRandomNewSkill(), getRandomNewSkill()}
		m.log = "✨ LEVEL UP! Escolha um dom."
	} else {
		m.state = StateMap
	}
	return m
}

func updateLevelUp(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	var choice int
	fmt.Sscanf(msg.String(), "%d", &choice)
	if choice >= 1 && choice <= len(m.levelOptions) {
		sel := m.levelOptions[choice-1]
		found := false
		for _, s := range m.player.Skills {
			if s.Description() == sel.Description() {
				s.Upgrade()
				found = true
				break
			}
		}
		if !found {
			m.player.Skills = append(m.player.Skills, sel)
		}
		m.state = StateMap
	}
	return m, nil
}
