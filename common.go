package config

import (
	"io"
	"net/http"
	"os"

	"github.com/goal-web/supports/logs"
)

type EnvProvider func() []byte

// File 以本地文件作为配置数据源，读取失败时返回空字节并记录调试日志。
func File(path string) EnvProvider {
    return func() []byte {
        tmpBytes, err := os.ReadFile(path)
        if err != nil {
            log.Debug("File: " + err.Error())
        }
        return tmpBytes
    }
}

// Url 以远程 URL 作为配置数据源，网络或读取失败时返回 nil 并记录调试日志。
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
