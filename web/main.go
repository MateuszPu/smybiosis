package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // here
	"pay.me/v4/database"
	"pay.me/v4/global"
	"pay.me/v4/logging"
	"pay.me/v4/mail"
	"pay.me/v4/payment"
	"pay.me/v4/payprovider"
	"pay.me/v4/server"
	"pay.me/v4/user"
	"time"
)

func main() {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(requestLogger())

	env := server.Env{
		Env:       server.EnvVariable("APP_ENV", "local"),
		Host:      server.EnvVariable("HOST", "http://localhost:8080/"),
		StripeKey: server.EnvVariable("STRIPE_KEY", "sk_test_51HTA7JDx7zNNd5t3lNXjrLaSX618luMWklkNUH86JVPfbfJpWtdnzTgQHU3w674dakLs4WyTbQQPenPXo7AF1yRP00SXmmlsYd"),
	}
	paymentProvider := payprovider.Service{Env: &env}.Init()
	baseServer := server.BaseSever{
		Router: router,
		Env:    &env,
		Logger: logging.Logger(env.Env),
	}

	db := database.SqlDatabase{
		DriverName: server.EnvVariable("DB_DRIVER", "postgres"),
		Name:       server.EnvVariable("DB_NAME", "postgres"),
		Host:       server.EnvVariable("DB_HOST", "localhost"),
		User:       server.EnvVariable("DB_USER", "user"),
		Password:   server.EnvVariable("DB_PASSWORD", "pass"),
		Logger:     baseServer.Logger,
	}.CreateDb()

	service := mail.Service{
		Host:     server.EnvVariable("MAIL_HOST", "some"),
		Port:     server.EnvVariable("MAIL_PORT", "port"),
		Email:    server.EnvVariable("MAIL_LOGIN", "mail"),
		Password: server.EnvVariable("MAIL_PASSWORD", "pass"),
	}

	paymentService := payment.Service{
		Repository:      payment.CreateSqlRepo(db),
		GlobalEnv:       &env,
		PaymentProvider: paymentProvider,
		Commission:      0.005,
		MailService:     &service}
	globalHandler(&baseServer)
	userHandlers(&baseServer, &paymentService, paymentProvider, db)
	paymentHandlers(&baseServer, &paymentService)

	router.Run(":8080")
}

func globalHandler(srv *server.BaseSever) {
	globalHandler := global.Handler{
		BaseServer: srv,
	}
	globalHandler.Routes()
}

func userHandlers(baseServer *server.BaseSever, paymentService *payment.Service, paymentProvider *payprovider.Service, db *sql.DB) {
	repo := user.CreateSqlRepo(db)
	service := user.Service{Repository: &repo, Env: baseServer.Env}
	userHandler := user.Handler{
		BaseSever:       baseServer,
		Service:         &service,
		PaymentService:  paymentService,
		PaymentProvider: paymentProvider,
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
