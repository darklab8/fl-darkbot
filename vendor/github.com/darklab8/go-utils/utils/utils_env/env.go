package utils_env

import (
	"os"
	"strconv"

	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/utils/utils_logus"
)

type EnvConfig struct{}

func NewEnvConfig() EnvConfig {
	return EnvConfig{}
}

func (e EnvConfig) GetEnvWithDefault(key string, default_ string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	} else {
		return default_
	}
}

func (e EnvConfig) GetEnvBool(key string) bool {
	return os.Getenv(key) == "true"
}

func (e EnvConfig) GetEnvBoolWithDefault(key string, default_ bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		return value == "true"
	}
	return default_
}

func (e EnvConfig) GetEnv(key string) string {
	return os.Getenv(key)
}

func (e EnvConfig) GetIntWithDefault(key string, default_ int) int {
	if value, ok := os.LookupEnv(key); ok {
		int_value, err := strconv.Atoi(value)
		utils_logus.Log.CheckPanic(err, "expected to be int", typelog.String("key", key))
		return int_value
	}
	return default_
}
