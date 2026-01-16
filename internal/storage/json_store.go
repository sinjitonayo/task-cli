package storage

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/sinjito/task-cli-go/internal/model"
)

const DefaultFileName = "tasks.json"

type JSONStore struct {
	FilePath string
}

func NewJSONStore(filePath string) *JSONStore {
	if filePath == "" {
		filePath = DefaultFileName
	}
	return &JSONStore{FilePath: filePath}
}

// if the files does not exist, it will create an empty one
func (s *JSONStore) LoadTasks() ([]model.Task, error) {
	// create file if not exist
	if _, err := os.Stat(s.FilePath); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(s.FilePath, []byte("[]"), 0644); err != nil {
			return nil, err
		}
		return []model.Task{}, nil
	}

	data, err := os.ReadFile(s.FilePath)
	if err != nil {
		return nil, err
	}

	// treat as empty slice if file is empty
	if len(data) == 0 {
		return []model.Task{}, nil
	}

	var tasks []model.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// overwrites the entire file with the provided tasks
func (s *JSONStore) SaveTasks(tasks []model.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.FilePath, data, 0644)
}
