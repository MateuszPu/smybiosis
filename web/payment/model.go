package payment

import "github.com/google/uuid"

type payment struct {
	id           uuid.UUID
	userId       uuid.UUID
	linkId       string
	stripeId     string
	clientSecret string
}

type paymentData struct {
	StripeAccountId    string
	StripeClientSecret string
}
