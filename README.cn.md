# Goal-Config

[![Go Reference](https://pkg.go.dev/badge/github.com/goal-web/config.svg)](https://pkg.go.dev/github.com/goal-web/config)
[![Go Report Card](https://goreportcard.com/badge/github.com/goal-web/config)](https://goreportcard.com/report/github.com/goal-web/config)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](../goal/LICENSE)
![GitHub Stars](https://img.shields.io/github/stars/goal-web/config?style=social)
![Release](https://img.shields.io/github/v/release/goal-web/config?include_prereleases)
![Go Version](https://img.shields.io/badge/go-%3E=%201.25.0-00ADD8?logo=go)

[Docs](https://pkg.go.dev/github.com/goal-web/config) · [Issues](https://github.com/goal-web/config/issues) · [Releases](https://github.com/goal-web/config/releases) · [CLI 加密](#cli-加密命令)

Goal 配置组件，支持从文件（TOML、YAML、DotEnv）和环境变量中加载配置，并自动合并优先级。

## 亮点

- 多源合并：本地文件与远程 URL 可组合，后者覆盖前者同名键。
- 环境优先：环境变量优先级最高，安全且可注入。
- 类型安全：提供类型化与 Optional 读取，避免隐式转换问题。
- 线程安全：读写分离锁，运行期可安全 `Set/Reload/Unset`。
- 可插拔：通过 `contracts.Env` 轻松扩展自定义源与预处理。

## 兼容性

- Go `>= 1.25.0`
- 模块路径：`github.com/goal-web/config`

## 目录

- 安装
- 快速开始
- 支持的文件格式
- 配置优先级
- 键命名规则
- 示例配置与环境变量示例
- 读取配置示例
- 可选配置获取方法
- 单例用法
- 组合多个数据源
- 环境变量覆盖示例
- 错误处理与诊断
- CLI 加密命令
- 注意事项

## 安装

```shell
go get github.com/goal-web/config
```

## 快速开始

### 基本用法

```go
package main

import (
    "github.com/goal-web/application"
    "github.com/goal-web/config"
    "github.com/goal-web/contracts"
    "fmt"
)

func main() {
    app := application.Default()

    // 注册配置服务
    app.RegisterService(
        config.NewService(
            config.NewToml(config.File("env.toml")), // 加载 TOML 文件
            map[string]contracts.ConfigProvider{},  // 自定义配置提供器（可选）
        ),
    )

    // 使用配置
    app.Call(func(conf contracts.Config) {
        debug := conf.GetBool("app.debug")
        if debug {
            fmt.Println("Debug mode enabled")
        }
    })
}
```

### 支持的文件格式

- TOML：使用 `config.NewToml(...)`
- YAML：使用 `config.NewYaml(...)`
- DotEnv（`.env`/键值对）：使用 `config.NewDotEnv(...)`

以上三种均可通过以下两类提供器读取源数据：
- `config.File(path)`：从本地文件加载
- `config.Url(url)`：从远程地址加载

示例：

```go
// 使用 YAML 文件
app.RegisterService(
    config.NewService(
        config.NewYaml(config.File("env.yaml")),
        map[string]contracts.ConfigProvider{},
    ),
)

// 使用 DotEnv 文件
app.RegisterService(
    config.NewService(
        config.NewDotEnv(config.File("config.env")),
        map[string]contracts.ConfigProvider{},
    ),
)

// 从 URL 加载（示例：TOML）
app.RegisterService(
    config.NewService(
        config.NewToml(config.Url("https://example.com/env.toml")),
        map[string]contracts.ConfigProvider{},
    ),
)
```

### 配置优先级

配置的加载优先级如下：
1. **环境变量**：最高优先级。
2. **配置文件**：次优先级（如 `env.toml`、`env.yaml`、`.env`）。
3. **默认值**：最低优先级（通过代码设置）。

#### 环境变量命名规则

- 将配置键中的 `.` 替换为 `_`。
- 转换为大写字母。

例如：
- 配置键 `app.debug` 对应的环境变量为 `APP_DEBUG`。
- 配置键 `database.host` 对应的环境变量为 `DATABASE_HOST`。

### 示例配置

#### `env.toml` 文件示例（完整示例）

```toml
[app]
name = "goal"
key = "dQcxsKvBZKNfWivwnhKlDwvseguknBZPEiiDRQlIatjKLLpbzK"
env = "local"
debug = true

[db.pgsql]
host = "localhost"
port = 55433
database = "postgres"
username = "postgres"
password = 123456
```

#### `env.yaml` 文件示例（完整示例）

```yaml
app:
  name: "goal"
  key: "dQcxsKvBZKNfWivwnhKlDwvseguknBZPEiiDRQlIatjKLLpbzK"
  env: "local"
  debug: true
db:
  pgsql:
    host: "localhost"
    port: 55433
    database: "postgres"
    username: "postgres"
    password: 123456
```

#### `.env`/`config.env` 文件示例（完整示例，键值对）

```env
# 注释
APP_KEY=dQcxsKvBZKNfWivwnhKlDwvseguknBZPEiiDRQlIatjKLLpbzK
APP_NAME=goal
APP_ENV=local
APP_DEBUG=true

HTTP_HOST=0.0.0.0
HTTP_PORT=8008

SQLITE_DATABASE=/Users/qbhy/project/go/goal-web/goal/example/database/db.sqlite

QUEUE_CONNECTION=nsq
QUEUE_KAFKA_BROKERS=localhost:9092
QUEUE_NSQ_ADDRESS=localhost:49162

DB_CONNECTION=sqlite
DB_HOST=localhost
DB_PORT=3308
DB_DATABASE=goal
DB_USERNAME=root
DB_PASSWORD=123456

DB_PGSQL_HOST=localhost
DB_PGSQL_PORT=55433
DB_PGSQL_DATABASE=postgres
DB_PGSQL_USERNAME=postgres
DB_PGSQL_PASSWORD=123456

REDIS_HOST=hsy
REDIS_PORT=6379
#REDIS_PASSWORD=123456

REDIS_CACHE_HOST=hsy
REDIS_CACHE_PORT=6379
REDIS_CACHE_DB=1

# 缓存配置
CACHE_DRIVER=redis
CACHE_CONNECTION=cache
CACHE_PREFIX=redis_

# 哈希配置
HASHING_DRIVER=bcrypt
HASHING_COST=14
HASHING_SALT=goal
# 自定义哈希
HASHING_HASHERS_MD5_DRIVER=md5
HASHING_HASHERS_MD5_SALT=goal

# 文件系统配置
FILESYSTEM_DRIVER=local
FILESYSTEM_ROOT=/Users/qbhy/project/go/goal/
FILESYSTEM_PERM=0777

QINIU_PRIVATE=false
QINIU_BUCKET=aa
QINIU_DOMAIN=https://xxx.xxx.com
QINIU_ACCESS_KEY=
QINIU_SECRET_KEY=


# session 配置
SESSION_ID=goal
SESSION_NAME=goal_session:
```

#### 环境变量示例

```shell
export APP_DEBUG=true
export DATABASE_HOST=mysql.example.com
```

### 读取配置示例

```go
app.Call(func(conf contracts.Config) {
    // 布尔值
    debug := conf.GetBool("app.debug")

    // 基本类型
    env := conf.GetString("app.env")
    port := conf.GetInt("db.pgsql.port")

    // 嵌套键（点号展开）
    host := conf.GetString("db.pgsql.host")

    fmt.Printf("env=%s, debug=%t, db=%s:%d\n", env, debug, host, port)
})
```

### 可选配置获取方法

当某个配置键不存在时，可以使用可选获取方法提供默认值。这些方法的优先级与普通获取一致：环境变量 > 配置文件 > 默认值。

常用方法包括：
- `StringOptional(key, default)`
- `IntOptional(key, default)` / `Int64Optional` / `UIntOptional` 等
- `BoolOptional(key, default)`
- `FloatOptional(key, default)` / `Float64Optional`

用法示例：

```go
package main

import (
    "fmt"
    "github.com/goal-web/application"
    "github.com/goal-web/config"
    "github.com/goal-web/contracts"
)

func main() {
    app := application.Default()
    app.RegisterService(
        config.NewService(
            config.NewToml(config.File("env.toml")),
            map[string]contracts.ConfigProvider{},
        ),
    )

    app.Call(func(conf contracts.Config) {
        // 如果 app.name 不存在，则返回默认值 "goal"
        name := conf.StringOptional("app.name", "goal")

        // 如果 hashing.cost 不存在，则返回默认值 10
        cost := conf.IntOptional("hashing.cost", 10)

        // 可选布尔值
        debug := conf.BoolOptional("app.debug", false)

        // 包级函数也可用（读取默认 config 单例）
        salt := config.StringOptional("hashing.salt", "default_salt")

        fmt.Printf("name=%s, cost=%d, debug=%t, salt=%s\n", name, cost, debug, salt)
    })
}
```

### 单例（Singleton）用法

当应用使用 `application.Default()` 并注册了配置服务后，可以直接使用 `config` 包的单例方法（见 `config/singleton.go`），无需在代码中显式传递 `contracts.Config`。

可用的单例方法包括：
- `config.Default()`：获取配置单例
- `config.Get(key)` / `config.Set(key, value)` / `config.Unset(key)`
- `config.Reload()`：根据注册的 `ConfigProvider` 重新加载配置
- 类型化辅助：`config.GetString`、`config.GetBool`、`config.IntOptional`、`config.StringOptional` 等（见 `helper.go`）

示例：

```go
package main

import (
    "fmt"
    "github.com/goal-web/application"
    "github.com/goal-web/config"
    "github.com/goal-web/contracts"
)

func main() {
    app := application.Default()
    app.RegisterService(
        config.NewService(
            config.NewToml(config.File("env.toml")),
            map[string]contracts.ConfigProvider{},
        ),
    )

    // 假设此时应用作为默认应用运行（例如框架启动流程中已设置）
    // 可直接使用包级单例方法
    env := config.GetString("app.env")
    debug := config.GetBool("app.debug")

    // 写入/覆盖某个键
    config.Set("feature.toggle", true)

    // 重新加载已注册的配置提供器（文件/URL 等）
    config.Reload()

    // 删除一个键
    config.Unset("db.pgsql.password")

    fmt.Printf("env=%s, debug=%t\n", env, debug)
}
```

### 组合多个数据源

```go
app.RegisterService(
    config.NewService(
        // 依次加载并合并（后者可覆盖前者同名键）
        config.NewToml(
            config.File("env.toml"),
            config.Url("https://example.com/override.toml"),
        ),
        map[string]contracts.ConfigProvider{},
    ),
)
```

### 错误处理与诊断

- 文件读取失败：本地 `File(path)` 读取错误会输出调试日志；返回空字节，后续解析会跳过该源。
- 远程 URL 失败：`Url(url)` 网络或读取失败会记录调试日志并返回 `nil`。
- 解析失败：TOML/YAML/DotEnv 解析异常会记录错误日志；该源不影响其他源的加载。
- 键覆盖冲突：后加载的数据源可覆盖先前同名键；最终值以环境变量为最高优先级。
- 类型不匹配：请使用类型化读取（如 `GetBool`、`GetInt`），避免隐式转换导致运行时错误。

### CLI 加密命令

配置文件可通过 CLI 进行加密/解密，避免明文环境泄露。

命令：

```shell
# 加密 .env 到 .env.encrypted，默认 AES 驱动，自动生成密钥
goal env encrypt --in .env --out .env.encrypted

# 使用指定密钥与驱动加密
goal env encrypt --driver AES --key "your-32-characters-key" --in .env --out .env.encrypted

# 解密 .env.encrypted 到 .env（如果文件存在，需 --force 覆盖）
goal env decrypt --in .env.encrypted --out .env --key "your-32-characters-key" --force
```

说明：
- 不指定 `--key` 时将自动生成 32 字符密钥，并在日志中输出。
- 通过 `--driver` 指定加密驱动（默认 `AES`）。
- 解密时如目标文件已存在，需添加 `--force` 才会覆盖。

### 环境变量覆盖示例

```go
// 将 app.env 的值覆盖为 testing（优先级最高）
os.Setenv("APP_ENV", "testing")
env := config.NewToml(config.File("env.toml"))
conf := config.New(env, nil)
fmt.Println(conf.GetString("app.env")) // 输出: testing
```

## Star History

<a href="https://star-history.com/#goal-web/config&Date"><img src="https://api.star-history.com/svg?repos=goal-web/config&type=Date" alt="Star History Chart"/></a>

![Stargazers over time](https://starchart.cc/goal-web/config.svg)

## Contributing

- 欢迎提交 Issue 与 Pull Request，建议遵循规范：
  - 保持 API 向后兼容（语义化改动通过新增方法或选项）。
  - 更新相关文档与示例；涉及行为变更请在 `注意事项/FAQ` 添加说明。
  - 提交前进行本地构建与基本测试校验。
- 讨论或建议：请使用 Issues 或 Discussions（如仓库开启）。

## Roadmap

- 子模块文档矩阵：为 `goal/config` 与 `goal-cli/config` 的各模块补充 README 与示例。
- 示例工程与冒烟测试：提供最小应用与覆盖组合源的验证。
- 英文文档：提供 `README.en.md`，方便国际开发者使用。

## Changelog

- 版本与变更记录请见 Releases：<https://github.com/goal-web/config/releases>

## Security

- 切勿将密钥与敏感信息提交到仓库；建议使用环境变量或 CLI 加密。
- 环境变量优先级最高，适用于在容器/CI 中注入配置。
- 对于生产环境的配置文件，建议加密并在部署时解密：见“CLI 加密命令”。

## FAQ

- 文件读取失败：`File(path)` 或 `Url(url)` 失败会记录日志并跳过该源，其它源不受影响。
- 值覆盖顺序：后加载的数据源覆盖前者；最终以环境变量为最高优先级。
- 类型不匹配：使用 `GetBool`、`GetInt` 等类型化方法，避免隐式转换导致意外行为。
### 注意事项

1. **文件格式支持**：内置支持 TOML、YAML 和 DotEnv。你可以通过实现 `contracts.Env` 接口扩展其他来源或处理流程（例如自定义拉取、预处理等）。
2. **动态加载**：应用启动时初始化配置；若运行期修改文件或远程源，可显式调用 `config.Reload()` 重新计算并生效（环境变量始终优先）。
3. **类型安全**：使用 `GetBool`、`GetString` 等方法确保类型安全。
