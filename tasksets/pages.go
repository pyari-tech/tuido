package tasksets

/*
tasksets::pages
*/

import tea "github.com/charmbracelet/bubbletea"

var (
	TuidoFile = ""
)

/* Tea Model */
type pageType int8

var pages []tea.Model

const (
	homePage pageType = iota
	taskFormPage
)

func CreatePages() {
	pages = []tea.Model{
		NewHome(),
		NewTaskForm(),
	}
}

func GetHomePage() tea.Model {
	return pages[homePage]
}

func GetTaskForm() tea.Model {
	return pages[taskFormPage]
}
