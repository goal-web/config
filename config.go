package config

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
	"strings"
	"sync"
)

func NewConfig(env contracts.Env, providers map[string]contracts.ConfigProvider) contracts.Config {
	return &config{
		writeMutex: sync.RWMutex{},
		providers:  providers,
		Env:        env,
		fields:     make(contracts.Fields),
		configs:    make(map[string]contracts.Config, 0),
	}
}

func WithFields(fields contracts.Fields) contracts.Config {
	return &config{
		fields:  fields,
		configs: make(map[string]contracts.Config, 0),
	}
}

type config struct {
	writeMutex sync.RWMutex
	fields     contracts.Fields
	configs    map[string]contracts.Config
	providers  map[string]contracts.ConfigProvider
	contracts.Env
}

func (this *config) Fields() contracts.Fields {
	return this.fields
}

func (this *config) Load(provider contracts.FieldsProvider) {
	utils.MergeFields(this.fields, provider.Fields())
}

func (this *config) Reload() {
	for name, provider := range this.providers {
		this.Set(name, provider(this.Env))
	}
}

func (this *config) Merge(key string, config contracts.Config) {
	this.fields[key] = config.Fields()
	this.configs[key] = config
}

func (this *config) Set(key string, value interface{}) {
	this.writeMutex.Lock()
	this.fields[key] = value
	this.writeMutex.Unlock()
}

func (this *config) Get(key string, defaultValue ...interface{}) interface{} {
	this.writeMutex.RLock()
	defer this.writeMutex.RUnlock()

	// 环境变量优先级最高
	if this.Env != nil {
		if envValue := this.Env.GetString(key); envValue != "" {
			return envValue
		}
	}

	if field, existsField := this.fields[key]; existsField {
		return field
	}

	// 尝试获取 fields
	var (
		fields = contracts.Fields{}
		prefix = key + "."
	)

	for fieldKey, fieldValue := range this.fields {
		if strings.HasPrefix(fieldKey, prefix) {
			fields[strings.Replace(fieldKey, prefix, "", 1)] = fieldValue
		}
	}

	if len(fields) > 0 {
		return fields
	}

	var keys = strings.Split(key, ".")

	if len(keys) > 1 {
		if subConfig, existsSubConfig := this.configs[keys[0]]; existsSubConfig {
			return subConfig.Get(strings.Join(keys[1:], "."), defaultValue...)
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return nil
}

func (this *config) GetConfig(key string) contracts.Config {
	return this.configs[key]
}

func (this *config) GetFields(key string) contracts.Fields {
	if field, isTypeRight := this.Get(key).(contracts.Fields); isTypeRight {
		return field
	}

	return nil
}

func (this *config) GetString(key string) string {
	if field, isTypeRight := this.Get(key).(string); isTypeRight {
		return field
	}

	return ""
}

func (this *config) GetInt(key string) int {
	if field := this.Get(key); field != nil {
		value := utils.ConvertToInt(field, 0)
		if value != 0 { // 缓存转换结果
			this.Set(key, value)
		}
		return value
	}

	return 0
}
func (this *config) GetInt64(key string) int64 {
	if field := this.Get(key); field != nil {
		value := utils.ConvertToInt64(field, 0)
		if value != 0 { // 缓存转换结果
			this.Set(key, value)
		}
		return value
	}

	return 0
}

func (this *config) Unset(key string) {
	delete(this.fields, key)
	delete(this.configs, key)
}

func (this *config) GetFloat(key string) float32 {
	if field := this.Get(key); field != nil {
		value := utils.ConvertToFloat(field, 0)
		if value != 0 { // 缓存转换结果
			this.Set(key, value)
		}
		return value
	}

	return 0
}
func (this *config) GetFloat64(key string) float64 {
	if field := this.Get(key); field != nil {
		value := utils.ConvertToFloat64(field, 0)
		if value != 0 { // 缓存转换结果
			this.Set(key, value)
		}
		return value
	}

	return 0
}

func (this *config) GetBool(key string) bool {
	if field := this.Get(key); field != nil {
		result := utils.ConvertToBool(field, false)
		this.Set(key, result)
		return result
	}

	return false
}
