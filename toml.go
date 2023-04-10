package config

import (
	"github.com/BurntSushi/toml"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports"
	"github.com/goal-web/supports/logs"
	"github.com/goal-web/supports/utils"
	"io"
	"net/http"
	"os"
)

type EnvProvider func() []byte

type tomlEnv[T any] struct {
	supports.BaseFields
	providers []EnvProvider
}

func NewToml(providers ...EnvProvider) contracts.Env {
	provider := &tomlEnv[any]{
		BaseFields: supports.BaseFields{Getter: func(key string) any {
			return os.Getenv(key)
		}},
		providers: providers,
	}

	provider.BaseFields.FieldsProvider = provider
	return provider
}

func File(path string) EnvProvider {
	return func() []byte {
		tmpBytes, err := os.ReadFile(path)
		if err != nil {
			log.Debug("File: " + err.Error())
		}
		return tmpBytes
	}
}

func Url(url string) EnvProvider {
	return func() []byte {
		res, err := http.Get(url)
		if err != nil {
			log.Debug("File: " + err.Error())
			return nil
		}
		tmpBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Debug("File: " + err.Error())
		}
		return tmpBytes
	}
}

var log = logs.Default()

func (env *tomlEnv[T]) Load() contracts.Fields {
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
