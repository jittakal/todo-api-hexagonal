package domain

type Task struct {
	ID    string
	Title string
	Done  bool
}

func NewTask(id string) *Task {
	return &Task{
		ID: id,
	}
}
