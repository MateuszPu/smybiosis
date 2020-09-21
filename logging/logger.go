package logging

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	*logrus.Entry
}

// HttpLogger initializes the standard logger
func HttpLogger(params gin.LogFormatterParams) *StandardLogger {
	var baseLogger = logrus.WithFields(logrus.Fields{
		//"env":          app.EnvVariable(env.APP_ENV, "local"),
		"ip":           params.ClientIP,
		"protocol":     params.Request.Proto,
		"userAgent":    params.Request.UserAgent(),
		"method":       params.Method,
		"path":         params.Path,
		"responseCode": params.StatusCode,
		"latency":      params.Latency,
	})
	baseLogger.Logger.Formatter = &logrus.JSONFormatter{}

	return &StandardLogger{baseLogger}
}

func Logger(env string) *StandardLogger {
	var baseLogger = logrus.WithFields(logrus.Fields{
		"env": env,
	})
	baseLogger.Logger.Formatter = &logrus.JSONFormatter{}

	return &StandardLogger{baseLogger}
}

