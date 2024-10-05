package config

import (
	"github.com/goal-web/contracts"
	"sync"
)

func New(env contracts.Env, providers map[string]contracts.ConfigProvider) contracts.Config {
	return &config{
		writeMutex: sync.RWMutex{},
		providers:  providers,
		Env:        env,
		fields:     make(contracts.Fields),
	}
}

func WithFields(fields contracts.Fields) contracts.Config {
	return &config{fields: fields}
}

type config struct {
	writeMutex sync.RWMutex
	fields     contracts.Fields
	providers  map[string]contracts.ConfigProvider
	contracts.Env
}

func (config *config) ToFields() contracts.Fields {
	return config.fields
}

func (config *config) Reload() {
	for name, provider := range config.providers {
		config.Set(name, provider(config.Env))
	}
}

func (config *config) Set(key string, value any) {
	config.writeMutex.Lock()
	config.fields[key] = value
	config.writeMutex.Unlock()
}

func (config *config) Get(key string) any {
	config.writeMutex.RLock()
	defer config.writeMutex.RUnlock()

	if field, existsField := config.fields[key]; existsField {
		return field
	}

	if config.Env == nil {
		return nil
	}
	return config.Env.Get(key)
}

func (config *config) Unset(key string) {
	config.writeMutex.Lock()
	delete(config.fields, key)
	config.writeMutex.Unlock()
}
