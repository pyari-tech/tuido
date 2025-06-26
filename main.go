package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Status int8

const (
	todo Status = iota
	wip
	done
	blocked
	abandon
)

const sections = 4

/*
Task would implement `list` with Title, Description, FilterValue
*/
type Task struct {
	status      Status
	title       string
	description string
}

func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

type Model struct {
	selected Status
	lists    []list.Model
	loaded   bool
	err      error
}

func New() *Model {
	return &Model{}
}

func initList(title string, status Status, width, height int) list.Model {
	var lst = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	lst.Title = title
	lst.SetItems([]list.Item{
		Task{status: status, title: fmt.Sprintf("%s this", title), description: "this is this"},
		Task{status: status, title: fmt.Sprintf("%s that", title), description: "this is that"},
		Task{status: status, title: fmt.Sprintf("%s what", title), description: "this is what"},
	})
	return lst
}

func (m *Model) initLists(width, height int) {
	width = width / sections
	var todo = initList("ToDo", todo, width, height)
	var wip = initList("W.I.P.", wip, width, height)
	var done = initList("Done", done, width, height)
	m.lists = []list.Model{
		todo,
		wip,
		done,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}
	}
	var cmd tea.Cmd
	m.lists[m.selected], cmd = m.lists[m.selected].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if !m.loaded {
		return "loading..."
	}
	return lipgloss.JoinHorizontal(lipgloss.Left,
		m.lists[todo].View(),
		m.lists[wip].View(),
		m.lists[done].View())
}

func main() {
	m := New()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("done.")
}
