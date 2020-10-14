package global

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"pay.me/v4/payprovider"
	"pay.me/v4/server"
	"strings"
)

type Handler struct {
	BaseServer *server.BaseSever
}

func (handler *Handler) Routes() *Handler {
	handler.BaseServer.Router.GET("/404", handler.error())
	handler.BaseServer.Router.GET("/", handler.index())
	return handler
}

func (handler *Handler) error() gin.HandlerFunc {
	t := template.Must(template.ParseFiles("templates/404.html"))
	return func(context *gin.Context) {
		t.Execute(context.Writer, nil)
	}
}

func (handler *Handler) index() gin.HandlerFunc {
	t := template.Must(template.ParseFiles("templates/index.html"))
	type data struct {
		Currencies []string
	}
	return func(context *gin.Context) {
		currencies := []string{}
		for _, currency := range payprovider.AllCurrencies() {
			currencies = append(currencies, strings.ToUpper(currency.Value))
		}
		d := data{currencies}
		t.Execute(context.Writer, d)
	}
}
