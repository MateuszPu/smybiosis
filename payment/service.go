package payment

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/paymentintent"
	"pay.me/mail"
	"pay.me/server"
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

func (service *Service) CreatePayment(stripeId string, userId uuid.UUID, email string) {
	//here we need to generate link to payment
	paramsPa := &stripe.PaymentIntentParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Amount:               stripe.Int64(1000),
		Currency:             stripe.String(string(stripe.CurrencyPLN)),
		ApplicationFeeAmount: stripe.Int64(100),
	}
	paramsPa.SetStripeAccount(stripeId)
	pi, _ := paymentintent.New(paramsPa)

	payment := payment{
		id:           uuid.New(),
		userId:       userId,
		linkId:       strings.ReplaceAll(uuid.New().String(), "-", ""),
		stripeId:     stripeId,
		clientSecret: pi.ClientSecret,
	}
	service.repository().save(payment)

	service.MailService.SendEmailWithPaymentLink(email, fmt.Sprintf("%spayment/%s", service.GlobalEnv.Host, payment.linkId))
}

func (service *Service) findPaymentByLinkId(linkId string) paymentData {
	payment, _ := service.repository().findByLinkId(linkId)
	return paymentData{
		StripeAccountId:    payment.stripeId,
		StripeClientSecret: payment.clientSecret,
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
