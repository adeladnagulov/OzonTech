package usecase

import (
	"fmt"
	"task-tracker/interal/domain"
	"time"
)

type Repo interface {
	GetNextId() (domain.TaskId, error)
	AddTask(task domain.Task) error
	GetTasksList() ([]domain.Task, error)
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
		return domain.ErrDeadlineInThePast
	}

	id, err := s.repo.GetNextId()
	if err != nil {
		return fmt.Errorf("failed to obtain nextId: %w", err)
	}
	task := domain.NewTask(id, title, deadline)
	err = s.repo.AddTask(task)
	if err != nil {
		return fmt.Errorf("failed to added new task: %w", err)
	}
	return nil
}

func (s *TaskService) GetAllTasks() ([]domain.Task, error) {
	tasks, err := s.repo.GetTasksList()
	if err != nil {
		return nil, fmt.Errorf("failed receiving tasks: %w", err)
	}
	return tasks, nil
}
