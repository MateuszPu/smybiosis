package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"html/template"
	"net/http"
	"pay.me/payment"
	"pay.me/server"
	"regexp"
)

type Handler struct {
	BaseSever      *server.BaseSever
	PaymentService *payment.Service
	Service        *Service
}

func (handler *Handler) Routes() *Handler {
	authGr := handler.BaseSever.ApiGroup().Group("/user")
	authGr.POST("", handler.createUser())

	//handler.BaseSever.Router.GET("/refresh/:id", handler.finishRegistration())
	handler.BaseSever.Router.GET("/confirm/:id", handler.finishRegistration())
	return handler
}

func (handler *Handler) finishRegistration() gin.HandlerFunc {
	t := template.Must(template.ParseFiles("templates/finish.html"))
	return func(context *gin.Context) {
		linkId := context.Param("id")
		userService := *handler.Service
		usr, err := userService.finishedStripeRegistration(linkId)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Email does not exist"})
			return
		}
		t.Execute(context.Writer, nil)
		go func(stripeId string, userId uuid.UUID) {
			handler.PaymentService.CreatePayment(stripeId, userId)
		}(usr.stripeId, usr.id)

	}
}

func (handler *Handler) createUser() gin.HandlerFunc {
	type request struct {
		Email string `json:"email"`
	}

	return func(context *gin.Context) {
		var json request
		if err := context.ShouldBindJSON(&json); err != nil {
			handler.BaseSever.Logger.Errorf("Error while parsing json %s", err)
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "error during parsing json"})
			return
		}
		re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		if !re.MatchString(json.Email) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Login should be valid email"})
			return
		}
		userService := *handler.Service

		//todo: alternative path when user already exist in database
		usr, err := userService.createUser(json.Email)
		if err != nil {
			handler.BaseSever.Logger.Errorf("Error while saving user in database %s", err)
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Error user creation"})
			return
		}
		link, err := userService.stripeLink(usr.stripeId, usr.linkId)
		if err != nil {
			handler.BaseSever.Logger.Errorf("Error while generating stripe link %s", err)
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Error stripe link creation"})
			return
		}
		context.Redirect(301, link)
	}
}
