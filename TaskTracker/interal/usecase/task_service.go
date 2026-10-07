package usecase

import (
	"errors"
	"task-tracker/interal/domain"
	"time"
)

var (
	ErrDeadlineInThePast = errors.New("set deadline in the past")
)

type Repo interface {
	GetNextId() (domain.TaskId, error)
	AddTask(task domain.Task) error
}

type TaskService struct {
	repo Repo
}

func NewTaskService(repo Repo) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) AddNewTask(title string, deadline time.Time) error {
	if deadline.Before(time.Now()) {
		return ErrDeadlineInThePast
	}

	id, err := s.repo.GetNextId()
	if err != nil {
		return err
	}
	task := domain.NewTask(id, title, deadline)
	err = s.repo.AddTask(*task)
	if err != nil {
		return err
	}
	return nil
}
