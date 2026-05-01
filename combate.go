package main

import (
	"fmt"
	"math/rand"

	tea "github.com/charmbracelet/bubbletea"
)

func generateEnemy(playerLvl int, isMutant bool) *Character {
	lvl := playerLvl
	nome := fmt.Sprintf("Monstro Nv.%d", lvl)
	element := "Neutral"

	// Monstro do miasma ganha 2 níveis grátis
	if isMutant {
		lvl += 2
		nome = fmt.Sprintf("Aberração Tóxica Nv.%d", lvl)
		element = "Poison"
	}

	scale := float64(lvl)
	hp := 20.0 * (1 + scale*0.3)
	return &Character{
		Name:    nome,
		Symbol:  EmojiEnemy,
		Stats:   Stats{MaxHP: hp, HP: hp, Str: 2.0 * (1 + scale*0.2)},
		Element: element,
	}
}

func updateCombat(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	var skillIndex int
	fmt.Sscanf(msg.String(), "%d", &skillIndex)
	if skillIndex < 1 || skillIndex > len(m.player.Skills) {
		return m, nil
	}

	skill := m.player.Skills[skillIndex-1]
	m.log = skill.Cast(m.player, m.enemy)

	if m.enemy.Stats.HP <= 0 {
		xp := 25.0 * m.enemy.Stats.Str
		m.player.Stats.XP += xp
		m.log = fmt.Sprintf("🏆 Venceu! +%.0f XP", xp)
		return checkLevelUp(m), nil
	}

	dmg := m.enemy.Stats.Str * (0.8 + rand.Float64()*0.4)
	m.player.Stats.HP -= dmg
	m.log += fmt.Sprintf("\n💀 %s atacou: -%.0f HP", m.enemy.Symbol, dmg)

	if m.player.Stats.HP <= 0 {
		m.state = StateGameOver
	}
	return m, nil
}
