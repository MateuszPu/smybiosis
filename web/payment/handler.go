package payment

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"pay.me/v4/server"
)

type Handler struct {
	BaseSever *server.BaseSever
	Service   *Service
}

func (handler *Handler) Routes() *Handler {
	handler.BaseSever.Router.GET("/payment/:id", handler.paymentLink())
	return handler
}

func (handler *Handler) paymentLink() gin.HandlerFunc {
	t := template.Must(template.ParseFiles("templates/payment.html"))
	return func(context *gin.Context) {
		param := context.Param("id")
		data, err := handler.Service.createStripePayment(param)
		if err != nil {
			//todo: logger
			context.Redirect(http.StatusFound, "/error")
		}
		t.Execute(context.Writer, data)
	}
}
