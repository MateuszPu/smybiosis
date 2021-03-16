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
	handler.BaseSever.Router.GET("/payments/:id/success", handler.success())
	handler.BaseSever.Router.GET("/payments/:id", handler.paymentLink())
	return handler
}

func (handler *Handler) paymentLink() gin.HandlerFunc {
	t := template.Must(template.ParseFiles("templates/payment/init.html"))
	return func(context *gin.Context) {
		param := context.Param("id")
		data, err := handler.Service.createStripePayment(param)
		if err != nil {
			//todo: logger
			context.Redirect(http.StatusFound, "/404")
		}
		t.Execute(context.Writer, data)
	}
}

func (handler *Handler) success() gin.HandlerFunc {
	t := template.Must(template.ParseFiles("templates/payment/success.html"))
	return func(context *gin.Context) {
		param := context.Param("id")
		if handler.Service.successPayment(param) != nil {
			//todo: logger
			context.Redirect(http.StatusFound, "/404")
		}
		t.Execute(context.Writer, paymentData{})
	}
}