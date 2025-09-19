package registry

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type FileRepository struct {
	baseDir string
}

func NewFileRepository(baseDir string) *FileRepository {
	err := os.MkdirAll(baseDir, 0755)

	if err != nil {
		panic(err)
	}

	return &FileRepository{baseDir: baseDir}
}

func (r *FileRepository) Save(app AppState) error {
	data, _ := json.MarshalIndent(app, "", "  ")
	return os.WriteFile(filepath.Join(r.baseDir, app.Name+".json"), data, 0644)
}

func (r *FileRepository) Find(name string) (*AppState, error) {
	data, err := os.ReadFile(filepath.Join(r.baseDir, name+".json"))
	if err != nil {
		return nil, err
	}
	var app AppState
	if err := json.Unmarshal(data, &app); err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *FileRepository) FindAll() ([]AppState, error) {
	entries, err := os.ReadDir(r.baseDir)
	if err != nil {
		return nil, err
	}
	var apps []AppState
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, _ := os.ReadFile(filepath.Join(r.baseDir, e.Name()))
		var app AppState
		if err := json.Unmarshal(data, &app); err == nil {
			apps = append(apps, app)
		}
	}
	return apps, nil
}

func (r *FileRepository) Delete(name string) error {
	return os.Remove(filepath.Join(r.baseDir, name+".json"))
}
