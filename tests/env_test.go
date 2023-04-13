package tests

import (
	"fmt"
	"github.com/goal-web/config"
	"github.com/stretchr/testify/assert"
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
