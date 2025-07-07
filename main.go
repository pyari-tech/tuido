package main

import (
	"flag"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"

	tasksets "tuido/tasksets"
)

var (
	license   = flag.Bool("license", false, "About license")
	tuidoFile = flag.String("file", "tuido.yaml", "YAML file path to load/persist Tuido Board")
)

func main() {
	flag.Parse()
	if *license {
		showLicense()
		return
	}
	tasksets.TuidoFile = *tuidoFile
	tasksets.CreatePages()
	home := tasksets.GetHomePage()
	p := tea.NewProgram(home, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Panic("Failed to start:", err.Error())
	}
}

func showLicense() {
	fmt.Println(`
	tuido  Copyright (C) 2025-infinity AbhishekKr
	This program comes with ABSOLUTELY NO WARRANTY.
	It's a free software available under GPLv3 License.

	For details visit: https://www.gnu.org/licenses/
			`)
}
