package tasksets

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const sections = 4

/* STYLING */
const padLeftRight = 3
const padTopBottom = 1

var (
	columnStyle   = lipgloss.NewStyle().Padding(padTopBottom, padLeftRight)
	selectedStyle = lipgloss.NewStyle().Padding(padTopBottom, padLeftRight).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

/* Struct */
type Model struct {
	selected Status
	lists    []list.Model
	loaded   bool
	exiting  bool
	err      error
}

func NewModel() *Model {
	return &Model{}
}

func initList(title string, status Status, width, height int) list.Model {
	var lst = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	lst.SetShowHelp(false)
	lst.Title = title
	lst.SetItems([]list.Item{
		Task{status: status, title: fmt.Sprintf("%s-1", title), description: "#1"},
		Task{status: status, title: fmt.Sprintf("%s-2", title), description: "#2"},
		Task{status: status, title: fmt.Sprintf("%s-3", title), description: "#3"},
	})
	return lst
}

func (m *Model) initLists(width, height int) {
	var todo = initList("ToDo", todo, width, height)
	var wip = initList("W.I.P.", wip, width, height)
	var done = initList("Done", done, width, height)
	m.lists = []list.Model{
		todo,
		wip,
		done,
	}
	m.lists[0].SetShowStatusBar(true)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			width := msg.Width / sections
			height := msg.Height - int(msg.Height/10)
			m.initLists(width, height)
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			m.PrevList()
		case "right", "j":
			m.NextList()
		case ">", "k":
			return m, m.MoveTaskToNext()
		case "<", "l":
			return m, m.MoveTaskToPrev()
		case "ctrl+c", "q":
			m.exiting = true
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.lists[m.selected], cmd = m.lists[m.selected].Update(msg)
	return m, cmd
}

func (m Model) columnView(status Status) string {
	if m.selected == status {
		return selectedStyle.Render(m.lists[status].View())
	}
	return columnStyle.Render(m.lists[status].View())
}

func (m Model) View() string {
	if m.exiting {
		return "[TBD] clean exit with persist state"
	}
	if !m.loaded {
		return "loading..."
	}

	var (
		todoView = m.columnView(todo)
		wipView  = m.columnView(wip)
		doneView = m.columnView(done)
	)
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		todoView,
		wipView,
		doneView,
	)
}

func (m *Model) NextList() {
	if m.selected == done {
		m.selected = todo
	} else {
		m.selected++
	}
}

func (m *Model) PrevList() {
	if m.selected == todo {
		m.selected = done
	} else {
		m.selected--
	}
}

func (m *Model) changeTaskStatus(targetStatus Status) {
	selectedItem := m.lists[m.selected].SelectedItem()
	selectedTask, ok := selectedItem.(Task)
	if !ok {
		return
	}
	selectedItemIndex := m.lists[m.selected].Index()
	m.lists[selectedTask.status].RemoveItem(selectedItemIndex)
	if m.selected > targetStatus {
		m.PrevList()
	} else {
		m.NextList()
	}
	selectedTask.status = m.selected
	selectedItemTargetIndex := len(m.lists[m.selected].Items()) + 1
	m.lists[m.selected].InsertItem(selectedItemTargetIndex, selectedTask)
}

func (m *Model) MoveTaskToNext() tea.Cmd {
	if m.selected == done {
		return nil
	} else if len(m.lists[m.selected].Items()) == 0 {
		return nil
	}

	m.changeTaskStatus(m.selected + 1)
	return nil
}

func (m *Model) MoveTaskToPrev() tea.Cmd {
	if m.selected == todo {
		return nil
	} else if len(m.lists[m.selected].Items()) == 0 {
		return nil
	}

	m.changeTaskStatus(m.selected - 1)
	return nil
}
