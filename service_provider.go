package config

import (
	"github.com/goal-web/contracts"
)

type serviceProvider struct {
	app             contracts.Application
	Env             string
	Paths           []string
	Sep             string
	ConfigProviders map[string]contracts.ConfigProvider
}

func NewService(env, path string, config map[string]contracts.ConfigProvider) contracts.ServiceProvider {
	return &serviceProvider{
		app:             nil,
		Env:             env,
		Paths:           []string{path},
		Sep:             "=",
		ConfigProviders: config,
	}
}

func (provider *serviceProvider) Stop() {

}

func (provider *serviceProvider) Start() error {
	return nil
}

func (provider *serviceProvider) Register(application contracts.Application) {
	provider.app = application

	application.Singleton("env", func() contracts.Env {
		return NewEnv(provider.Paths, provider.Sep)
	})

	application.Singleton("config", func(env contracts.Env) contracts.Config {
		configInstance := New(env, provider.ConfigProviders)

		for key, provider := range provider.ConfigProviders {
			configInstance.Set(key, provider(env))
		}
		return configInstance
	})
}
