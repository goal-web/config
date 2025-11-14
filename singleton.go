package config

import (
	"sync"

	"github.com/goal-web/application"
	"github.com/goal-web/contracts"
)

var singleton contracts.Config
var once sync.Once

// Default 返回配置单例。仅在首次调用时从默认应用中获取并缓存。
func Default() contracts.Config {
    once.Do(func() {
        singleton = application.Get("config").(contracts.Config)
    })

    return singleton
}

// Get 按点号路径获取配置值（如 "app.debug"）。
func Get(key string) any {
    return Default().Get(key)
}

// Set 设置或覆盖一个配置键的值，线程安全。
func Set(key string, value any) {
    Default().Set(key, value)
}

// Reload 重新计算并加载所有已注册的配置提供器（文件/URL/环境）。
func Reload() {
    Default().Reload()
}

// Unset 删除指定配置键。
func Unset(key string) {
    Default().Unset(key)
}
