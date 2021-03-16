package payprovider

import (
	"fmt"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/account"
	"github.com/stripe/stripe-go/v71/accountlink"
	"github.com/stripe/stripe-go/v71/checkout/session"
	"pay.me/v4/server"
	"strings"
)

type PaymentProvider interface {
	CreateUser(email string) (string, error)
	RegistrationLink(externalAccId string, linkId string) (string, error)
	CreatePayment(amount int64, commission int64, currency string, description string, stripeAccId string, confirmedHash string, canceledHash string) (string, error)
}

type StripeProvider struct {
	Env *server.Env
}

func (service StripeProvider) Init() PaymentProvider {
	stripe.Key = service.Env.StripeKey
	return service
}

func (service StripeProvider) CreateUser(email string) (string, error) {
	params := &stripe.AccountParams{
		Type:  stripe.String(string(stripe.AccountTypeStandard)),
		Email: stripe.String(email),
	}
	acc, err := account.New(params)
	if err != nil {
		return "", err
	}
	return acc.ID, nil
}

func (service StripeProvider) RegistrationLink(externalAccId string, linkId string) (string, error) {
	refreshUrl := fmt.Sprintf("%srefresh/%s", service.Env.Host, linkId)
	returnUrl := fmt.Sprintf("%sconfirm/%s", service.Env.Host, linkId)
	par := &stripe.AccountLinkParams{
		Account:    stripe.String(externalAccId),
		RefreshURL: stripe.String(refreshUrl),
		ReturnURL:  stripe.String(returnUrl),
		Type:       stripe.String(string(stripe.AccountLinkTypeAccountOnboarding)),
	}

	newAcc, err := accountlink.New(par)
	if err != nil {
		return "", err
	}
	return newAcc.URL, nil
}

func (service StripeProvider) CreatePayment(amount int64, commission int64, currency string, description string, stripeAccId string, successHash string, canceledHash string) (string, error) {
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			string(stripe.PaymentMethodTypeCard),
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Name:     stripe.String(description),
				Amount:   stripe.Int64(amount),
				Currency: stripe.String(strings.ToLower(currency)),
				Quantity: stripe.Int64(1),
			},
		},
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			ApplicationFeeAmount: stripe.Int64(int64(commission)),
		},
		//todo: implement two pages for this finished will mark payment as done
		//todo: canceled will ask user do you realy want cancel or you want to back there?
		//todo; canceled also appears when user re-open the payment link
		SuccessURL: stripe.String(fmt.Sprintf("%spayments/%s/success", service.Env.Host, successHash)),
		CancelURL:  stripe.String(fmt.Sprintf("%spayments/%s/cancel", service.Env.Host, canceledHash)),
	}

	params.SetStripeAccount(stripeAccId)
	s, err := session.New(params)
	if err != nil {
		return "", err
	}
	return s.ID, nil
}
