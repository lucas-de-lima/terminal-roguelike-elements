package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"

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

var MiniBossPool = []EnemyDef{
	{"Minotauro Furioso", "🦍", ElemEarth},
	{"Aranha Viúva", "🕷️", ElemNone},
	{"Elemental Maior", "🦂", ElemFire},
}

var BossPool = []EnemyDef{
	{"Dragão Ancião", "🐉", ElemFire},
	{"Kraken Abissal", "🦑", ElemWater},
	{"T-Rex Zumbi", "🦖", ElemNone},
}

const (
	TypeNormal = iota
	TypeMiniBoss
	TypeBoss
)

func generateEnemy(playerLvl int, isMutant bool, enemyType int) *Character {
	lvl := playerLvl

	var baseEnemy EnemyDef
	switch enemyType {
	case TypeBoss:
		baseEnemy = BossPool[rand.Intn(len(BossPool))]
		lvl += 5
	case TypeMiniBoss:
		baseEnemy = MiniBossPool[rand.Intn(len(MiniBossPool))]
		lvl += 2
	default:
		baseEnemy = EnemyPool[rand.Intn(len(EnemyPool))]
	}

	nome := baseEnemy.Name
	elemento := baseEnemy.Element

	if isMutant {
		lvl += 3
		nome = nome + " (Tóxico)"
	}

	scale := float64(lvl - 1)
	hp := 25.0 * math.Pow(1.15, scale)
	str := 3.0 * math.Pow(1.15, scale)

	if enemyType == TypeBoss {
		hp *= 3.5
		str *= 1.8
	} else if enemyType == TypeMiniBoss {
		hp *= 2.0
		str *= 1.3
	}

	if isMutant {
		hp *= 1.5
		str *= 1.3
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

		goldDrop := 10 + rand.Intn(10)
		if strings.Contains(m.enemy.Name, "Dragão") || strings.Contains(m.enemy.Name, "Kraken") || strings.Contains(m.enemy.Name, "T-Rex") {
			goldDrop = 150 + rand.Intn(50)
		}
		m.player.Stats.Gold += goldDrop

		m.log = fmt.Sprintf("🏆 Venceu! +%.0f XP | +%d Ouro", xp, goldDrop)

		// Verifica se há Boss vivo
		hasBoss := false
		for _, e := range m.enemies {
			if strings.Contains(e.Name, "Dragão") || strings.Contains(e.Name, "Kraken") || strings.Contains(e.Name, "T-Rex") {
				hasBoss = true
				break
			}
		}

		if !hasBoss && len(m.enemies) < m.totalEnemies/2 {
			m.grid[m.playerY][m.playerX] = TilePortal
			m.log += " 🌀 O Mestre deste andar caiu. Um PORTAL apareceu."
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
