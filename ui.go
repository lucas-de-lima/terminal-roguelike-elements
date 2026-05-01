package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// --- 0. ESTILOS E EMOJIS ---

const (
	EmojiWarrior = "⚔️ "
	EmojiMage    = "🧙"
	EmojiRogue   = "🥷"
	EmojiEnemy   = "👹"
	EmojiWall    = "🧱 "
	EmojiFloor   = " . "
	EmojiChest   = "📦 "
	EmojiMiasma  = "🟣 "
	EmojiPortal  = "🌀 "
)

var (
	colorPlayer = lipgloss.Color("#00FF99")
	colorEnemy  = lipgloss.Color("#FF3333")
	colorMenu   = lipgloss.Color("#FFD700")

	stylePlayer = lipgloss.NewStyle().Foreground(colorPlayer)
	styleEnemy  = lipgloss.NewStyle().Foreground(colorEnemy)
	styleWall   = lipgloss.NewStyle()
	styleFloor  = lipgloss.NewStyle().Foreground(lipgloss.Color("#333333"))
	styleFog    = lipgloss.NewStyle().Foreground(lipgloss.Color("#0a0a0a"))

	// Novos estilos
	styleChest  = lipgloss.NewStyle()
	styleMiasma = lipgloss.NewStyle().Foreground(lipgloss.Color("#A020F0")) // Roxo Tóxico
	stylePortal = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Bold(true)

	styleHUD = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#555555")).
			Padding(0, 1).
			MarginBottom(1)

	styleLog    = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).Italic(true)
	styleTitle  = lipgloss.NewStyle().Foreground(colorMenu).Bold(true).MarginBottom(1)
	styleCursor = lipgloss.NewStyle().Foreground(colorMenu).Bold(true)
)

func progressBar(width int, current, total float64, color lipgloss.Color) string {
	percent := math.Max(0, math.Min(1, current/total))
	filled := int(float64(width) * percent)
	return lipgloss.NewStyle().Foreground(color).Render(strings.Repeat("█", filled) + strings.Repeat("░", width-filled))
}

