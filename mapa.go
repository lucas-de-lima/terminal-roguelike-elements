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

	TileWall   = '#'
	TileFloor  = '.'
	TileMiasma = '~'
	TileChest  = 'C'
	TilePortal = 'O'
)

type Biome struct {
	Name        string
	EmojiWall   string
	EmojiHazard string
	HazardDesc  string
	ColorFloor  string
	ColorHazard string
}

var Biomes = []Biome{
	{"Masmorra Sombria", "🧱", "🟣", "O Miasma corrói sua pele!", "#333333", "#A020F0"},
	{"Floresta Anciã", "🌲", "🌺", "O pólen tóxico te sufoca!", "#228B22", "#FF1493"},
	{"Caverna Vulcânica", "⛰️", "🔥", "A lava derrete suas botas!", "#8B0000", "#FF4500"},
	{"Pico Congelado", "🧊", "🌨️", "O frio extremo drena sua vida!", "#B0E0E6", "#00FFFF"},
}

func generateFloor(m *model) {
	// 1. Sorteia o Bioma do Andar
	m.currentBiome = Biomes[rand.Intn(len(Biomes))]

	// 2. Limpa e Gera a estrutura do mapa
	var newGrid [MapH][MapW]rune
	for y := 0; y < MapH; y++ {
		for x := 0; x < MapW; x++ {
			if x == 0 || x == MapW-1 || y == 0 || y == MapH-1 || rand.Float64() < 0.15 {
				newGrid[y][x] = TileWall
			} else {
				newGrid[y][x] = TileFloor
			}
		}
	}

	// 3. Gera Zonas de Perigo (Hazard/Miasma)
	for i := 0; i < 8; i++ {
		mx, my := rand.Intn(MapW-10)+5, rand.Intn(MapH-10)+5
		for dy := -3; dy <= 3; dy++ {
			for dx := -3; dx <= 3; dx++ {
				if rand.Float64() < 0.6 {
					newGrid[my+dy][mx+dx] = TileMiasma
				}
			}
		}
	}

	// 4. Área Segura do Jogador
	sx, sy := MapW/2, MapH/2
	for y := sy - 2; y <= sy+2; y++ {
		for x := sx - 2; x <= sx+2; x++ {
			newGrid[y][x] = TileFloor
		}
	}
	m.playerX, m.playerY = sx, sy

	// 5. Baús
	for i := 0; i < 10; i++ {
		cx, cy := rand.Intn(MapW-2)+1, rand.Intn(MapH-2)+1
		if newGrid[cy][cx] != TileWall {
			newGrid[cy][cx] = TileChest
		}
	}

	// 6. Inimigos Escalonados pelo Andar!
	m.enemies = []*Character{}
	qtdInimigos := 25 + (m.floor * 5)

	for i := 0; i < qtdInimigos; i++ {
		ex, ey := rand.Intn(MapW-2)+1, rand.Intn(MapH-2)+1
		tileAtual := newGrid[ey][ex]
		if tileAtual == TileFloor || tileAtual == TileMiasma {
			isMutant := (tileAtual == TileMiasma)
			enemyLvl := m.player.Stats.Level + m.floor - 1
			enemy := generateEnemy(enemyLvl, isMutant)
			enemy.X = ex
			enemy.Y = ey
			m.enemies = append(m.enemies, enemy)
		}
	}

	m.grid = newGrid
	m.totalEnemies = len(m.enemies)
}

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

	// 1.5 Lógica do Portal
	if tile == TilePortal {
		m.floor++
		generateFloor(&m)
		m.log = fmt.Sprintf("🌀 Você desceu para o Andar %d: %s!", m.floor, m.currentBiome.Name)
		return m, nil
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
		m.log = "☣️ " + m.currentBiome.HazardDesc + " (-1 HP)"
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
