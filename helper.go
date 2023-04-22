package config

import (
	"github.com/goal-web/contracts"
	"os"
	"strings"
)

var singleton contracts.Config

func Get(key string) any {
	return singleton.Get(key)
}

func osEnvGetter(key string, defaultValue any) any {
	if value, ok := os.LookupEnv(strings.ToUpper(strings.ReplaceAll(key, ".", "_"))); ok {
		return value
	}
	return defaultValue
}

func GetString(key string) string {
	return singleton.GetString(key)
}
func GetInt64(key string) int64 {
	return singleton.GetInt64(key)
}
func GetInt32(key string) int32 {
	return singleton.GetInt32(key)
}
func GetInt16(key string) int16 {
	return singleton.GetInt16(key)
}
func GetInt8(key string) int8 {
	return singleton.GetInt8(key)
}
func GetInt(key string) int {
	return singleton.GetInt(key)
}
func GetUInt64(key string) uint64 {
	return singleton.GetUInt64(key)
}
func GetUInt32(key string) uint32 {
	return singleton.GetUInt32(key)
}
func GetUInt16(key string) uint16 {
	return singleton.GetUInt16(key)
}
func GetUInt8(key string) uint8 {
	return singleton.GetUInt8(key)
}
func GetUInt(key string) uint {
	return singleton.GetUInt(key)
}
func GetFloat64(key string) float64 {
	return singleton.GetFloat64(key)
}
func GetFloat(key string) float32 {
	return singleton.GetFloat(key)
}
func GetBool(key string) bool {
	return singleton.GetBool(key)
}

func Optional(key string, defaultValue any) any {
	return singleton.Optional(key, defaultValue)
}

func StringOptional(key string, defaultValue string) string {
	return singleton.StringOptional(key, defaultValue)
}
func Int64Optional(key string, defaultValue int64) int64 {
	return singleton.Int64Optional(key, defaultValue)
}
func Int32Optional(key string, defaultValue int32) int32 {
	return singleton.Int32Optional(key, defaultValue)
}
func Int16Optional(key string, defaultValue int16) int16 {
	return singleton.Int16Optional(key, defaultValue)
}
func Int8Optional(key string, defaultValue int8) int8 {
	return singleton.Int8Optional(key, defaultValue)
}
func IntOptional(key string, defaultValue int) int {
	return singleton.IntOptional(key, defaultValue)
}
func UInt64Optional(key string, defaultValue uint64) uint64 {
	return singleton.UInt64Optional(key, defaultValue)
}
func UInt32Optional(key string, defaultValue uint32) uint32 {
	return singleton.UInt32Optional(key, defaultValue)
}
func UInt16Optional(key string, defaultValue uint16) uint16 {
	return singleton.UInt16Optional(key, defaultValue)
}
func UInt8Optional(key string, defaultValue uint8) uint8 {
	return singleton.UInt8Optional(key, defaultValue)
}
func UIntOptional(key string, defaultValue uint) uint {
	return singleton.UIntOptional(key, defaultValue)
}
func Float64Optional(key string, defaultValue float64) float64 {
	return singleton.Float64Optional(key, defaultValue)
}
func FloatOptional(key string, defaultValue float32) float32 {
	return singleton.FloatOptional(key, defaultValue)
}
func BoolOptional(key string, defaultValue bool) bool {
	return singleton.BoolOptional(key, defaultValue)
}
