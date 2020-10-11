package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"html/template"
	"net/http"
	"pay.me/v4/payment"
	"pay.me/v4/server"
	"regexp"
	"strconv"
)

type Handler struct {
	BaseSever      *server.BaseSever
	PaymentService *payment.Service
	Service        *Service
}

func (handler *Handler) Routes() *Handler {
	//authGr := handler.BaseSever.ApiGroup().Group("/user")
	handler.BaseSever.Router.POST("/user", handler.createUser())
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
			context.Redirect(http.StatusFound, "/404")
			return
		}
		t.Execute(context.Writer, nil)
		go func(userId uuid.UUID) {
			handler.BaseSever.Logger.Infof("Sending payment link for user %s", userId.String())
			handler.PaymentService.GenerateFirstPaymentLink(userId)
		}(usr.id)

	}
}

func (handler *Handler) createUser() gin.HandlerFunc {
	type request struct {
		Email       string  `json:"email"`
		Amount      float64 `json:"amount"`
		Currency    string  `json:"currency"`
		Description string  `json:"description"`
	}

	return func(context *gin.Context) {
		amount, _ := strconv.ParseFloat(context.PostForm("amount"), 64)
		json := request{
			Email:       context.PostForm("email"),
			Amount:      amount,
			Currency:    context.PostForm("currency"),
			Description: context.PostForm("description"),
		}
		re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		if !re.MatchString(json.Email) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Login should be valid email"})
			return
		}
		userService := *handler.Service

		user, _ := userService.findByEmail(json.Email)
		if user.email == json.Email {
			//it means user alread exist in database send link to him
			//TODO: deal with situation when user was found but not finish registration process
			paymentId, err := handler.PaymentService.CreatePayment(user.id, json.Currency, json.Amount, json.Description)
			if err != nil {
				//todo: log here
			}
			handler.PaymentService.GeneratePaymentLink(paymentId)
			context.Redirect(http.StatusMovedPermanently, "/")
		} else {
			usr, err := userService.createUser(json.Email)
			if err != nil {
				handler.BaseSever.Logger.Errorf("Error while saving user in database %s", err)
				context.Redirect(http.StatusFound, "/404")
				return
			}
			link, err := userService.stripeLink(usr.stripeId, usr.linkId)
			if err != nil {
				handler.BaseSever.Logger.Errorf("Error while generating stripe link %s", err)
				context.Redirect(http.StatusFound, "/404")
				return
			}
			go func(userId uuid.UUID, json request) {
				//todo: error handling
				handler.PaymentService.CreatePayment(userId, json.Currency, json.Amount, json.Description)
			}(usr.id, json)
			context.Redirect(http.StatusMovedPermanently, link)
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
			context.Redirect(http.StatusFound, "/404")
			return
		}

		link, err := userService.stripeLink(usr.stripeId, usr.linkId)
		if err != nil {
			handler.BaseSever.Logger.Errorf("Error while generating stripe link %s", err)
			context.Redirect(http.StatusFound, "/404")
			return
		}
		context.Redirect(http.StatusMovedPermanently, link)
	}

}
