package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	tasksets "tuido/tasksets"
)

func main() {
	tasksets.CreatePages()
	home := tasksets.GetHomePage()
	p := tea.NewProgram(home, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Panic("Failed to start:", err.Error())
	}
}
