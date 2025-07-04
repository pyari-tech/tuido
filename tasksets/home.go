package tasksets

/*
tasksets::home
*/

import (
	"time"
	"tuido/persist"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const sections = 4
const todoListTitle = "ToDo"
const wipListTitle = "W.I.P."
const doneListTitle = "Done"

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
	selected    Status
	lists       []list.Model
	loaded      bool
	exiting     bool
	err         error
	updateIndex int
}

func NewHome() *Home {
	return &Home{updateIndex: -1}
}

func initList(title string, width, height int) list.Model {
	var lst = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	lst.SetShowHelp(false)
	lst.Title = title
	return lst
}

func (h *Home) initLists(width, height int) {
	if isLoaded := h.Load(width, height); isLoaded {
		return
	}
	var todoModel = initList(todoListTitle, width, height)
	var wipModel = initList(wipListTitle, width, height)
	var doneModel = initList(doneListTitle, width, height)
	h.lists = []list.Model{
		todoModel,
		wipModel,
		doneModel,
	}
	h.lists[0].InsertItem(
		0,
		Task{
			status:      todo,
			title:       "create your own todos",
			description: "this is an empty board, edit/delete this, add new ones",
		},
	)
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
		case "+", "n":
			pages[homePage] = h
			pages[taskFormPage] = NewTaskForm()
			ctrlN := tea.KeyMsg(tea.Key{Type: tea.KeyCtrlN, Runes: []rune{'c', 't', 'r', 'l', '+', 'n'}})
			return pages[taskFormPage].Update(ctrlN)
		case "-", "delete":
			h.DeleteTask()
			return h, nil
		case "insert", "u":
			h.UpdateTask()
			ctrlU := tea.KeyMsg(tea.Key{Type: tea.KeyCtrlU, Runes: []rune{'c', 't', 'r', 'l', '+', 'u'}})
			return pages[taskFormPage].Update(ctrlU)
		case "tab", "r":
			h.UpdateTask()
			ctrlR := tea.KeyMsg(tea.Key{Type: tea.KeyCtrlR, Runes: []rune{'c', 't', 'r', 'l', '+', 'r'}})
			return pages[taskFormPage].Update(ctrlR)
		case "ctrl+c", "q":
			h.exiting = true
			h.Persist()
			return h, tea.Quit
		}
	case Task:
		tsk := msg
		tskIndex := len(h.lists[todo].Items())
		if h.updateIndex > -1 {
			tskIndex = h.updateIndex
			h.lists[tsk.status].RemoveItem(tskIndex)
		}
		tsk.updated = time.Now()
		cmd := h.lists[tsk.status].InsertItem(tskIndex, tsk)
		h.updateIndex = -1
		return h, cmd
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

func (h *Home) currentTask() (Task, int) {
	selectedItem := h.lists[h.selected].SelectedItem()
	selectedTask, ok := selectedItem.(Task)
	if !ok {
		return Task{}, -1
	}
	selectedItemIndex := h.lists[h.selected].Index()
	return selectedTask, selectedItemIndex
}

func (h *Home) deleteTask() Task {
	selectedTask, selectedItemIndex := h.currentTask()
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
	selectedTask.updated = time.Now()
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

func (h *Home) UpdateTask() {
	tisTask, tisIndex := h.currentTask()
	h.updateIndex = tisIndex
	pages[homePage] = h
	tf := NewTaskForm()
	tf.status = h.selected
	tf.title.SetValue(tisTask.title)
	tf.description.SetValue(tisTask.description)
	tf.title.Focus()
	tf.title.CursorStart()
	pages[taskFormPage] = tf
}

func (h *Home) Persist() {
	board := persist.Tuido{
		Lists: make([]persist.TuidoList, 3),
	}
	for idx, lst := range h.lists {
		items := lst.Items()
		board.Lists[idx].Title = lst.Title
		board.Lists[idx].Tasks = make([]persist.Task, len(items))
		for tskIdx, item := range items {
			if tsk, ok := item.(Task); ok {
				board.Lists[idx].Tasks[tskIdx].Index = tskIdx
				board.Lists[idx].Tasks[tskIdx].Title = tsk.title
				board.Lists[idx].Tasks[tskIdx].Description = tsk.description
				board.Lists[idx].Tasks[tskIdx].Created = persist.CustomTime{Time: tsk.created}
				board.Lists[idx].Tasks[tskIdx].Updated = persist.CustomTime{Time: tsk.updated}
			}
		}
	}
	board.Persist(TuidoFile)
}

func (h *Home) Load(width, height int) bool {
	var savedBoard = persist.LoadTuido(TuidoFile)
	if len(savedBoard.Lists) == 0 {
		return false
	}
	h.lists = make([]list.Model, len(savedBoard.Lists))

	for idx, savedList := range savedBoard.Lists {
		var lst = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
		lst.Title = savedList.Title
		var listStatus Status
		switch lst.Title {
		case todoListTitle:
			listStatus = todo
		case wipListTitle:
			listStatus = wip
		case doneListTitle:
			listStatus = done
		}
		for _, item := range savedList.Tasks {
			tsk := Task{
				status:      listStatus,
				title:       item.Title,
				description: item.Description,
				created:     item.Created.Time,
				updated:     item.Updated.Time,
			}
			lst.InsertItem(item.Index, tsk)
		}
		lst.SetShowHelp(false)
		h.lists[idx] = lst
	}
	h.lists[0].SetShowStatusBar(true)
	return true
}
