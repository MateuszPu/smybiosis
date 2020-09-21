package server

import (
	"github.com/gin-gonic/gin"
	"os"
	"pay.me/logging"
)

type BaseSever struct {
	Router *gin.RouterGroup
	Logger *logging.StandardLogger
}

type Env struct {
	Env string
}

func EnvVariable(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
