package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"html/template"
	"net/http"
	"pay.me/v4/payment"
	"pay.me/v4/server"
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

	handler.BaseSever.Router.GET("/refresh/:id", handler.refreshRegistration())
	handler.BaseSever.Router.GET("/confirm/:id", handler.finishRegistration())
	return handler
}

func (handler *Handler) finishRegistration() gin.HandlerFunc {
	t := template.Must(template.ParseFiles("templates/finish.html"))
	return func(context *gin.Context) {
		linkId := context.Param("id")
		handler.BaseSever.Logger.Infof("Registration finished for %s", linkId)
		userService := *handler.Service
		usr, err := userService.finishedStripeRegistration(linkId)
		if err != nil {
			handler.BaseSever.Logger.Errorf("Account does not exist for linkId: %s. Cannot finish registration", linkId)
			context.Redirect(http.StatusSeeOther, "/error")
			return
		}
		t.Execute(context.Writer, nil)
		go func(stripeId string, userId uuid.UUID, email string) {
			handler.BaseSever.Logger.Infof("Generate payment for user %s", userId.String())
			handler.PaymentService.CreatePayment(stripeId, userId, email)
		}(usr.stripeId, usr.id, usr.email)

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

		user := userService.findByEmail(json.Email)
		if user.email == json.Email {
			//TODO: deal with situation when user was found but not finish registration process
			handler.PaymentService.CreatePayment(user.stripeId, user.id, user.email)
		} else {
			usr, err := userService.createUser(json.Email)
			if err != nil {
				handler.BaseSever.Logger.Errorf("Error while saving user in database %s", err)
				context.Redirect(http.StatusSeeOther, "/error")
				return
			}
			link, err := userService.stripeLink(usr.stripeId, usr.linkId)
			if err != nil {
				handler.BaseSever.Logger.Errorf("Error while generating stripe link %s", err)
				context.Redirect(http.StatusSeeOther, "/error")
				return
			}
			context.Redirect(301, link)
		}

	}
}

func (handler *Handler) refreshRegistration() gin.HandlerFunc {
	return func(context *gin.Context) {
		userService := *handler.Service
		linkId := context.Param("id")
		usr, err := userService.findByLinkId(linkId)

		if err != nil {
			handler.BaseSever.Logger.Errorf("User not found in database %s", err)
			context.Redirect(http.StatusSeeOther, "/error")
			return
		}

		link, err := userService.stripeLink(usr.stripeId, usr.linkId)
		if err != nil {
			handler.BaseSever.Logger.Errorf("Error while generating stripe link %s", err)
			context.Redirect(http.StatusSeeOther, "/error")
			return
		}
		context.Redirect(301, link)
	}

}
