package config

import (
	"os"
	"strings"
)

// ToEnvKey 将 key 转换为大写，并将 . 替换为 _
func ToEnvKey(key string) string {
	return strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
}

func osEnvGetter(key string, defaultValue any) any {
	if value, ok := os.LookupEnv(ToEnvKey(key)); ok {
		return value
	}
	return defaultValue
}

func GetString(key string) string {
	return Default().GetString(key)
}
func GetInt64(key string) int64 {
	return Default().GetInt64(key)
}
func GetInt32(key string) int32 {
	return Default().GetInt32(key)
}
func GetInt16(key string) int16 {
	return Default().GetInt16(key)
}
func GetInt8(key string) int8 {
	return Default().GetInt8(key)
}
func GetInt(key string) int {
	return Default().GetInt(key)
}
func GetUInt64(key string) uint64 {
	return Default().GetUInt64(key)
}
func GetUInt32(key string) uint32 {
	return Default().GetUInt32(key)
}
func GetUInt16(key string) uint16 {
	return Default().GetUInt16(key)
}
func GetUInt8(key string) uint8 {
	return Default().GetUInt8(key)
}
func GetUInt(key string) uint {
	return Default().GetUInt(key)
}
func GetFloat64(key string) float64 {
	return Default().GetFloat64(key)
}
func GetFloat(key string) float32 {
	return Default().GetFloat(key)
}
func GetBool(key string) bool {
	return Default().GetBool(key)
}

func Optional(key string, defaultValue any) any {
	return Default().Optional(key, defaultValue)
}

func StringOptional(key string, defaultValue string) string {
	return Default().StringOptional(key, defaultValue)
}
func Int64Optional(key string, defaultValue int64) int64 {
	return Default().Int64Optional(key, defaultValue)
}
func Int32Optional(key string, defaultValue int32) int32 {
	return Default().Int32Optional(key, defaultValue)
}
func Int16Optional(key string, defaultValue int16) int16 {
	return Default().Int16Optional(key, defaultValue)
}
func Int8Optional(key string, defaultValue int8) int8 {
	return Default().Int8Optional(key, defaultValue)
}
func IntOptional(key string, defaultValue int) int {
	return Default().IntOptional(key, defaultValue)
}
func UInt64Optional(key string, defaultValue uint64) uint64 {
	return Default().UInt64Optional(key, defaultValue)
}
func UInt32Optional(key string, defaultValue uint32) uint32 {
	return Default().UInt32Optional(key, defaultValue)
}
func UInt16Optional(key string, defaultValue uint16) uint16 {
	return Default().UInt16Optional(key, defaultValue)
}
func UInt8Optional(key string, defaultValue uint8) uint8 {
	return Default().UInt8Optional(key, defaultValue)
}
func UIntOptional(key string, defaultValue uint) uint {
	return Default().UIntOptional(key, defaultValue)
}
func Float64Optional(key string, defaultValue float64) float64 {
	return Default().Float64Optional(key, defaultValue)
}
func FloatOptional(key string, defaultValue float32) float32 {
	return Default().FloatOptional(key, defaultValue)
}
func BoolOptional(key string, defaultValue bool) bool {
	return Default().BoolOptional(key, defaultValue)
}
