package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	UserService    *Service
}

func (handler *Handler) Routes() *Handler {
	handler.BaseSever.Router.POST("/user", handler.createUser())
	handler.BaseSever.Router.POST("/email", handler.customerEmail())
	handler.BaseSever.Router.GET("/refresh/:id", handler.refreshRegistration())
	handler.BaseSever.Router.GET("/return/:id", handler.returnFromRegistration())
	return handler
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
			Description: context.PostForm("title"),
		}
		re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		if !re.MatchString(json.Email) {
			//todo: error handling
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Login should be valid email"})
			return
		}

		user, err := handler.UserService.findByEmail(json.Email)
		if err != nil {
			logrus.Info("User not found creating new one")
		}
		if user.email == json.Email {
			//it means user already exist in database send link to him
			//TODO: deal with situation when user was found but not finish registration process
			paymentId, err := handler.PaymentService.InitPayment(user.id, json.Currency, json.Amount, json.Description)
			if err != nil {
				//todo: log here
			}
			handler.PaymentService.GeneratePaymentLink(paymentId)
			context.SetCookie(server.COOKIE_NAME, user.cookieId, 0, "/", handler.BaseSever.Env.CookieHost, false, false)
			context.Redirect(http.StatusMovedPermanently, "/")
		} else {
			link, userId, err := handler.UserService.createUser(json.Email, context.GetHeader("user-agent"))
			if err != nil {
				handler.BaseSever.Logger.Errorf("User not found in database %s", err)
				context.Redirect(http.StatusFound, "/404")
			}
			go func(userId uuid.UUID, json request) {
				//todo: error handling
				handler.PaymentService.InitPayment(userId, json.Currency, json.Amount, json.Description)
			}(*userId, json)
			context.Redirect(http.StatusMovedPermanently, *link)
		}

	}
}

func (handler *Handler) customerEmail() gin.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Question string `json:"question"`
	}
	return func(context *gin.Context) {
		email := context.PostForm("email")
		question := context.PostForm("question")
		handler.UserService.sendCustomerEmail(email, question)
		context.Redirect(http.StatusMovedPermanently, "/")
	}
}

func (handler *Handler) refreshRegistration() gin.HandlerFunc {
	return func(context *gin.Context) {
		linkId := context.Param("id")
		usr, err := handler.UserService.findByLinkId(linkId)

		if err != nil {
			handler.BaseSever.Logger.Errorf("User not found in database %s", err)
			context.Redirect(http.StatusFound, "/404")
			return
		}

		if usr.status == STRIPE_CONFIRMED {
			handler.BaseSever.Logger.Errorf("User is already registered in stripe", err)
			context.Redirect(http.StatusFound, "/404")
			return
		}

		link, err := handler.UserService.registrationLink(usr.stripeId, usr.linkId)
		if err != nil {
			handler.BaseSever.Logger.Errorf("Error while generating stripe link %s", err)
			context.Redirect(http.StatusFound, "/404")
			return
		}
		context.Redirect(http.StatusMovedPermanently, link)
	}

}

func (handler *Handler) returnFromRegistration() gin.HandlerFunc {
	t := template.Must(template.ParseFiles("templates/return.html"))
	return func(context *gin.Context) {
		handler.BaseSever.Logger.Infof("Registration finished for %s", context.Param("id"))
		linkId := context.Param("id")

		usr, err := handler.UserService.findByLinkId(linkId)
		if err != nil {
			handler.BaseSever.Logger.Errorf("Cannot finish registration for linkId: %s. Err: %s", linkId, err)
			context.Redirect(http.StatusFound, "/404")
			return
		}
		handler.PaymentService.GenerateFirstPaymentLink(usr.id)

		link, err := handler.UserService.registrationLink(usr.stripeId, usr.linkId)
		println(link)

		//todo: here we need to decide what information dispaly on webpage and should we send email? or maybe just the same webpage but the email content will be based on staus of stripe registration
		context.SetCookie(server.COOKIE_NAME, usr.cookieId, 0, "/", handler.BaseSever.Env.CookieHost, false, false)
		t.Execute(context.Writer, nil)
	}
}
