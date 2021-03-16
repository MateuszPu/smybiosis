package payment

import (
	"github.com/google/uuid"
	models "pay.me/v4/database-models"
)

const PAYMENT_CREATED = "PAYMENT_CREATED"
const PAYMENT_INITIALED = "PAYMENT_INITIALED"
const PAYMENT_SUCCESS = "PAYMENT_SUCCESS"
const PAYMENT_CANCELED = "PAYMENT_CANCELED"

type payment struct {
	id              uuid.UUID
	linkHash        uuid.UUID
	confirmedHash   uuid.UUID
	canceledHash    uuid.UUID
	currency        string
	amount          float64
	description     string
	stripeAccId     string
	stripeIdPayment string
	email           string
	status          string
}

func (p *payment) from(dbPayment *models.Payment, user *models.User) payment {
	return payment{
		id:              uuid.MustParse(dbPayment.ID),
		linkHash:        uuid.MustParse(dbPayment.LinkHash),
		confirmedHash:   uuid.MustParse(dbPayment.ConfirmedHash),
		canceledHash:    uuid.MustParse(dbPayment.CanceledHash),
		status:          dbPayment.Status,
		currency:        dbPayment.Currency,
		amount:          dbPayment.Amount,
		description:     dbPayment.Description,
		stripeAccId:     user.StripeAccount,
		stripeIdPayment: dbPayment.StripeIDPayment.String,
		email:           user.Email,
	}
}

type paymentData struct {
	StripeConnectedAccountId string
	StripeClientSecret       string
}
