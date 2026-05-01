package main

import (
	"fmt"
	"math/rand"

	tea "github.com/charmbracelet/bubbletea"
)

// Definição de moldes de inimigos
type EnemyDef struct {
	Name    string
	Symbol  string
	Element string
}

var EnemyPool = []EnemyDef{
	{"Slime Aquático", "🐸", ElemWater},
	{"Goblin das Chamas", "👺", ElemFire},
	{"Golem de Granito", "🪨", ElemEarth},
	{"Harpia Cortante", "🦅", ElemWind},
	{"Espírito Faísca", "🌩️", ElemLight},
}

func generateEnemy(playerLvl int, isMutant bool) *Character {
	lvl := playerLvl
	baseEnemy := EnemyPool[rand.Intn(len(EnemyPool))]

	nome := baseEnemy.Name
	elemento := baseEnemy.Element // <--- MANTÉM O ELEMENTO NATURAL!

	// Mutantes ganham nível, mas continuam vulneráveis às fraquezas normais
	if isMutant {
		lvl += 2
		nome = nome + " (Tóxico)"
	}

	scale := float64(lvl)
	hp := 20.0 * (1 + scale*0.3)

	return &Character{
		Name:    nome,
		Symbol:  baseEnemy.Symbol,
		Element: elemento,
		Stats:   Stats{MaxHP: hp, HP: hp, Str: 2.0 * (1 + scale*0.2)},
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

	// Turno do Inimigo
	enemyElemMult := getElementalMultiplier(m.enemy.Element, m.player.Element)
	dmg := m.enemy.Stats.Str * (0.8 + rand.Float64()*0.4) * enemyElemMult

	m.player.Stats.HP -= dmg

	feedbackInimigo := ""
	if enemyElemMult > 1.0 {
		feedbackInimigo = " (Dano Crítico!)"
	}

	m.log += fmt.Sprintf("\n💀 %s atacou: -%.0f HP%s", m.enemy.Symbol, dmg, feedbackInimigo)

	if m.player.Stats.HP <= 0 {
		m.state = StateGameOver
	}
	return m, nil
}
