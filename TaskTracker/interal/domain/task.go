package domain

import "time"

type Task struct {
	Id       TaskId    `json:"id"`
	Title    string    `json:"title"`
	Status   Status    `json:"status"`
	Deadline time.Time `json:"deadline"`
}
