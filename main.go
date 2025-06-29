package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	tasksets "tuido/tasksets"
)

func main() {
	m := tasksets.NewHome()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Panic("Failed to start:", err.Error())
	}
}
