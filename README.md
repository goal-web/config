# Goal-Config

[![Go Reference](https://pkg.go.dev/badge/github.com/goal-web/config.svg)](https://pkg.go.dev/github.com/goal-web/config)
[![Go Report Card](https://goreportcard.com/badge/github.com/goal-web/config)](https://goreportcard.com/report/github.com/goal-web/config)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](../goal/LICENSE)
![GitHub Stars](https://img.shields.io/github/stars/goal-web/config?style=social)
![Release](https://img.shields.io/github/v/release/goal-web/config?include_prereleases)
![Go Version](https://img.shields.io/badge/go-%3E=%201.25.0-00ADD8?logo=go)
![CI](https://img.shields.io/github/actions/workflow/status/goal-web/config/ci.yml?branch=master&label=CI)
![Lint](https://img.shields.io/github/actions/workflow/status/goal-web/config/lint.yml?branch=master&label=Lint)
![Commit Activity](https://img.shields.io/github/commit-activity/m/goal-web/config)

[Docs](https://pkg.go.dev/github.com/goal-web/config) · [Issues](https://github.com/goal-web/config/issues) · [Releases](https://github.com/goal-web/config/releases) · [CLI Encryption](#cli-encryption) · [中文文档](./README.cn.md)

Goal Config provides configuration loading from files (TOML, YAML, DotEnv) and environment variables, with automatic merging and clear precedence.

## Highlights

- Multi-source merge: local files and remote URLs can be composed; later sources override earlier ones.
- Environment-first: OS environment variables have the highest precedence.
- Type-safe: typed getters and optional variants avoid implicit conversions.
- Thread-safe: RW locks for reads/writes; runtime-safe `Set/Reload/Unset`.
- Pluggable: implement `contracts.Env` to extend sources and preprocessing.

## Compatibility

- Go `>= 1.25.0`
- Module path: `github.com/goal-web/config`

## Table of Contents

- Installation
- Quick Start
- Supported Formats
- Config Priority
- Key Naming Rules
- Example Configs & Env Examples
- Reading Config
- Optional Getters
- Singleton Usage
- Combining Sources
- Env Override Example
- Error Handling & Diagnostics
- CLI Encryption
- Notes

## Installation

```shell
go get github.com/goal-web/config
```

## Quick Start

### Basic usage

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

    // register config service
    app.RegisterService(
        config.NewService(
            config.NewToml(config.File("env.toml")),
            map[string]contracts.ConfigProvider{},
        ),
    )

    // use config
	app.Call(func(conf contracts.Config) {
		debug := conf.GetBool("app.debug")
		if debug {
			fmt.Println("Debug mode enabled")
		}
	})
}
```

### Supported Formats

- TOML: `config.NewToml(...)`
- YAML: `config.NewYaml(...)`
- DotEnv (`.env` key-value): `config.NewDotEnv(...)`

Sources for all formats:
- `config.File(path)`: load from local file
- `config.Url(url)`: load from remote URL

Examples:

```go
// YAML file
app.RegisterService(
    config.NewService(
        config.NewYaml(config.File("env.yaml")),
        map[string]contracts.ConfigProvider{},
    ),
)

// DotEnv file
app.RegisterService(
    config.NewService(
        config.NewDotEnv(config.File("config.env")),
        map[string]contracts.ConfigProvider{},
    ),
)

// Load from URL (TOML)
app.RegisterService(
    config.NewService(
        config.NewToml(config.Url("https://example.com/env.toml")),
        map[string]contracts.ConfigProvider{},
    ),
)
```

### Config Priority

Precedence:
1. Environment variables (highest)
2. Config files (e.g. `env.toml`, `env.yaml`, `.env`)
3. Defaults (set in code)

#### Key Naming Rules for Env

- Replace `.` with `_` in keys
- Uppercase letters

Examples:
- `app.debug` → `APP_DEBUG`
- `database.host` → `DATABASE_HOST`

### Example Configs

#### `env.toml`

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

#### `env.yaml`

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

#### `.env` / `config.env`

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

#### Env examples

```shell
export APP_DEBUG=true
export DATABASE_HOST=mysql.example.com
```

### Reading Config

```go
app.Call(func(conf contracts.Config) {
    // booleans
    debug := conf.GetBool("app.debug")

    // primitives
    env := conf.GetString("app.env")
    port := conf.GetInt("db.pgsql.port")

    // nested keys (dot notation)
    host := conf.GetString("db.pgsql.host")

    fmt.Printf("env=%s, debug=%t, db=%s:%d\n", env, debug, host, port)
})
```

### Optional Getters

When a key is missing, optional getters provide defaults. Same precedence: env > files > defaults.

Common methods:
- `StringOptional(key, default)`
- `IntOptional(key, default)` / `Int64Optional` / `UIntOptional` 等
- `BoolOptional(key, default)`
- `FloatOptional(key, default)` / `Float64Optional`

Example:

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
        // default if app.name missing
        name := conf.StringOptional("app.name", "goal")

        // default if hashing.cost missing
        cost := conf.IntOptional("hashing.cost", 10)

        // optional boolean
        debug := conf.BoolOptional("app.debug", false)

        // package-level helpers work with singleton
        salt := config.StringOptional("hashing.salt", "default_salt")

        fmt.Printf("name=%s, cost=%d, debug=%t, salt=%s\n", name, cost, debug, salt)
    })
}
```

### Singleton Usage

With `application.Default()` and config service registered, you can use package-level singleton helpers (see `config/singleton.go`) without passing `contracts.Config` explicitly.

Available helpers:
- `config.Default()`：获取配置单例
- `config.Get(key)` / `config.Set(key, value)` / `config.Unset(key)`
- `config.Reload()`：根据注册的 `ConfigProvider` 重新加载配置
- Typed helpers: `config.GetString`, `config.GetBool`, `config.IntOptional`, `config.StringOptional` (see `helper.go`)

Example:

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

    // use package-level singleton methods
    env := config.GetString("app.env")
    debug := config.GetBool("app.debug")

    // set/override
    config.Set("feature.toggle", true)

    // reload registered providers (files/URL)
    config.Reload()

    // unset a key
    config.Unset("db.pgsql.password")

    fmt.Printf("env=%s, debug=%t\n", env, debug)
}
```

### Combining Multiple Sources

```go
app.RegisterService(
    config.NewService(
        // merge in order (later overrides earlier)
        config.NewToml(
            config.File("env.toml"),
            config.Url("https://example.com/override.toml"),
        ),
        map[string]contracts.ConfigProvider{},
    ),
)
```

### Error Handling & Diagnostics

- File read errors: `File(path)` logs debug and returns empty bytes; parser skips this source.
- Remote URL errors: `Url(url)` logs debug and returns `nil`.
- Parse errors: TOML/YAML/DotEnv errors are logged; other sources continue.
- Key conflicts: later sources override earlier; env takes highest precedence.
- Type mismatch: use typed getters like `GetBool`, `GetInt`.

### CLI Encryption

Encrypt/decrypt config files via CLI to avoid plaintext secrets.

Commands:

```shell
# Encrypt .env to .env.encrypted (AES, auto-generate key)
goal env encrypt --in .env --out .env.encrypted

# Encrypt with specified key/driver
goal env encrypt --driver AES --key "your-32-characters-key" --in .env --out .env.encrypted

# Decrypt .env.encrypted to .env (use --force to overwrite)
goal env decrypt --in .env.encrypted --out .env --key "your-32-characters-key" --force
```

Notes:
- If `--key` missing, a 32-char key is generated and logged.
- Use `--driver` to select encryption driver (default `AES`).
- `--force` overwrites target on decrypt.

### Env Override Example

```go
// override app.env via OS env (highest precedence)
os.Setenv("APP_ENV", "testing")
env := config.NewToml(config.File("env.toml"))
conf := config.New(env, nil)
fmt.Println(conf.GetString("app.env")) // prints: testing
```

## Star History

<a href="https://star-history.com/#goal-web/config&Date"><img src="https://api.star-history.com/svg?repos=goal-web/config&type=Date" alt="Star History Chart"/></a>

![Stargazers over time](https://starchart.cc/goal-web/config.svg)

## CI Insights

- CI Status: ![CI](https://img.shields.io/github/actions/workflow/status/goal-web/config/ci.yml?branch=main&label=CI)
- Lint Status: ![Lint](https://img.shields.io/github/actions/workflow/status/goal-web/config/lint.yml?branch=main&label=Lint)
- Commit Activity: ![Commit Activity](https://img.shields.io/github/commit-activity/m/goal-web/config)
- Last Commit: ![Last Commit](https://img.shields.io/github/last-commit/goal-web/config)


## Contributing

- Issues and PRs are welcome. Please:
  - Keep API backward compatible.
  - Update docs/examples; document behavioral changes in Notes/FAQ.
  - Build and run basic tests before submitting.
- Use Issues or Discussions (if enabled) for proposals.

## Roadmap

- Module docs matrix for `goal/config` and `goal-cli/config`.
- Example app and smoke tests for combined sources.
- English/Chinese docs maintained.

## Changelog

- See Releases: <https://github.com/goal-web/config/releases>

## Security

- Do not commit secrets; prefer env vars or CLI encryption.
- Env vars have highest precedence; ideal for CI/containers.
- Encrypt production configs; decrypt at deploy. See CLI.

## FAQ

- File/URL load failures are logged and skipped; other sources continue.
- Override order: later sources win; env is highest.
- Use typed getters to avoid implicit conversions.
### Notes

1. Formats: built-in TOML/YAML/DotEnv; implement `contracts.Env` for custom sources.
2. Dynamic reload: call `config.Reload()` to re-compute sources at runtime; env stays highest.
3. Type-safety: use typed getters like `GetBool`, `GetString`.
