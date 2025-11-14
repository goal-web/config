package config

import (
	"os"

	"github.com/goal-web/contracts"
	"github.com/goal-web/supports"
	"github.com/joho/godotenv"
)

type dotEnv struct {
    supports.BaseFields
    providers []EnvProvider
    fields    contracts.Fields
}

// NewDotEnv 创建基于 DotEnv（键值对）数据源的环境读取器。
// 支持从本地文件或远程地址加载，OS 环境变量优先级最高。
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
		for key, value := range strFields {
			key = ToEnvKey(key)
			if _, exists := os.LookupEnv(key); exists {
				continue
			}
			err = os.Setenv(key, value)
			if err != nil {
				log.Error("dotEnv.load: " + err.Error())
				continue
			}
		}
	}
	return envs
}
