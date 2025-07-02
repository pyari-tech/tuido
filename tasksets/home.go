package tasksets

/*
tasksets::home
*/

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

/* Struct Home for homePage's Model */
type Home struct {
	selected Status
	lists    []list.Model
	loaded   bool
	exiting  bool
	err      error
}

func NewHome() *Home {
	return &Home{}
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

func (h *Home) initLists(width, height int) {
	var todo = initList("ToDo", todo, width, height)
	var wip = initList("W.I.P.", wip, width, height)
	var done = initList("Done", done, width, height)
	h.lists = []list.Model{
		todo,
		wip,
		done,
	}
	h.lists[0].SetShowStatusBar(true)
}

func (h Home) Init() tea.Cmd {
	return nil
}

func (h Home) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !h.loaded {
			width := msg.Width / sections
			height := msg.Height - int(msg.Height/10)
			h.initLists(width, height)
			h.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			h.PrevList()
		case "right", "j":
			h.NextList()
		case ">", "k":
			return h, h.MoveTaskToNext()
		case "<", "l":
			return h, h.MoveTaskToPrev()
		case "+", "n", "insert":
			pages[homePage] = h
			pages[taskFormPage] = NewTaskForm()
			ctrlN := tea.KeyMsg(tea.Key{Type: tea.KeyCtrlN, Runes: []rune{'c', 't', 'r', 'l', '+', 'n'}})
			return pages[taskFormPage].Update(ctrlN)
		case "-", "delete":
			h.DeleteTask()
			return h, nil
		case "ctrl+c", "q":
			h.exiting = true
			return h, tea.Quit
		}
	case Task:
		tsk := msg
		todoIndex := len(h.lists[todo].Items())
		return h, h.lists[tsk.status].InsertItem(todoIndex, tsk)
	}
	var cmd tea.Cmd
	h.lists[h.selected], cmd = h.lists[h.selected].Update(msg)
	return h, cmd
}

func (h Home) columnView(status Status) string {
	if h.selected == status {
		return selectedStyle.Render(h.lists[status].View())
	}
	return columnStyle.Render(h.lists[status].View())
}

func (h Home) View() string {
	if h.exiting {
		return "[TBD] clean exit with persist state"
	}
	if !h.loaded {
		return "loading..."
	}

	var (
		todoView = h.columnView(todo)
		wipView  = h.columnView(wip)
		doneView = h.columnView(done)
	)
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		todoView,
		wipView,
		doneView,
	)
}

func (h *Home) NextList() {
	if h.selected == done {
		h.selected = todo
	} else {
		h.selected++
	}
}

func (h *Home) PrevList() {
	if h.selected == todo {
		h.selected = done
	} else {
		h.selected--
	}
}

func (h *Home) deleteTask() Task {
	selectedItem := h.lists[h.selected].SelectedItem()
	selectedTask, ok := selectedItem.(Task)
	if !ok {
		return Task{}
	}
	selectedItemIndex := h.lists[h.selected].Index()
	h.lists[selectedTask.status].RemoveItem(selectedItemIndex)
	return selectedTask
}

func (h *Home) changeTaskStatus(targetStatus Status) {
	selectedTask := h.deleteTask()
	if h.selected > targetStatus {
		h.PrevList()
	} else {
		h.NextList()
	}
	selectedTask.status = h.selected
	selectedItemTargetIndex := len(h.lists[h.selected].Items())
	h.lists[h.selected].InsertItem(selectedItemTargetIndex, selectedTask)
}

func (h *Home) DeleteTask() {
	h.deleteTask()
}

func (h *Home) MoveTaskToNext() tea.Cmd {
	if h.selected == done {
		return nil
	} else if len(h.lists[h.selected].Items()) == 0 {
		return nil
	}

	h.changeTaskStatus(h.selected + 1)
	return nil
}

func (h *Home) MoveTaskToPrev() tea.Cmd {
	if h.selected == todo {
		return nil
	} else if len(h.lists[h.selected].Items()) == 0 {
		return nil
	}

	h.changeTaskStatus(h.selected - 1)
	return nil
}
