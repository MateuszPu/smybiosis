package user

import (
	"github.com/google/uuid"
	models "pay.me/v4/database-models"
	"time"
)

const ACCOUNT_CREATED = "ACCOUNT_CREATED"
const STRIPE_CONFIRMED = "STRIPE_CONFIRMED"

type user struct {
	id        uuid.UUID
	cookieId  string
	stripeId  string
	linkId    string
	email     string
	status    string
	createdAt time.Time
}

func (u *user) from(dbUser *models.User) user {
	return user{
		id:        uuid.MustParse(dbUser.ID),
		cookieId:  dbUser.CookieID,
		stripeId:  dbUser.StripeAccount,
		linkId:    dbUser.LinkRegistration,
		email:     dbUser.Email,
		status:    dbUser.Status,
		createdAt: dbUser.CreatedAt,
	}
}
