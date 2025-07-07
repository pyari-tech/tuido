package tasksets

/*
tasksets::home
*/

import (
	"time"
	"tuido/persist"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const sections = 4

var listTitles = map[Status]string{
	todo:    "ToDo",
	wip:     "W.I.P.",
	done:    "Done",
	blocked: "Blocked",
	abandon: "Abandon",
	archive: "Archive",
}

/* STYLING */
const padLeftRight = 3
const padTopBottom = 1

var (
	columnStyle = lipgloss.NewStyle().Padding(padTopBottom, padLeftRight)

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
	lastColumn  Status
}

func NewHome() *Home {
	return &Home{updateIndex: -1, lastColumn: done}
}

func initList(title string, width, height int) list.Model {
	var lst = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	lst.SetShowHelp(false)
	lst.Title = title
	return lst
}

func (h *Home) initLists(width, height int) {
	var todoModel = initList(listTitles[todo], width, height)
	var wipModel = initList(listTitles[wip], width, height)
	var doneModel = initList(listTitles[done], width, height)
	var blockedModel = initList(listTitles[blocked], width, height)
	var abandonModel = initList(listTitles[abandon], width, height)
	var archiveModel = initList(listTitles[archive], width, height)
	h.lists = []list.Model{
		todoModel,
		wipModel,
		doneModel,
		blockedModel,
		abandonModel,
		archiveModel,
	}
	if isLoaded := h.Load(width, height); !isLoaded {
		h.lists[0].InsertItem(0, InitTask())
	}
	h.lists[0].SetShowStatusBar(true)
	h.lists[0].SetShowHelp(true)
	h.lists[0].AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.selectListOnLeft,
			listKeys.selectListOnRight,
			listKeys.nextListPage,
			listKeys.prevListPage,
			listKeys.moveTaskUp,
			listKeys.moveTaskToNextList,
			listKeys.moveTaskDown,
			listKeys.moveTaskToPrevList,
			listKeys.createTask,
			listKeys.readTask,
			listKeys.updateTask,
			listKeys.deleteTask,
			listKeys.showDoneLast,
			listKeys.showBlockedLast,
			listKeys.showAbandonLast,
			listKeys.showArchiveLast,
		}
	}
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
		case "right", "l":
			h.NextList()
		case "pageup":
			h.lists[h.selected].PrevPage()
		case "pagedown":
			h.lists[h.selected].NextPage()
			return h, nil
		case "[", "w":
			return h, h.MoveTaskUp()
		case "<", "a":
			return h, h.MoveTaskToPrev()
		case "]", "s":
			return h, h.MoveTaskDown()
		case ">", "d":
			return h, h.MoveTaskToNext()
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
		case "ctrl+e":
			h.lastColumn = done
			return h, nil
		case "ctrl+r":
			h.lastColumn = archive
			return h, nil
		case "ctrl+t":
			h.lastColumn = blocked
			return h, nil
		case "ctrl+y":
			h.lastColumn = abandon
			return h, nil
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
		lastView string
	)

	if h.lastColumn == blocked {
		lastView = h.columnView(blocked)
	} else if h.lastColumn == abandon {
		lastView = h.columnView(abandon)
	} else if h.lastColumn == archive {
		lastView = h.columnView(archive)
	} else {
		lastView = h.columnView(done)
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		todoView,
		wipView,
		lastView,
	)
}

func (h *Home) NextList() {
	targetStatus := h.selected + 1
	if h.selected == h.lastColumn {
		h.selected = todo
	} else if targetStatus < done {
		h.selected = targetStatus
	} else {
		h.selected = h.lastColumn
	}
}

func (h *Home) PrevList() {
	if h.selected == todo {
		h.selected = h.lastColumn
	} else if h.selected == h.lastColumn {
		h.selected = wip
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
	h.lists[h.selected].Select(selectedItemTargetIndex)
}

func (h *Home) changeTaskOrder(step int) {
	var indexToSwitch = h.lists[h.selected].Index() + step
	if step > 0 {
		maxIndexAvailable := len(h.lists[h.selected].Items()) - 1
		if indexToSwitch > maxIndexAvailable {
			return
		}
	} else if indexToSwitch < 0 {
		return
	}
	selectedTask := h.deleteTask()
	selectedTask.updated = time.Now()
	h.lists[h.selected].InsertItem(indexToSwitch, selectedTask)
	h.lists[h.selected].Select(indexToSwitch)
}

func (h *Home) DeleteTask() {
	h.deleteTask()
}

func (h *Home) MoveTaskToNext() tea.Cmd {
	if h.selected == h.lastColumn {
		return nil
	} else if len(h.lists[h.selected].Items()) == 0 {
		return nil
	}

	targetStatus := h.selected + 1
	if targetStatus < done {
		h.changeTaskStatus(targetStatus)
	} else {
		h.changeTaskStatus(h.lastColumn)
	}
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

func (h *Home) MoveTaskDown() tea.Cmd {
	h.changeTaskOrder(1)
	return nil
}

func (h *Home) MoveTaskUp() tea.Cmd {
	h.changeTaskOrder(-1)
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
		Lists: make([]persist.TuidoList, totalStatus),
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

	for idx, savedList := range savedBoard.Lists {
		var lst = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
		lst.Title = savedList.Title
		var listStatus Status
		switch lst.Title {
		case listTitles[todo]:
			listStatus = todo
		case listTitles[wip]:
			listStatus = wip
		case listTitles[done]:
			listStatus = done
		case listTitles[blocked]:
			listStatus = blocked
		case listTitles[abandon]:
			listStatus = abandon
		case listTitles[archive]:
			listStatus = archive
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
	return true
}
