package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"

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

	// 1. Lógica do Baú
	if tile == TileChest {
		xpGain := 40.0 + (float64(m.player.Stats.Level) * 15.0)
		m.player.Stats.XP += xpGain
		m.player.Stats.HP = math.Min(m.player.Stats.HP+30, m.player.Stats.MaxHP)
		m.log = fmt.Sprintf("🎁 Baú! Você recuperou HP e ganhou %.0f XP.", xpGain)
		m.grid[ny][nx] = TileFloor
		m.playerX, m.playerY = nx, ny
		return checkLevelUp(m), nil
	}

	// 2. Colisão Dinâmica de Combate (Procura na Lista)
	for i, e := range m.enemies {
		if e.X == nx && e.Y == ny {
			m.state = StateCombat
			m.enemy = e
			m.enemies = append(m.enemies[:i], m.enemies[i+1:]...) // Remove do mapa
			m.playerX, m.playerY = nx, ny

			if strings.Contains(e.Name, "Tóxico") {
				m.log = fmt.Sprintf("☣️ COMBATE TÓXICO! %s atacou!", m.enemy.Name)
			} else {
				m.log = fmt.Sprintf("⚔️ COMBATE! %s te ataca!", m.enemy.Name)
			}
			return m, nil
		}
	}

	// 3. Atualiza Jogador e Miasma
	m.playerX, m.playerY = nx, ny
	if tile == TileMiasma {
		m.player.Stats.HP -= 1
		m.log = "☣️ O Miasma queima sua pele! (-1 HP)"
		if m.player.Stats.HP <= 0 {
			m.state = StateGameOver
		}
	} else {
		m.log = "Avançando pelas masmorras..."
	}

	// 4. INTELIGÊNCIA ARTIFICIAL: Patrulha dos Monstros
	for _, e := range m.enemies {
		// 40% de chance de andar quando você anda
		if rand.Float64() < 0.4 {
			edx, edy := 0, 0
			switch rand.Intn(4) {
			case 0:
				edy = -1
			case 1:
				edy = 1
			case 2:
				edx = -1
			case 3:
				edx = 1
			}

			enx, eny := e.X+edx, e.Y+edy

			// Checa limites do mundo
			if enx >= 0 && enx < MapW && eny >= 0 && eny < MapH {
				targetTile := m.grid[eny][enx]

				// Monstro não entra na parede, em baú e não sobrepõe o jogador
				if targetTile != TileWall && targetTile != TileChest && !(enx == m.playerX && eny == m.playerY) {
					// Impede dois monstros no mesmo tile
					occupied := false
					for _, other := range m.enemies {
						if other.X == enx && other.Y == eny {
							occupied = true
							break
						}
					}
					if !occupied {
						e.X, e.Y = enx, eny // Anda!
					}
				}
			}
		}
	}

	return m, nil
}
