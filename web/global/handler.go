package global

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"pay.me/v4/server"
)

type Handler struct {
	BaseServer *server.BaseSever
}

func (handler *Handler) Routes() *Handler {
	handler.BaseServer.Router.GET("/error", handler.error())
	return handler
}

func (handler *Handler) error() gin.HandlerFunc {
	t := template.Must(template.ParseFiles("templates/error.html"))
	return func(context *gin.Context) {
		t.Execute(context.Writer, nil)
	}
}
