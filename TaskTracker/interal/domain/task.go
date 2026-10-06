package domain

import "time"

type Task struct {
	Id       TaskId    `json:"id"`
	Title    string    `json:"title"`
	Status   Status    `json:"status"`
	Deadline time.Time `json:"deadline"`
}

func NewTask(id TaskId, title string, deadline time.Time) *Task {
	return &Task{
		Id:       id,
		Title:    title,
		Status:   PLANNED,
		Deadline: deadline,
	}
}
