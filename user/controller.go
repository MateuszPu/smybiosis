package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pay.me/server"
	"regexp"
)

type Handler struct {
	BaseSever *server.BaseSever
	Service   *Service
}

func (handler *Handler) Routes() *Handler {
	authGr := handler.BaseSever.Router.Group("/user")
	authGr.POST("", handler.createUser())

	return handler
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
		link, err := userService.stripeLink(usr.stripeId)
		if err != nil {
			handler.BaseSever.Logger.Errorf("Error while generating stripe link %s", err)
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Error stripe link creation"})
			return
		}
		context.Redirect(301, link)
	}
}
