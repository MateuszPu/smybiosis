package payment

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/checkout/session"
	"pay.me/v4/mail"
	"pay.me/v4/server"
	"strings"
)

type Service struct {
	Repository  *repository
	MailService *mail.Service
	GlobalEnv   *server.Env
}

func (service *Service) repository() repository {
	return *service.Repository
}

func (service *Service) CreatePayment(stripConnectedAcc string, userId uuid.UUID, email string) {
	//here we need to generate link to payment
	//paramsPa := &stripe.PaymentIntentParams{
	//	PaymentMethodTypes: stripe.StringSlice([]string{
	//		"card",
	//	}),
	//	Amount:               stripe.Int64(1000),
	//	Currency:             stripe.String(string(stripe.CurrencyPLN)),
	//	ApplicationFeeAmount: stripe.Int64(100),
	//}
	//paramsPa.SetStripeAccount(stripConnectedAcc)
	//pi, _ := paymentintent.New(paramsPa)
	payment := payment{
		id:                uuid.New(),
		userId:            userId,
		linkId:            strings.ReplaceAll(uuid.New().String(), "-", ""),
		stripConnectedAcc: stripConnectedAcc,
	}
	service.repository().save(payment)

	service.MailService.SendEmailWithPaymentLink(email, fmt.Sprintf("%spayment/%s", service.GlobalEnv.Host, payment.linkId))
}

func (service *Service) createStripePayment(linkId string) paymentData {
	//TODO: error handling?
	payment, _ := service.repository().findByLinkId(linkId)
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Name:     stripe.String("Stainless Steel Water Bottle"),
				Amount:   stripe.Int64(50000000),
				Currency: stripe.String(string(stripe.CurrencyPLN)),
				Quantity: stripe.Int64(1),
			},
		},
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			ApplicationFeeAmount: stripe.Int64(5000),
		},
		SuccessURL: stripe.String("https://example.com/success"),
		CancelURL:  stripe.String("https://example.com/cancel"),
	}

	params.SetStripeAccount(payment.stripConnectedAcc)
	s, _ := session.New(params)

	return paymentData{
		StripeConnectedAccountId: payment.stripConnectedAcc,
		StripeClientSecret:       s.ID,
	}

}

type repository interface {
	save(payment)
	findByLinkId(string) (payment, error)
}

func CreateInMemoryRepo() repository {
	return &RepositoryInMemory{inMemory: make(map[uuid.UUID]payment)}
}

type RepositoryInMemory struct {
	inMemory map[uuid.UUID]payment
}

func (repo *RepositoryInMemory) save(payment payment) {
	repo.inMemory[payment.id] = payment
}

func (repo *RepositoryInMemory) findByLinkId(linkId string) (payment, error) {
	for _, pymnt := range repo.inMemory {
		if pymnt.linkId == linkId {
			return pymnt, nil
		}
	}
	return payment{}, nil
}
