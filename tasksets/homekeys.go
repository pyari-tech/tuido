package tasksets

import (
	"github.com/charmbracelet/bubbles/key"
)

type listKeyMap struct {
	selectListOnLeft   key.Binding
	selectListOnRight  key.Binding
	nextListPage       key.Binding
	prevListPage       key.Binding
	moveTaskUp         key.Binding
	moveTaskToNextList key.Binding
	moveTaskDown       key.Binding
	moveTaskToPrevList key.Binding
	createTask         key.Binding
	readTask           key.Binding
	updateTask         key.Binding
	deleteTask         key.Binding
	showDoneLast       key.Binding
	showBlockedLast    key.Binding
	showAbandonLast    key.Binding
	showArchiveLast    key.Binding
}

var (
	listKeys = newListKeyMap()
)

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		selectListOnLeft: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("left", "or 'h' to move left")),
		selectListOnRight: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("right", "or 'l' to move right")),
		nextListPage: key.NewBinding(
			key.WithKeys("pageup"),
			key.WithHelp("pageup", "previous list page")),
		prevListPage: key.NewBinding(
			key.WithKeys("pagedown"),
			key.WithHelp("pagedown", "next list page")),
		moveTaskUp: key.NewBinding(
			key.WithKeys("w", "["),
			key.WithHelp("w", "or '[' move task up")),
		moveTaskToNextList: key.NewBinding(
			key.WithKeys("a", "<"),
			key.WithHelp("a", "or '<' move task left")),
		moveTaskDown: key.NewBinding(
			key.WithKeys("s", "]"),
			key.WithHelp("s", "or ']' move task down")),
		moveTaskToPrevList: key.NewBinding(
			key.WithKeys("d", ">"),
			key.WithHelp("d", "or '>' move task right")),
		createTask: key.NewBinding(
			key.WithKeys("n", "+"),
			key.WithHelp("+", "Add task")),
		readTask: key.NewBinding(
			key.WithKeys("delete", "-"),
			key.WithHelp("-", "Delete task")),
		updateTask: key.NewBinding(
			key.WithKeys("insert", "u"),
			key.WithHelp("u", "Update task")),
		deleteTask: key.NewBinding(
			key.WithKeys("tab", "r"),
			key.WithHelp("tab", "Read task")),
		showDoneLast: key.NewBinding(
			key.WithKeys("ctrl+e"),
			key.WithHelp("ctrl+e", "Done list")),
		showBlockedLast: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("ctrl+t", "Blocked list")),
		showAbandonLast: key.NewBinding(
			key.WithKeys("ctrl+y"),
			key.WithHelp("ctrl+y", "Abandon list")),
		showArchiveLast: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("ctrl+r", "Archive list")),
	}
}
