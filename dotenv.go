package config

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports"
	"github.com/goal-web/supports/utils"
	"github.com/joho/godotenv"
)

type dotEnv struct {
	supports.BaseFields
	providers []EnvProvider
	fields    contracts.Fields
}

func NewDotEnv(providers ...EnvProvider) contracts.Env {
	provider := &dotEnv{
		BaseFields: supports.BaseFields{OptionalGetter: osEnvGetter},
		providers:  providers,
	}

	provider.BaseFields.Provider = provider
	return provider
}

func (env *dotEnv) ToFields() contracts.Fields {
	if env.fields == nil {
		env.fields = env.Load()
	}

	return env.fields
}

func (env *dotEnv) Load() contracts.Fields {
	var envs = make(contracts.Fields)
	for _, provider := range env.providers {

		strFields, err := godotenv.UnmarshalBytes(provider())
		if err != nil {
			log.Error("tomlEnv.load: " + err.Error())
			continue
		}
		fields, _ := utils.ToFields(strFields)
		utils.Flatten(envs, fields, ".")
	}
	return envs
}
