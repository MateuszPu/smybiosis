package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v71"
	"html/template"
	"pay.me/logging"
	"pay.me/server"
	"pay.me/user"
	"time"
)

func main() {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(requestLogger())
	env := server.Env{
		Env: server.EnvVariable("APP_ENV", "local"),
	}
	stripe.Key = "sk_test_51HTA7JDx7zNNd5t3lNXjrLaSX618luMWklkNUH86JVPfbfJpWtdnzTgQHU3w674dakLs4WyTbQQPenPXo7AF1yRP00SXmmlsYd"

	baseServer := server.BaseSever{
		Router: router.Group("/api/v1"),
		Logger: logging.Logger(env.Env),
	}
	userHandlers(baseServer)

	t := template.Must(template.ParseFiles("templates/index.html"))
	router.GET("/hello/:name", func(c *gin.Context) {
		param := c.Param("name")
		t.Execute(c.Writer, Data{"dupadupa", param})
	})

	router.Run(":8080")
}

type Data struct {
	Title string
	Text string
}

func userHandlers(baseServer server.BaseSever) {
	service := user.Service{Repository: user.CreateInMemoryRepo()}
	userHandler := user.Handler{
		BaseSever: &baseServer,
		Service:   &service,
	}
	userHandler.Routes()
}

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()

		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)
		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.BodySize = c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path
		//TODO: put it to ELK? stack?
		logging.HttpLogger(param).Info()
	}
}
