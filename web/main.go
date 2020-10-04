package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // here
	"github.com/stripe/stripe-go/v71"
	"pay.me/global"
	"pay.me/logging"
	"pay.me/mail"
	"pay.me/payment"
	"pay.me/server"
	"pay.me/user"
	"time"
)

func main() {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(requestLogger())
	stripe.Key = "sk_test_51HTA7JDx7zNNd5t3lNXjrLaSX618luMWklkNUH86JVPfbfJpWtdnzTgQHU3w674dakLs4WyTbQQPenPXo7AF1yRP00SXmmlsYd"

	env := server.Env{
		Env:  server.EnvVariable("APP_ENV", "local"),
		Host: server.EnvVariable("HOST", "http://localhost:8080/"),
	}

	baseServer := server.BaseSever{
		Router: router,
		Env:    &env,
		Logger: logging.Logger(env.Env),
	}
	service := mail.Service{
		Host:     server.EnvVariable("MAIL_HOST", "some"),
		Port:     server.EnvVariable("MAIL_PORT", "port"),
		Email:    server.EnvVariable("MAIL_LOGIN", "mail"),
		Password: server.EnvVariable("MAIL_PASSWORD", "pass"),
	}

	repo := payment.CreateInMemoryRepo()
	paymentService := payment.Service{
		Repository:  &repo,
		GlobalEnv:   &env,
		MailService: &service}
	globalHandler(&baseServer)
	userHandlers(&baseServer, &paymentService)
	paymentHandlers(&baseServer, &paymentService)

	router.Run(":8080")
}

func globalHandler(srv *server.BaseSever) {
	globalHandler := global.Handler{
		BaseServer: srv,
	}
	globalHandler.Routes()
}

func userHandlers(baseServer *server.BaseSever, paymentService *payment.Service) {
	repo := user.CreateInMemoryRepo()
	service := user.Service{Repository: &repo, Env: baseServer.Env}
	userHandler := user.Handler{
		BaseSever:      baseServer,
		Service:        &service,
		PaymentService: paymentService,
	}
	userHandler.Routes()
}

func paymentHandlers(baseServer *server.BaseSever, paymentService *payment.Service) {
	paymentHandler := payment.Handler{
		BaseSever: baseServer,
		Service:   paymentService,
	}
	paymentHandler.Routes()
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
