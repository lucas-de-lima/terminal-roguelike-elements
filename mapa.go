package main

import (
	"fmt"
	"math"

	tea "github.com/charmbracelet/bubbletea"
)

// --- 1. CONFIGURAÇÕES DO MUNDO ---

const (
	MapW, MapH   = 80, 40
	ViewW, ViewH = 30, 12
	VisionRadius = 6.0

	// Tipos de Tiles na Matriz
	TileWall   = '#'
	TileFloor  = '.'
	TileMiasma = '~'
	TileChest  = 'C'
)

func updateMap(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	dx, dy := 0, 0
	switch msg.String() {
	case "w", "up":
		dy = -1
	case "s", "down":
		dy = 1
	case "a", "left":
		dx = -1
	case "d", "right":
		dx = 1
	}

	if dx == 0 && dy == 0 {
		return m, nil
	}

	nx, ny := m.playerX+dx, m.playerY+dy
	if nx < 0 || nx >= MapW || ny < 0 || ny >= MapH {
		return m, nil
	}

	tile := m.grid[ny][nx]
	if tile == TileWall {
		return m, nil
	}

	// Lógica do Baú
	if tile == TileChest {
		xpGain := 40.0 + (float64(m.player.Stats.Level) * 15.0)
		m.player.Stats.XP += xpGain
		m.player.Stats.HP = math.Min(m.player.Stats.HP+30, m.player.Stats.MaxHP) // Cura 30 HP

		m.log = fmt.Sprintf("🎁 Baú! Você recuperou HP e ganhou %.0f XP.", xpGain)
		m.grid[ny][nx] = TileFloor // O baú some após ser aberto
		m.playerX, m.playerY = nx, ny

		return checkLevelUp(m), nil
	}

	// Lógica de Combate ('E' normal, 'M' mutante)
	if tile == 'E' || tile == 'M' {
		m.state = StateCombat
		isMutant := (tile == 'M')

		// 1. Procura qual inimigo está nesta exata coordenada
		var foundEnemy *Character
		var foundIdx int
		for i, e := range m.enemies {
			if e.X == nx && e.Y == ny {
				foundEnemy = e
				foundIdx = i
				break
			}
		}

		m.enemy = foundEnemy

		// 2. Remove o inimigo da lista global para ele não existir mais no mapa
		m.enemies = append(m.enemies[:foundIdx], m.enemies[foundIdx+1:]...)

		// 3. Volta o tile para o chão original
		if isMutant {
			m.grid[ny][nx] = TileMiasma
		} else {
			m.grid[ny][nx] = TileFloor
		}
		m.playerX, m.playerY = nx, ny

		if isMutant {
			m.log = fmt.Sprintf("☣️ COMBATE TÓXICO! %s atacou!", m.enemy.Name)
		} else {
			m.log = fmt.Sprintf("⚔️ COMBATE! %s te ataca!", m.enemy.Name)
		}
		return m, nil
	}

	// Lógica de andar e Miasma
	m.playerX, m.playerY = nx, ny

	if tile == TileMiasma {
		m.player.Stats.HP -= 1 // Dano contínuo
		m.log = "☣️ O Miasma queima sua pele! (-1 HP)"
		if m.player.Stats.HP <= 0 {
			m.state = StateGameOver
		}
	} else {
		m.log = "Avançando pelas masmorras..."
	}

	return m, nil
}
