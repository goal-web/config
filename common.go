package config

import (
	"io"
	"net/http"
	"os"

	"github.com/goal-web/supports/logs"
)

type EnvProvider func() []byte

func File(path string) EnvProvider {
	return func() []byte {
		tmpBytes, err := os.ReadFile(path)
		if err != nil {
			log.Debug("File: " + err.Error())
		}
		return tmpBytes
	}
}

func Url(url string) EnvProvider {
	return func() []byte {
		res, err := http.Get(url)
		if err != nil {
			log.Debug("File: " + err.Error())
			return nil
		}
		tmpBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Debug("File: " + err.Error())
		}
		return tmpBytes
	}
}

var log = logs.Default()
