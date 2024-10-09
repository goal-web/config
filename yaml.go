package config

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports"
	"github.com/goal-web/supports/utils"
	"gopkg.in/yaml.v3"
)

type yamlEnv struct {
	supports.BaseFields
	providers []EnvProvider
	fields    contracts.Fields
}

func NewYaml(providers ...EnvProvider) contracts.Env {
	provider := &yamlEnv{
		BaseFields: supports.BaseFields{OptionalGetter: osEnvGetter},
		providers:  providers,
	}

	provider.BaseFields.Provider = provider
	return provider
}
func (env *yamlEnv) ToFields() contracts.Fields {
	if env.fields == nil {
		env.fields = env.Load()
	}

	return env.fields
}

func (env *yamlEnv) Load() contracts.Fields {
	var envs = make(contracts.Fields)
	for _, provider := range env.providers {

		var data = make(map[any]any)
		err := yaml.Unmarshal(provider(), &data)
		if err != nil {
			log.Error("yamlEnv.load: " + err.Error())
			continue
		}
		fields, err := utils.ToFields(data)
		if err != nil {
			log.Error("yamlEnv.load: " + err.Error())
			continue
		}
		utils.Flatten(envs, fields, ".")
	}
	return envs
}
