package config

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports"
	"github.com/goal-web/supports/utils"
	"github.com/joho/godotenv"
	"os"
)

type dotEnvProvider[T any] struct {
	supports.BaseFields
	providers []EnvProvider
}

func NewDotEnv(providers ...EnvProvider) contracts.Env {
	provider := &dotEnvProvider[any]{
		BaseFields: supports.BaseFields{Getter: func(key string) any {
			return os.Getenv(key)
		}},
		providers: providers,
	}

	provider.BaseFields.FieldsProvider = provider
	return provider
}

func (env *dotEnvProvider[T]) Load() contracts.Fields {
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
