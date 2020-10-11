package server

import (
	"github.com/gin-gonic/gin"
	"os"
	"pay.me/v4/logging"
)

type BaseSever struct {
	Router *gin.Engine
	Env    *Env
	Logger *logging.StandardLogger
}

func (base *BaseSever) ApiGroup() *gin.RouterGroup {
	return base.Router.Group("/api/v1")
}

type Env struct {
	Env  string
	Host string
	StripeKey string
}

func EnvVariable(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
