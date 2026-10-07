package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"task-tracker/interal/domain"
)

type fileData struct {
	NextId domain.TaskId `json:"next_id"`
	Tasks  []domain.Task `json:"tasks"`
}

type TaskRepo struct {
	path string
}

func NewTaskRepo(path string) *TaskRepo {
	return &TaskRepo{
		path: path,
	}
}

func (r *TaskRepo) load() (*fileData, error) {
	data, err := os.ReadFile(r.path)
	if errors.Is(err, os.ErrNotExist) {
		return &fileData{NextId: 1}, nil
	}
	if err != nil {
		return nil, err
	}

	var fileData fileData
	err = json.Unmarshal(data, &fileData)
	if err != nil {
		return nil, err
	}
	return &fileData, nil
}

func (r *TaskRepo) AddTask(task domain.Task) error {
	fileData, err := r.load()
	if err != nil {
		return fmt.Errorf("search error in the file: %w", err)
	}

	fileData.Tasks = append(fileData.Tasks, task)
	fileData.NextId++

	data, err := json.Marshal(fileData)
	if err != nil {
		return domain.ErrIncurrectData
	}
	err = os.WriteFile(r.path, data, 0o644)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepo) GetNextId() (domain.TaskId, error) {
	fieldData, err := r.load()
	if err != nil {
		return -1, fmt.Errorf("search error in the file: %w", err)
	}
	return fieldData.NextId, nil
}
