package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func updateShop(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "up", "w":
		m.menuCursor--
		if m.menuCursor < 0 {
			m.menuCursor = 3
		}
	case "down", "s":
		m.menuCursor++
		if m.menuCursor > 3 {
			m.menuCursor = 0
		}
	case "enter", " ":
		switch m.menuCursor {
		case 0: // Refil de Poções (50g)
			if m.player.Stats.Gold >= 50 && m.player.Stats.Potions < m.player.Stats.MaxPotions {
				m.player.Stats.Gold -= 50
				m.player.Stats.Potions = m.player.Stats.MaxPotions
				m.log = "Frascos reabastecidos!"
			} else {
				m.log = "Ouro insuficiente ou frascos já cheios."
			}
		case 1: // Upgrade: +1 Frasco Máximo (150g)
			if m.player.Stats.Gold >= 150 {
				m.player.Stats.Gold -= 150
				m.player.Stats.MaxPotions++
				m.player.Stats.Potions++
				m.log = "Você comprou um Frasco Extra!"
			} else {
				m.log = "Ouro insuficiente."
			}
		case 2: // Upgrade: Potência da Poção (200g)
			if m.player.Stats.Gold >= 200 {
				m.player.Stats.Gold -= 200
				m.player.Stats.PotionPower += 20.0
				m.log = "Suas poções agora curam mais vida!"
			} else {
				m.log = "Ouro insuficiente."
			}
		case 3: // Descer para o próximo andar
			m.floor++
			generateFloor(&m)
			m.state = StateMap
			m.log = fmt.Sprintf("Você desceu para o Andar %d.", m.floor)
		}
	}
	return m, nil
}
