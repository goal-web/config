package config

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
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
	return &config{
		fields: fields,
	}
}

type config struct {
	writeMutex sync.RWMutex
	fields     contracts.Fields
	providers  map[string]contracts.ConfigProvider
	contracts.Env
}

func (config *config) Fields() contracts.Fields {
	return config.fields
}

func (config *config) Load(provider contracts.FieldsProvider) {
	utils.MergeFields(config.fields, provider.Fields())
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

	return nil
}

func (config *config) Unset(key string) {
	delete(config.fields, key)
}
