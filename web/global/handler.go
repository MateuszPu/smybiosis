package global

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"pay.me/v4/payprovider"
	"pay.me/v4/server"
	"sort"
	"strings"
)

type Handler struct {
	BaseServer  *server.BaseSever
	Repository *repository
}

func (handler *Handler) Routes() *Handler {
	handler.BaseServer.Router.GET("/404", handler.error())
	handler.BaseServer.Router.GET("/", handler.index())
	return handler
}

func (handler *Handler) error() gin.HandlerFunc {
	t := template.Must(template.ParseFiles("templates/404.html"))
	return func(context *gin.Context) {
		context.AbortWithStatus(404)
		t.Execute(context.Writer, nil)
	}
}

func (handler *Handler) index() gin.HandlerFunc {
	t := template.Must(template.ParseFiles("templates/index.html"))
	type data struct {
		Currencies []string
		ButtonName string
		Email      string
		Amount     string
		Title      string
	}
	return func(context *gin.Context) {
		currencies := []string{}
		for _, currency := range payprovider.AllCurrencies() {
			if currency.Value == "USD" || currency.Value == "EUR" || currency.Value == "GBP" {
				continue
			}
			currencies = append(currencies, strings.ToUpper(currency.Value))
		}
		sort.Strings(currencies)

		userCookieId, err := context.Cookie(server.COOKIE_NAME)
		if err != nil {
			d := data{Currencies: currencies, ButtonName: "Next"}
			t.Execute(context.Writer, d)
		} else {
			userDetails, _ := handler.repository().findUserDetailsBy(userCookieId)
			d := data{Currencies: currencies, ButtonName: "Send me a link", Email: userDetails.Email, Amount: fmt.Sprintf("%g", userDetails.Amount), Title: userDetails.Title}
			t.Execute(context.Writer, d)
		}

	}
}

func (h *Handler) repository() repository {
	return *h.Repository
}
