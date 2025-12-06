package loader

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type JSON struct{}

func NewJSON() *JSON {
	return &JSON{}
}

func (l *JSON) Load(cfg any, fileName string) error {
	data, err := os.ReadFile(fileName) //nolint:gosec
	if err != nil {
		return err
	}
	return json.Unmarshal(data, cfg)
}

func (l *JSON) Save(cfg any, fileName string) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(fileName), 0o750); err != nil {
		return err
	}

	return os.WriteFile(fileName, data, 0644)
}

func (l *JSON) GetAllowFileExtensions() []string {
	return []string{"json"}
}
