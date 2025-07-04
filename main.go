package main

import (
	"flag"
	"log"

	tea "github.com/charmbracelet/bubbletea"

	tasksets "tuido/tasksets"
)

var (
	tuidoFile = flag.String("file", "tuido.yaml", "YAML file path to load/persist Tuido Board")
)

func main() {
	flag.Parse()
	tasksets.TuidoFile = *tuidoFile
	tasksets.CreatePages()
	home := tasksets.GetHomePage()
	p := tea.NewProgram(home, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Panic("Failed to start:", err.Error())
	}
}
