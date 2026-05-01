package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Erro:", err)
		os.Exit(1)
	}
}
