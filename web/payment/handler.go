package payment

import (
	"github.com/gin-gonic/gin"
	"html/template"
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
		data := handler.Service.createStripePayment(param)
		t.Execute(context.Writer, data)
	}
}
