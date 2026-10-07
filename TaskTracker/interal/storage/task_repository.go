package storage

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"task-tracker/interal/domain"
	"time"
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

func (r *TaskRepo) GetTasksList() ([]domain.Task, error) {
	filedata, err := r.load()
	if err != nil {
		return nil, fmt.Errorf("search error in the file: %w", err)
	}
	slices.SortFunc(filedata.Tasks, func(a, b domain.Task) int {
		if n := time.Time.Compare(a.Deadline, b.Deadline); n != 0 {
			return n
		}
		return cmp.Compare(a.Id, b.Id)
	})
	return filedata.Tasks, nil
}

func (r *TaskRepo) GetTask(id domain.TaskId) (domain.Task, error) {
	filedata, err := r.load()
	if err != nil {
		return domain.Task{}, fmt.Errorf("search error in the file: %w", err)
	}

	for _, v := range filedata.Tasks {
		if v.Id == id {
			return v, nil
		}
	}
	return domain.Task{}, domain.ErrTaskNotFound
}
