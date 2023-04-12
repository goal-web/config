package config

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports"
	"github.com/goal-web/supports/utils"
	"github.com/joho/godotenv"
	"os"
)

type dotEnv[T any] struct {
	supports.BaseFields
	providers []EnvProvider
	fields    contracts.Fields
}

func NewDotEnv(providers ...EnvProvider) contracts.Env {
	provider := &dotEnv[any]{
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

func (env *dotEnv[T]) Fields() contracts.Fields {
	if env.fields == nil {
		env.fields = env.Load()
	}

	return env.fields
}

func (env *dotEnv[T]) Load() contracts.Fields {
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
