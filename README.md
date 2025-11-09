# Goal-Config

Goal 配置组件，支持从文件（TOML、YAML、DotEnv）和环境变量中加载配置，并自动合并优先级。

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
    app := application.New()

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
    app := application.New()
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

### 环境变量覆盖示例

```go
// 将 app.env 的值覆盖为 testing（优先级最高）
os.Setenv("APP_ENV", "testing")
env := config.NewToml(config.File("env.toml"))
conf := config.New(env, nil)
fmt.Println(conf.GetString("app.env")) // 输出: testing
```

### 注意事项

1. **文件格式支持**：内置支持 TOML、YAML 和 DotEnv。你可以通过实现 `contracts.Env` 接口扩展其他来源或处理流程（例如自定义拉取、预处理等）。
2. **动态加载**：配置在应用启动时加载，运行时修改文件或环境变量需重启应用生效。
3. **类型安全**：使用 `GetBool`、`GetString` 等方法确保类型安全。