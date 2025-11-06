package config

import (
	"sync"

	"github.com/goal-web/application"
	"github.com/goal-web/contracts"
)

var singleton contracts.Config
var once sync.Once

func Default() contracts.Config {
	once.Do(func() {
		singleton = application.Get("config").(contracts.Config)
	})

	return singleton
}

func Get(key string) any {
	return Default().Get(key)
}

func Set(key string, value any) {
	Default().Set(key, value)
}

func Reload() {
	Default().Reload()
}

func Unset(key string) {
	Default().Unset(key)
}
