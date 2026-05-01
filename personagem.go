package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// --- ELEMENTOS ---
const (
	ElemFire  = "Fogo"
	ElemWind  = "Vento"
	ElemLight = "Raio"
	ElemEarth = "Terra"
	ElemWater = "Água"
	ElemNone  = "Neutro"
)

// Estrutura para exibir no menu
type ElementDef struct {
	Name   string
	Symbol string
}

var AvailableElements = []ElementDef{
	{ElemFire, "🔥"},
	{ElemWater, "💧"},
	{ElemWind, "🍃"},
	{ElemEarth, "🪨"},
	{ElemLight, "⚡"},
}

// Lógica de Vantagem: Retorna o multiplicador de dano
func getElementalMultiplier(atkElem, defElem string) float64 {
	if atkElem == defElem || atkElem == ElemNone || defElem == ElemNone {
		return 1.0 // Dano normal
	}

	// Quem ganha de quem (+30% de dano)
	strongAgainst := map[string]string{
		ElemFire:  ElemWind,
		ElemWind:  ElemLight,
		ElemLight: ElemEarth,
		ElemEarth: ElemWater,
		ElemWater: ElemFire,
	}

	// Quem perde pra quem (-30% de dano)
	weakAgainst := map[string]string{
		ElemWind:  ElemFire,
		ElemLight: ElemWind,
		ElemEarth: ElemLight,
		ElemWater: ElemEarth,
		ElemFire:  ElemWater,
	}

	if strongAgainst[atkElem] == defElem {
		return 1.3
	} else if weakAgainst[atkElem] == defElem {
		return 0.7
	}

	return 1.0 // Se não interagir, dano normal
}

type Stats struct {
	MaxHP, HP      float64
	Str, Int       float64
	CritChance     float64 // Ex: 0.10 para 10%
	CritMultiplier float64 // Ex: 1.5 para +50% de dano
	Level          int
	XP, NextXP     float64
}

type ClassDef struct {
	Name           string
	Desc           string
	Symbol         string
	HP             float64
	Str            float64
	Int            float64
	CritChance     float64
	CritMultiplier float64
	Element        string
}

var AvailableClasses = []ClassDef{
	{"Guerreiro", "Dano físico massivo e defesa.", EmojiWarrior, 150, 12, 4, 0.05, 1.5, ElemEarth},
	{"Mago Elementalista", "Domínio mágico total.", EmojiMage, 80, 4, 15, 0.05, 1.5, ElemFire},
	{"Ladino", "Ataques furtivos e letais.", EmojiRogue, 100, 10, 8, 0.15, 2.0, ElemWind},
}

type Character struct {
	Name    string
	Symbol  string
	Stats   Stats
	Skills  []Skill
	Element string
	X, Y    int
}

// --- SISTEMA DE MARCOS (A cada 5 níveis) ---
type StatReward struct {
	Name  string
	Desc  string
	Apply func(c *Character)
}

var MilestoneRewards = []StatReward{
	{"Treino de Força", "+15% Força Base", func(c *Character) { c.Stats.Str *= 1.15 }},
	{"Meditação Profunda", "+15% Inteligência Base", func(c *Character) { c.Stats.Int *= 1.15 }},
	{"Golpe Vital", "+5% Chance Crítica", func(c *Character) { c.Stats.CritChance += 0.05 }},
	{"Brutalidade", "+30% Multiplicador Crítico", func(c *Character) { c.Stats.CritMultiplier += 0.30 }},
	{"Vigor", "+20% Vida Máxima", func(c *Character) { c.Stats.MaxHP *= 1.20 }},
}

// Helper para gerenciar Level Up (Agora usado pelo Baú e pelo Combate)
func checkLevelUp(m model) model {
	if m.player.Stats.XP >= m.player.Stats.NextXP {
		m.player.Stats.XP -= m.player.Stats.NextXP
		m.player.Stats.NextXP *= 1.5
		m.player.Stats.Level++
		m.player.Stats.MaxHP *= 1.05
		m.player.Stats.HP = m.player.Stats.MaxHP // Cura total no level up

		if m.player.Stats.Level%5 == 0 {
			// A cada 5 níveis, escolhe status!
			m.state = StateMilestone
			m.milestoneOptions = MilestoneRewards
			m.log = "🌟 MARCO ATINGIDO! Melhore seus atributos."
		} else {
			// Nível normal, escolhe skills filtradas
			m.state = StateLevelUp
			m.levelOptions = getSkillOptions(m.player)
			m.log = "✨ LEVEL UP! Escolha um dom."
		}
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

func updateMilestone(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	var choice int
	fmt.Sscanf(msg.String(), "%d", &choice)
	if choice >= 1 && choice <= len(m.milestoneOptions) {
		reward := m.milestoneOptions[choice-1]
		reward.Apply(m.player)
		m.log = fmt.Sprintf("🌟 %s aplicado!", reward.Name)
		m.state = StateMap
	}
	return m, nil
}
