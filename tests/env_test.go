package tests

import (
	"fmt"
	"github.com/goal-web/config"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestToml(t *testing.T) {
	env := config.NewToml(config.File("config.toml"))
	var fields = env.Load()
	fmt.Println(fields)
	assert.NotNil(t, fields)
	conf := config.New(env, nil)
	assert.True(t, conf.GetBool("app.debug"))
}

func TestYaml(t *testing.T) {
	env := config.NewYaml(config.File("config.yaml"))
	var fields = env.Load()
	fmt.Println(fields)
	assert.NotNil(t, fields)
	conf := config.New(env, nil)
	assert.True(t, conf.GetBool("app.debug"))
}

func TestDotEnv(t *testing.T) {
	env := config.NewDotEnv(config.File("config.env"))
	var fields = env.Load()
	fmt.Println(fields)
	assert.NotNil(t, fields)
	conf := config.New(env, nil)
	assert.True(t, conf.GetBool("app.debug"))
}

func TestOsEnvGetter(t *testing.T) {
	key := "app.env"
	osEnvKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
	assert.Equal(t, osEnvKey, "APP_ENV")
	err := os.Setenv("APP_ENV", "testing")
	assert.NoError(t, err, err)
	env := config.NewToml(config.File("config.toml"))
	var fields = env.Load()
	assert.NotNil(t, fields)
	conf := config.New(env, nil)
	assert.Equal(t, conf.GetString("app.env"), "testing")
}
