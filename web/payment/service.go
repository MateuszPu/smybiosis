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
	Repository  *repository
	MailService *mail.Service
	GlobalEnv   *server.Env
}

func (service *Service) repository() repository {
	return *service.Repository
}

func (service *Service) CreatePayment(userId uuid.UUID, currency string, amount float64, description string) (uuid.UUID, error) {
	payment := payment{
		id:                uuid.New(),
		currency:          currency,
		amount:            amount,
		description:       description,
	}
	err := service.repository().save(payment, userId)
	return payment.id, err
}

func (service *Service) GenerateFirstPaymentLink(userId uuid.UUID) {
	payment, err := service.repository().findPaymentByUserId(userId)
	if err != nil {
		//todo; log here
		return
	}
	service.MailService.SendEmailWithPaymentLink(payment.email, fmt.Sprintf("%spayment/%s", service.GlobalEnv.Host, payment.linkHash.String()))
}


func (service *Service) GeneratePaymentLink(id uuid.UUID) {
	payment, err := service.repository().findById(id)
	if err != nil {
		//todo; log here
		return
	}
	service.MailService.SendEmailWithPaymentLink(payment.email, fmt.Sprintf("%spayment/%s", service.GlobalEnv.Host, payment.linkHash.String()))
}

func (service *Service) createStripePayment(linkId string) paymentData {
	//TODO: error handling?
	payment, _ := service.repository().findByLinkHash(linkId)

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
	//TODO: error handling?
	s, _ := session.New(params)

	return paymentData{
		StripeConnectedAccountId: payment.stripeAccId,
		StripeClientSecret:       s.ID,
	}

}

