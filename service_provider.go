package config

import (
	"github.com/goal-web/contracts"
)

type serviceProvider struct {
	app             contracts.Application
	Env             contracts.Env
	ConfigProviders map[string]contracts.ConfigProvider
}

func NewService(env contracts.Env, config map[string]contracts.ConfigProvider) contracts.ServiceProvider {
	return &serviceProvider{
		app:             nil,
		Env:             env,
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
		return provider.Env
	})

	application.Singleton("config", func(env contracts.Env) contracts.Config {
		configInstance := New(env, provider.ConfigProviders)

		for key, provider := range provider.ConfigProviders {
			configInstance.Set(key, provider(env))
		}
		return configInstance
	})
}
