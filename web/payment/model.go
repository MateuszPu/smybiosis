package payment

import "github.com/google/uuid"

type payment struct {
	id                uuid.UUID
	userId            uuid.UUID
	linkId            string
	stripConnectedAcc string
}

type paymentData struct {
	StripeConnectedAccountId string
	StripeClientSecret       string
}
