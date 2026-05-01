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
	{"Masmorra Sombria", "🧱 ", "🟣 ", "O Miasma corrói sua pele!", "#333333", "#A020F0"},
	{"Floresta Anciã", "🌲 ", "🌺 ", "O pólen tóxico te sufoca!", "#228B22", "#FF1493"},
	{"Caverna Vulcânica", "⛰️ ", "🔥 ", "A lava derrete suas botas!", "#8B0000", "#FF4500"},
	{"Pico Congelado", "🧊 ", "🌨️ ", "O frio extremo drena sua vida!", "#B0E0E6", "#00FFFF"},
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

	// ==========================================
	// NOVO: Algoritmo Flood Fill para spawn seguro
	// ==========================================
	var reachable [MapH][MapW]bool
	queue := []struct{ x, y int }{{sx, sy}}
	reachable[sy][sx] = true

	dirs := []struct{ dx, dy int }{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for _, d := range dirs {
			nx, ny := curr.x+d.dx, curr.y+d.dy
			if nx >= 0 && nx < MapW && ny >= 0 && ny < MapH {
				if newGrid[ny][nx] != TileWall && !reachable[ny][nx] {
					reachable[ny][nx] = true
					queue = append(queue, struct{ x, y int }{nx, ny})
				}
			}
		}
	}

	type Coord struct{ x, y int }
	validSpawns := []Coord{}
	for y := 0; y < MapH; y++ {
		for x := 0; x < MapW; x++ {
			if reachable[y][x] && newGrid[y][x] != TileWall && !(x == sx && y == sy) {
				validSpawns = append(validSpawns, Coord{x, y})
			}
		}
	}

	rand.Shuffle(len(validSpawns), func(i, j int) {
		validSpawns[i], validSpawns[j] = validSpawns[j], validSpawns[i]
	})

	// 5. Baús usando coordenadas alcançáveis
	spawnIndex := 0
	for i := 0; i < 10 && spawnIndex < len(validSpawns); i++ {
		c := validSpawns[spawnIndex]
		newGrid[c.y][c.x] = TileChest
		spawnIndex++
	}

	// 6. Inimigos Escalonados pelo Andar!
	m.enemies = []*Character{}

	// 1. Spawna 1 Boss (Garantido)
	if spawnIndex < len(validSpawns) {
		c := validSpawns[spawnIndex]
		boss := generateEnemy(m.player.Stats.Level+m.floor, false, TypeBoss)
		boss.X, boss.Y = c.x, c.y
		m.enemies = append(m.enemies, boss)
		spawnIndex++
	}

	// 2. Spawna Mini-Bosses (1 a cada 2 andares)
	for i := 0; i < m.floor/2 && spawnIndex < len(validSpawns); i++ {
		c := validSpawns[spawnIndex]
		mini := generateEnemy(m.player.Stats.Level+m.floor, false, TypeMiniBoss)
		mini.X, mini.Y = c.x, c.y
		m.enemies = append(m.enemies, mini)
		spawnIndex++
	}

	// 3. Spawna o resto dos inimigos comuns
	qtdInimigos := 20 + (m.floor * 4)
	for i := 0; i < qtdInimigos && spawnIndex < len(validSpawns); i++ {
		c := validSpawns[spawnIndex]
		isMutant := (newGrid[c.y][c.x] == TileMiasma)
		enemy := generateEnemy(m.player.Stats.Level+m.floor, isMutant, TypeNormal)
		enemy.X, enemy.Y = c.x, c.y
		m.enemies = append(m.enemies, enemy)
		spawnIndex++
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

	// 0. Curar com Poção (Tecla 'h')
	if msg.String() == "h" {
		if m.player.Stats.Potions > 0 && m.player.Stats.HP < m.player.Stats.MaxHP {
			m.player.Stats.Potions--
			cura := m.player.Stats.PotionPower
			m.player.Stats.HP = math.Min(m.player.Stats.HP+cura, m.player.Stats.MaxHP)
			m.log = fmt.Sprintf("🧪 Bebeu poção! Curou %.0f HP. Restam: %d", cura, m.player.Stats.Potions)
		} else if m.player.Stats.Potions <= 0 {
			m.log = "⚠️ Sem poções!"
		} else {
			m.log = "Sua vida já está cheia."
		}
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
		m.state = StateShop
		m.menuCursor = 0
		m.log = "Você encontrou o Mercador do Abismo..."
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
