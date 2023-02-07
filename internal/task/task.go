package task

type Status int

const (
	Todo Status = iota
	InProgress
	Done
)

type Task struct {
	status      Status
	title       string
	description string
}

func New(status Status, title, description string) Task {
	return Task{
		status:      status,
		title:       title,
		description: description,
	}
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

func (t Task) Status() Status {
	return t.status
}

func (t *Task) Next() {
	if t.status == Done {
		t.status = Todo
	} else {
		t.status++
	}
}
