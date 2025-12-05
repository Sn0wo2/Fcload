package Fcload

import "errors"

var ErrConfigNotFound = errors.New("config file not found")

var Path string

type Loader interface {
	Load(cfg any, fileName string) error
	Save(cfg any, fileName string) error
	GetTag() string
	GetAllowFileExtensions() []string
}
