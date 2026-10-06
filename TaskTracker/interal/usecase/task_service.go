package usecase

import (
	"task-tracker/interal/domain"
	"task-tracker/interal/storage"
	"time"
)

type TaskService struct {
	repo *storage.TaskRepo
}

func NewTaskService(repo *storage.TaskRepo) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) AddNewTask(title string, deadline time.Time) error {
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
