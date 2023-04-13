package config

import (
	"github.com/BurntSushi/toml"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports"
	"github.com/goal-web/supports/utils"
	"os"
)

type tomlEnv struct {
	supports.BaseFields
	providers []EnvProvider
	fields    contracts.Fields
}

func NewToml(providers ...EnvProvider) contracts.Env {
	provider := &tomlEnv{
		BaseFields: supports.BaseFields{OptionalGetter: func(key string, defaultValue any) any {
			if value, ok := os.LookupEnv(key); ok {
				return value
			}
			return defaultValue
		}},
		providers: providers,
	}

	provider.BaseFields.FieldsProvider = provider
	return provider
}
func (env *tomlEnv) Fields() contracts.Fields {
	if env.fields == nil {
		env.fields = env.Load()
	}

	return env.fields
}

func (env *tomlEnv) Load() contracts.Fields {
	var envs = make(contracts.Fields)
	for _, provider := range env.providers {

		var fields contracts.Fields
		err := toml.Unmarshal(provider(), &fields)
		if err != nil {
			log.Error("tomlEnv.load: " + err.Error())
			continue
		}
		utils.Flatten(envs, fields, ".")
	}
	return envs
}
