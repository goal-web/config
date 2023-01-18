package config

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
	"strings"
	"sync"
)

func New(env contracts.Env, providers map[string]contracts.ConfigProvider) contracts.Config {
	return &config{
		writeMutex: sync.RWMutex{},
		providers:  providers,
		Env:        env,
		fields:     make(contracts.Fields),
	}
}

func WithFields(fields contracts.Fields) contracts.Config {
	return &config{
		fields: fields,
	}
}

type config struct {
	writeMutex sync.RWMutex
	fields     contracts.Fields
	providers  map[string]contracts.ConfigProvider
	contracts.Env
}

func (config *config) Fields() contracts.Fields {
	return config.fields
}

func (config *config) Load(provider contracts.FieldsProvider) {
	utils.MergeFields(config.fields, provider.Fields())
}

func (config *config) Reload() {
	for name, provider := range config.providers {
		config.Set(name, provider(config.Env))
	}
}

func (config *config) Set(key string, value interface{}) {
	config.writeMutex.Lock()
	config.fields[key] = value
	config.writeMutex.Unlock()
}

func (config *config) Get(key string, defaultValue ...interface{}) interface{} {
	config.writeMutex.RLock()
	defer config.writeMutex.RUnlock()

	// 环境变量优先级最高
	if config.Env != nil {
		if envValue := config.Env.GetString(key); envValue != "" {
			return envValue
		}
	}

	if field, existsField := config.fields[key]; existsField {
		return field
	}

	// 尝试获取 fields
	var (
		fields = contracts.Fields{}
		prefix = key + "."
	)

	for fieldKey, fieldValue := range config.fields {
		if strings.HasPrefix(fieldKey, prefix) {
			fields[strings.Replace(fieldKey, prefix, "", 1)] = fieldValue
		}
	}

	if len(fields) > 0 {
		return fields
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return nil
}

func (config *config) GetFields(key string) contracts.Fields {
	if field, isTypeRight := config.Get(key).(contracts.Fields); isTypeRight {
		return field
	}

	return nil
}

func (config *config) GetString(key string) string {
	if field, isTypeRight := config.Get(key).(string); isTypeRight {
		return field
	}

	return ""
}

func (config *config) GetInt(key string) int {
	if field := config.Get(key); field != nil {
		value := utils.ConvertToInt(field, 0)
		if value != 0 { // 缓存转换结果
			config.Set(key, value)
		}
		return value
	}

	return 0
}
func (config *config) GetInt64(key string) int64 {
	if field := config.Get(key); field != nil {
		value := utils.ConvertToInt64(field, 0)
		if value != 0 { // 缓存转换结果
			config.Set(key, value)
		}
		return value
	}

	return 0
}

func (config *config) Unset(key string) {
	delete(config.fields, key)
}

func (config *config) GetFloat(key string) float32 {
	if field := config.Get(key); field != nil {
		value := utils.ConvertToFloat(field, 0)
		if value != 0 { // 缓存转换结果
			config.Set(key, value)
		}
		return value
	}

	return 0
}
func (config *config) GetFloat64(key string) float64 {
	if field := config.Get(key); field != nil {
		value := utils.ConvertToFloat64(field, 0)
		if value != 0 { // 缓存转换结果
			config.Set(key, value)
		}
		return value
	}

	return 0
}

func (config *config) GetBool(key string) bool {
	if field := config.Get(key); field != nil {
		result := utils.ConvertToBool(field, false)
		config.Set(key, result)
		return result
	}

	return false
}
