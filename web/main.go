package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // here
	"github.com/stripe/stripe-go/v71"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"pay.me/v4/global"
	"pay.me/v4/logging"
	"pay.me/v4/mail"
	"pay.me/v4/payment"
	"pay.me/v4/server"
	"pay.me/v4/user"
	"time"
)

func main() {
	db, _ := sql.Open("postgres", "dbname=postgres host=localhost user=user password=pass sslmode=disable")
	//if err1 != nil {
	//	println("a")
	//}
	//err2 := db.Ping()
	//if err2 != nil {
	//	println("a")
	//}
	boil.SetDB(db)
	boil.DebugMode = true

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

	paymentService := payment.Service{
		Repository:  payment.CreateSqlRepo(db),
		GlobalEnv:   &env,
		MailService: &service}
	globalHandler(&baseServer)
	userHandlers(&baseServer, &paymentService, db)
	paymentHandlers(&baseServer, &paymentService)

	router.Run(":8080")
}

func globalHandler(srv *server.BaseSever) {
	globalHandler := global.Handler{
		BaseServer: srv,
	}
	globalHandler.Routes()
}

func userHandlers(baseServer *server.BaseSever, paymentService *payment.Service, db *sql.DB) {
	repo := user.CreateSqlRepo(db)
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
