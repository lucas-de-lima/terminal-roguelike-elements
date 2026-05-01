package main

import (
	"fmt"
	"math"
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
	elemento := baseEnemy.Element

	if isMutant {
		lvl += 3 // Mutantes ganham 3 níveis de cara!
		nome = nome + " (Tóxico)"
	}

	// Curva de Dificuldade Exponencial!
	// A cada nível, os monstros ficam 15% mais parrudos.
	scale := float64(lvl - 1)
	hp := 25.0 * math.Pow(1.15, scale)
	str := 3.0 * math.Pow(1.15, scale)

	// O BUFF DO MIASMA: Inimigos ali dentro são mini-chefes
	if isMutant {
		hp *= 1.8  // 80% mais vida
		str *= 1.4 // 40% mais dano
	}

	return &Character{
		Name:    nome,
		Symbol:  baseEnemy.Symbol,
		Element: elemento,
		Stats: Stats{
			MaxHP:          hp,
			HP:             hp,
			Str:            str,
			Int:            str,
			CritChance:     0.05,
			CritMultiplier: 1.5,
		},
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

		// GATILHO DO PORTAL
		if len(m.enemies) == 0 {
			m.grid[m.playerY][m.playerX] = TilePortal
			m.log += " 🌀 A área está limpa! Um PORTAL se abriu sob você."
		}

		return checkLevelUp(m), nil
	}

	// Turno do Inimigo (usando sistema de crítico também)
	baseDmg := m.enemy.Stats.Str * (0.8 + rand.Float64()*0.4)
	elemMult := getElementalMultiplier(m.enemy.Element, m.player.Element)

	isCrit := rand.Float64() < m.enemy.Stats.CritChance
	critMult := 1.0
	critText := ""
	if isCrit {
		critMult = m.enemy.Stats.CritMultiplier
		critText = " ⚡CRÍTICO!"
	}

	finalDmg := baseDmg * elemMult * critMult
	m.player.Stats.HP -= finalDmg

	feedbackInimigo := ""
	if elemMult > 1.0 {
		feedbackInimigo = " (Super Efetivo!)"
	}
	if elemMult < 1.0 {
		feedbackInimigo = " (Resistido)"
	}

	m.log += fmt.Sprintf("\n💀 %s atacou: -%.0f HP%s%s", m.enemy.Symbol, finalDmg, critText, feedbackInimigo)

	if m.player.Stats.HP <= 0 {
		m.state = StateGameOver
	}
	return m, nil
}
