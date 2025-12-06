package loader

import (
	"errors"
	"os"

	"github.com/Sn0wo2/Fcload"
	"gopkg.in/yaml.v3"
)

type YAML struct{}

func NewYAML() *YAML {
	return &YAML{}
}

func (l *YAML) Load(cfg any, fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, cfg)
}

func (l *YAML) Save(cfg any, fileName string) error {
	var node yaml.Node

	data, err := os.ReadFile(fileName)
	if errors.Is(err, os.ErrNotExist) {
		newData, marshalErr := yaml.Marshal(cfg)
		if marshalErr != nil {
			return marshalErr
		}
		return os.WriteFile(fileName, newData, 0644)
	} else if err != nil {
		return err
	}

	if err = yaml.Unmarshal(data, &node); err != nil {
		return err
	}

	Fcload.MergeYamlNode(&node, &cfg)

	newData, err := yaml.Marshal(&node)
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, newData, 0644)
}

func (l *YAML) GetTag() string {
	return "yaml"
}

func (l *YAML) GetAllowFileExtensions() []string {
	return []string{l.GetTag(), "yml"}
}