func (m model) View() string {
	if m.state == StateMainMenu {
		s := "\n" + styleTitle.Render(" ⚔️  BEM-VINDO AO TERMINAL RPG ⚔️ ") + "\n\n"

		if m.menuStep == 0 {
			s += " Passo 1: Escolha sua vocação (W/S para mover, ENTER para confirmar):\n\n"
			for i, class := range AvailableClasses {
				cursor := "  "
				if m.menuCursor == i {
					cursor = styleCursor.Render("> ")
				}
				s += fmt.Sprintf("%s %s %s\n     HP: %.0f | FOR: %.0f | INT: %.0f\n     %s\n\n",
					cursor, class.Symbol, class.Name, class.HP, class.Str, class.Int, class.Desc)
			}
		} else {
			s += fmt.Sprintf(" Vocação: %s %s\n", AvailableClasses[m.chosenClass].Symbol, AvailableClasses[m.chosenClass].Name)
			s += " Passo 2: Escolha sua afinidade elemental:\n\n"
			for i, elem := range AvailableElements {
				cursor := "  "
				if m.menuCursor == i {
					cursor = styleCursor.Render("> ")
				}
				s += fmt.Sprintf("%s %s %s\n", cursor, elem.Symbol, elem.Name)
			}
			s += "\n (Pressione ESC para voltar)"
		}
		return s
	}

	hpBar := progressBar(15, m.player.Stats.HP, m.player.Stats.MaxHP, lipgloss.Color("#FF0055"))
	xpBar := progressBar(10, m.player.Stats.XP, m.player.Stats.NextXP, lipgloss.Color("#AF87FF"))

	// --- CÁLCULO DA BÚSSOLA ---
	trackerText := ""
	limiteBussola := int(math.Ceil(float64(m.totalEnemies) * 0.10)) // Calcula 10% (arredondado pra cima)

	if len(m.enemies) > 0 && (len(m.enemies) <= limiteBussola || len(m.enemies) <= 5) {
		minDist := 9999.0
		var closest *Character
		for _, e := range m.enemies {
			// Distância Euclidiana
			d := math.Hypot(float64(e.X-m.playerX), float64(e.Y-m.playerY))
			if d < minDist {
				minDist = d
				closest = e
			}
		}

		// Calcula o ângulo em Graus
		dx := float64(closest.X - m.playerX)
		dy := float64(closest.Y - m.playerY)
		angle := math.Atan2(dy, dx) * 180 / math.Pi

		arrow := ""
		switch {
		case angle > -22.5 && angle <= 22.5:
			arrow = "➡️"
		case angle > 22.5 && angle <= 67.5:
			arrow = "↘️"
		case angle > 67.5 && angle <= 112.5:
			arrow = "⬇️"
		case angle > 112.5 && angle <= 157.5:
			arrow = "↙️"
		case angle > 157.5 || angle <= -157.5:
			arrow = "⬅️"
		case angle > -157.5 && angle <= -112.5:
			arrow = "↖️"
		case angle > -112.5 && angle <= -67.5:
			arrow = "⬆️"
		case angle > -67.5 && angle <= -22.5:
			arrow = "↗️"
		}
		trackerText = fmt.Sprintf(" | 🧭 Presa a %d passos %s", int(minDist), arrow)
	}

	// --- HUD ---
	hudText := fmt.Sprintf(" Andar %d | %s %s Nv.%d | HP %s %.0f/%.0f | XP %s | 🪙 %d | 🧪 %d/%d %s",
		m.floor, m.player.Symbol, m.player.Name, m.player.Stats.Level,
		hpBar, m.player.Stats.HP, m.player.Stats.MaxHP, xpBar,
		m.player.Stats.Gold, m.player.Stats.Potions, m.player.Stats.MaxPotions, trackerText)

	view := styleHUD.Render(hudText) + "\n"

	currentFloorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(m.currentBiome.ColorFloor))
	currentHazardStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(m.currentBiome.ColorHazard))

	if m.state == StateMap {
		camX, camY := m.playerX-(ViewW/2), m.playerY-(ViewH/2)

		for y := 0; y < ViewH; y++ {
			line := ""
			for x := 0; x < ViewW; x++ {
				wx, wy := camX+x, camY+y
				dist := math.Sqrt(math.Pow(float64(wx-m.playerX), 2) + math.Pow(float64(wy-m.playerY), 2))

				if wx < 0 || wx >= MapW || wy < 0 || wy >= MapH {
					line += "   "
				} else if dist > VisionRadius {
					line += styleFog.Render(" . ")
				} else {
					// 1. Verifica Entidades primeiro (Jogador e Inimigos)
					hasEnemy := false
					isMutant := false
					emoji := ""
					for _, e := range m.enemies {
						if e.X == wx && e.Y == wy {
							hasEnemy = true
							emoji = e.Symbol
							isMutant = strings.Contains(e.Name, "Tóxico")
							break
						}
					}

					if wx == m.playerX && wy == m.playerY {
						line += stylePlayer.Render(m.player.Symbol + " ")
					} else if hasEnemy {
						if isMutant {
							line += currentHazardStyle.Render(emoji + " ")
						} else {
							line += styleEnemy.Render(emoji + " ")
						}
						// 2. Se não houver entidade, renderiza o terreno (grid)
					} else if m.grid[wy][wx] == TileWall {
						line += styleWall.Render(m.currentBiome.EmojiWall)
					} else if m.grid[wy][wx] == TileMiasma {
						line += currentHazardStyle.Render(m.currentBiome.EmojiHazard + " ")
					} else if m.grid[wy][wx] == TileChest {
						line += styleChest.Render(EmojiChest + " ")
					} else if m.grid[wy][wx] == TilePortal {
						line += stylePortal.Render(EmojiPortal + " ")
					} else {
						line += currentFloorStyle.Render(" . ")
					}
				}
			}
			view += line + "\n"
		}
		view += "\n" + styleLog.Render(m.log)

	} else if m.state == StateCombat {
		enemyHP := progressBar(20, m.enemy.Stats.HP, m.enemy.Stats.MaxHP, colorEnemy)
		combatView := fmt.Sprintf("\n %s COMBATE: %s\n HP: %s\n\n", m.enemy.Symbol, m.enemy.Name, enemyHP)
		for i, s := range m.player.Skills {
			combatView += fmt.Sprintf(" [%d] %s\n", i+1, s.Name())
		}
		view += lipgloss.NewStyle().Padding(2).Render(combatView)
		view += "\n\n" + styleLog.Render(m.log)

	} else if m.state == StateLevelUp {
		lvlView := "\n ✨ ESCOLHA UMA RECOMPENSA ✨\n\n"
		for i, opt := range m.levelOptions {
			lvlView += fmt.Sprintf(" [%d] %s\n     %s\n\n", i+1, opt.Name(), opt.Description())
		}
		view += lipgloss.NewStyle().Padding(2).Foreground(lipgloss.Color("#AF87FF")).Render(lvlView)
	} else if m.state == StateMilestone {
		msView := "\n 🌟 MARCO ATINGIDO! ESCOLHA UMA MELHORIA 🌟\n\n"
		for i, reward := range m.milestoneOptions {
			msView += fmt.Sprintf(" [%d] %s\n     %s\n\n", i+1, reward.Name, reward.Desc)
		}
		view += lipgloss.NewStyle().Padding(2).Foreground(lipgloss.Color("#FFD700")).Render(msView)
	} else if m.state == StateShop {
		s := "\n ⛺ O MERCADOR DO ABISMO ⛺\n\n"
		s += fmt.Sprintf(" Seu Ouro: 🪙 %d\n\n", m.player.Stats.Gold)

		opcoes := []string{
			"Recarregar Poções (50g)",
			"Comprar +1 Frasco Máximo (150g)",
			"Aumentar Potência da Cura (200g)",
			"Continuar descendo o Abismo",
		}

		for i, opt := range opcoes {
			cursor := "  "
			if m.menuCursor == i {
				cursor = styleCursor.Render("> ")
			}
			s += fmt.Sprintf("%s %s\n", cursor, opt)
		}
		view += s + "\n\n" + styleLog.Render(m.log)
	} else if m.state == StateGameOver {
		view += "\n\n 💀 VOCÊ MORREU 💀\n (Pressione R para ir ao Menu)"
	}

	return view
}
