package domain

import (
	"fmt"
	"time"
)

type Task struct {
	Id       TaskId    `json:"id"`
	Title    string    `json:"title"`
	Status   Status    `json:"status"`
	Deadline time.Time `json:"deadline"`
}

func NewTask(id TaskId, title string, deadline time.Time) Task {
	return Task{
		Id:       id,
		Title:    title,
		Status:   PLANNED,
		Deadline: deadline,
	}
}

func (t Task) String() string {
	deadlene := t.Deadline.Format("2006-01-02")
	return fmt.Sprintf("Id: %v\nTitle: %s\nStatus: %v\nDeadline: %v\n", t.Id, t.Title, t.Status, deadlene)
}
