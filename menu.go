package main

import (
	"math/rand"
)

func startGame(classIndex, elementIndex int) model {
	m := model{state: StateMap}
	c := AvailableClasses[classIndex]

	element := AvailableElements[elementIndex].Name

	m.player = &Character{
		Name:    c.Name,
		Symbol:  c.Symbol,
		Stats:   Stats{MaxHP: c.HP, HP: c.HP, Str: c.Str, Int: c.Int, Level: 1, XP: 0, NextXP: 50},
		Skills:  []Skill{&HeavyStrike{BaseSkill{name: "Ataque Base", level: 1}}},
		Element: element,
	}

	// 1. Gera o Chão e Paredes
	for y := 0; y < MapH; y++ {
		for x := 0; x < MapW; x++ {
			if x == 0 || x == MapW-1 || y == 0 || y == MapH-1 || rand.Float64() < 0.15 {
				m.grid[y][x] = TileWall
			} else {
				m.grid[y][x] = TileFloor
			}
		}
	}

	// 2. Gera "Poças" de Miasma
	for i := 0; i < 8; i++ {
		mx, my := rand.Intn(MapW-10)+5, rand.Intn(MapH-10)+5
		for dy := -3; dy <= 3; dy++ {
			for dx := -3; dx <= 3; dx++ {
				if rand.Float64() < 0.6 {
					m.grid[my+dy][mx+dx] = TileMiasma
				}
			}
		}
	}

	// Limpa área segura inicial do jogador
	sx, sy := MapW/2, MapH/2
	for y := sy - 2; y <= sy+2; y++ {
		for x := sx - 2; x <= sx+2; x++ {
			m.grid[y][x] = TileFloor
		}
	}

	// 3. Spawna Baús e Inimigos
	for i := 0; i < 10; i++ { // Baús
		cx, cy := rand.Intn(MapW-2)+1, rand.Intn(MapH-2)+1
		if m.grid[cy][cx] != TileWall {
			m.grid[cy][cx] = TileChest
		}
	}

	for i := 0; i < 35; i++ { // Inimigos
		ex, ey := rand.Intn(MapW-2)+1, rand.Intn(MapH-2)+1
		tileAtual := m.grid[ey][ex]

		if tileAtual == TileFloor || tileAtual == TileMiasma {
			isMutant := (tileAtual == TileMiasma)

			// Gera o inimigo e anota onde ele está
			enemy := generateEnemy(m.player.Stats.Level, isMutant)
			enemy.X = ex
			enemy.Y = ey
			m.enemies = append(m.enemies, enemy)

			// Marca no grid
			if isMutant {
				m.grid[ey][ex] = 'M'
			} else {
				m.grid[ey][ex] = 'E'
			}
		}
	}

	m.playerX, m.playerY = sx, sy
	m.log = "A jornada começa. Cuidado com as zonas roxas!"
	return m
}
