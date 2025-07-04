package tasksets

/*
tasksets::taskform
*/

import (
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TaskAction uint8

const (
	noTask TaskAction = iota
	createTask
	updateTask
	readTask
)

/* Struct TaskForm for taskFormPage's Model */

type TaskForm struct {
	status      Status
	title       textinput.Model
	description textarea.Model
	action      TaskAction
}

func NewTaskForm() *TaskForm {
	tf := &TaskForm{
		status:      todo, // any fresh task goes to ToDo first
		title:       textinput.New(),
		description: textarea.New(),
	}
	tf.title.Placeholder = "to do.."
	tf.title.PlaceholderStyle = lipgloss.NewStyle().Italic(true)
	tf.title.Focus()
	tf.title.Cursor.Blink = true
	return tf
}

func (tf TaskForm) Init() tea.Cmd {
	return nil
}

func (tf TaskForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+n":
			tf.action = createTask
			return tf, textinput.Blink
		case "ctrl+u":
			tf.action = updateTask
			return tf, textinput.Blink
		case "ctrl+r":
			tf.action = readTask
			tf.title.Blur()
			return tf, textinput.Blink
		case "enter":
			if tf.action == readTask {
				tf.action = noTask
				return pages[homePage], nil
			}
			if tf.title.Focused() {
				tf.title.Blur()
				tf.description.Focus()
				return tf, textarea.Blink
			} else if tf.description.Focused() {
				tf.description.Blur()
				tf.action = noTask
				pages[taskFormPage] = tf
				return pages[homePage], tf.AddTaskToHome
			}
		case "ctrl+k":
			tf.action = noTask
			return pages[homePage], nil
		case "ctrl+c":
			return tf, tea.Quit
		}
	}

	var cmd tea.Cmd
	if tf.title.Focused() {
		tf.title, cmd = tf.title.Update(msg)
	} else if tf.description.Focused() {
		tf.description, cmd = tf.description.Update(msg)
	}
	return tf, cmd
}

func (tf TaskForm) View() string {
	var tfTitle string
	switch tf.action {
	case createTask:
		tfTitle = "Create a new task (ctrl+k to skip)"
	case updateTask:
		tfTitle = "Update existing task (ctrl+k to skip)"
	case readTask:
		tfTitle = "Task Details (ctrl+k or ENTER for Task Board)"
	default:
		tf.action = noTask
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		tfTitle,
		tf.title.View(),
		tf.description.View(),
	)
}

func (tf TaskForm) AddTaskToHome() tea.Msg {
	return Task{
		status:      tf.status,
		title:       tf.title.Value(),
		description: tf.description.Value(),
		created:     time.Now(),
	}
}
