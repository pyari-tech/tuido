package tasksets

import "time"

type Status int8

const (
	todo Status = iota
	wip
	done
	blocked
	abandon
	archive
)

/*
Task would implement `list` with Title, Description, FilterValue
*/
type Task struct {
	status      Status
	title       string
	description string
	created     time.Time
	updated     time.Time
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
