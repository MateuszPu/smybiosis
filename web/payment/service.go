package payment

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/checkout/session"
	"pay.me/v4/mail"
	"pay.me/v4/server"
)

type Service struct {
	Repository  repository
	MailService *mail.Service
	GlobalEnv   *server.Env
}

func (service *Service) CreatePayment(userId uuid.UUID, currency string, amount float64, description string) (uuid.UUID, error) {
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

	amount := int64(payment.amount * 100)
	commission := payment.amount * 100 * 0.005
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			string(stripe.PaymentMethodTypeCard),
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Name:     stripe.String(payment.description),
				Amount:   stripe.Int64(amount),
				Currency: stripe.String(string(stripe.CurrencyPLN)),
				Quantity: stripe.Int64(1),
			},
		},
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			ApplicationFeeAmount: stripe.Int64(int64(commission)),
		},
		//todo: implement two pages for this finished will mark payment as done
		//todo: canceled will ask user do you realy want cancel or you want to back there?
		SuccessURL: stripe.String(fmt.Sprintf("%spayment/finished/%s", service.GlobalEnv.Host, payment.confirmedHash)),
		CancelURL:  stripe.String(fmt.Sprintf("%spayment/canceled/%s", service.GlobalEnv.Host, payment.canceledHash)),
	}

	params.SetStripeAccount(payment.stripeAccId)
	s, err := session.New(params)
	if err != nil {
		//todo: logger
		return paymentData{}, err
	}
	err = service.Repository.statusChange(linkId, PAYMENT_INITIALED, s.ID)
	if err != nil {
		//todo: logger
		return paymentData{}, err
	}

	return paymentData{
		StripeConnectedAccountId: payment.stripeAccId,
		StripeClientSecret:       s.ID,
	}, nil

}
