package tests

import (
	"fmt"
	"github.com/goal-web/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToml(t *testing.T) {
	toml := config.NewToml(config.File("env.toml"))
	var fields = toml.Load()
	fmt.Println(fields)
	assert.NotNil(t, fields)
}

func TestDotEnv(t *testing.T) {
	toml := config.NewDotEnv(config.File("env.env"))
	var fields = toml.Load()
	fmt.Println(fields)
	assert.NotNil(t, fields)
}
