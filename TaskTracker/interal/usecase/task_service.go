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
	GetTask(id domain.TaskId) (domain.Task, error)
	ChangeStatus(id domain.TaskId, newStatus domain.Status) error
	DeleteTask(id domain.TaskId) error
	IncNextId() error
}

type TaskService struct {
	repo Repo
}

func NewTaskService(repo Repo) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) AddNewTask(taskIn domain.TaskInpuut) error {
	title, deadline := *taskIn.Title, *taskIn.Deadline
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
	err = s.repo.IncNextId()
	if err != nil {
		return fmt.Errorf("failed increment id: %w", err)
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

func (s *TaskService) GetTask(id int64) (domain.Task, error) {
	task, err := s.repo.GetTask(domain.TaskId(id))
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed receiving task by id: %w", err)
	}
	return task, nil
}

func (s *TaskService) ChangeStatus(id int64, newStatus string) error {
	status, err := domain.ValidStatus(newStatus)
	if err != nil {
		return fmt.Errorf("failed receiving status: %w", err)
	}
	err = s.repo.ChangeStatus(domain.TaskId(id), status)
	if err != nil {
		return fmt.Errorf("failed update status: %w", err)
	}
	return nil
}

func (s *TaskService) DeleteTask(id int64) error {
	err := s.repo.DeleteTask(domain.TaskId(id))
	if err != nil {
		return fmt.Errorf("failed delete task: %w", err)
	}
	return nil
}

func (s *TaskService) EditTask(id int64, taskIn domain.TaskInpuut) error {
	task, err := s.GetTask(id)
	if err != nil {
		return fmt.Errorf("failed getting task by edit him: %w", err)
	}

	var title string
	var deadline time.Time
	if taskIn.Title != nil {
		title = *taskIn.Title
		if title == "" {
			return domain.ErrEmptyTitle
		}
		task.Title = title
	}
	if taskIn.Deadline != nil {
		deadline = *taskIn.Deadline
		if deadline.Before(time.Now()) {
			return domain.ErrDeadlineInThePast
		}
		task.Deadline = deadline
	}

	err = s.DeleteTask(id)
	if err != nil {
		return fmt.Errorf("failed delete task by edit him: %w", err)
	}

	err = s.repo.AddTask(task)
	if err != nil {
		return fmt.Errorf("failed return updated task: %w", err)
	}
	return nil
}
