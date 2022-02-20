package config

import (
	"github.com/goal-web/contracts"
)

func Service(env, path string, config map[string]contracts.ConfigProvider) contracts.ServiceProvider {
	return &ServiceProvider{
		app:             nil,
		Env:             env,
		Paths:           []string{path},
		Sep:             "=",
		ConfigProviders: config,
	}
}

type ServiceProvider struct {
	app             contracts.Application
	Env             string
	Paths           []string
	Sep             string
	ConfigProviders map[string]contracts.ConfigProvider
}

func (this *ServiceProvider) Stop() {

}

func (this *ServiceProvider) Start() error {
	return nil
}

func (this *ServiceProvider) Register(application contracts.Application) {
	this.app = application

	application.Singleton("env", func() contracts.Env {
		return NewEnv(this.Paths, this.Sep)
	})

	application.Singleton("config", func(env contracts.Env) contracts.Config {
		configInstance := NewConfig(env, this.ConfigProviders)

		for key, provider := range this.ConfigProviders {
			configInstance.Set(key, provider(env))
		}
		return configInstance
	})
}
