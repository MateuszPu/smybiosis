package payment

import (
	"fmt"
	"github.com/google/uuid"
	"pay.me/v4/mail"
	"pay.me/v4/payprovider"
	"pay.me/v4/server"
	"strings"
)

type Service struct {
	Repository      repository
	MailService     *mail.Service
	PaymentProvider *payprovider.PaymentProvider
	GlobalEnv       *server.Env
	Commission      float64
}

func (service *Service) paymentProvider() payprovider.PaymentProvider {
	return *service.PaymentProvider
}

func (service *Service) InitPayment(userId uuid.UUID, currency string, amount float64, description string) (uuid.UUID, error) {
	payment := payment{
		id:          uuid.New(),
		currency:    currency,
		amount:      amount,
		description: description,
	}
	err := service.Repository.save(payment, userId)
	return payment.id, err
}

func (service *Service) GenerateFirstPaymentLink(userId uuid.UUID) {
	payment, err := service.Repository.byUserId(userId)
	if err != nil {
		//todo; log here
		return
	}
	service.MailService.SendEmailWithPaymentLink(payment.email, fmt.Sprintf("%spayment/%s", service.GlobalEnv.Host, payment.linkHash.String()))
}

func (service *Service) GeneratePaymentLink(id uuid.UUID) {
	payment, err := service.Repository.byId(id)
	if err != nil {
		//todo; log here
		return
	}
	service.MailService.SendEmailWithPaymentLink(payment.email, fmt.Sprintf("%spayment/%s", service.GlobalEnv.Host, payment.linkHash.String()))
}

func (service *Service) createStripePayment(linkId string) (paymentData, error) {
	//TODO: error handling?
	payment, err := service.Repository.byLinkHash(linkId)
	if err != nil {
		//todo: logger
		return paymentData{}, err
	}
	//if !strings.EqualFold(payment.status, PAYMENT_CREATED) {
	//	//todo: logger
	//	return paymentData{}, errors.New("Payment in wrong status")
	//}
	if payment.stripeIdPayment != "" {
		return paymentData{
			StripeConnectedAccountId: payment.stripeAccId,
			StripeClientSecret:       payment.stripeIdPayment,
		}, nil
	}
	currency := payprovider.AllCurrencies()[strings.ToLower(payment.currency)]

	amount, commission := currency.CalculatePayment(payment.amount, service.Commission)
	id, err := service.paymentProvider().CreatePayment(amount, commission, payment.currency, payment.description, payment.stripeAccId, payment.confirmedHash.String(), payment.canceledHash.String())

	if err != nil {
		//todo: logger
		return paymentData{}, err
	}
	err = service.Repository.statusChange(linkId, PAYMENT_INITIALED, id)
	if err != nil {
		//todo: logger
		return paymentData{}, err
	}

	return paymentData{
		StripeConnectedAccountId: payment.stripeAccId,
		StripeClientSecret:       id,
	}, nil

}
