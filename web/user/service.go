package user

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/account"
	"github.com/stripe/stripe-go/v71/accountlink"
	"pay.me/v4/server"
	"strings"
)

type Service struct {
	Env        *server.Env
	Repository *repository
}

func (service *Service) repository() repository {
	return *service.Repository
}

func (service *Service) createUser(email string) (user, error) {
	stripeId, err := service.createUserInStripe(email)
	if err != nil {
		return user{}, err
	}
	createdUser := user{
		id: uuid.New(),
		email:     email,
		stripeId:  stripeId,
		linkId:    strings.ReplaceAll(uuid.New().String(), "-", "")}
	err = service.repository().save(createdUser)
	if err != nil {
		return user{}, err
	}
	return createdUser, nil
}

func (service *Service) createUserInStripe(email string) (string, error) {
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

func (service *Service) finishedStripeRegistration(linkId string) (user, error) {
	return service.repository().updateUserStatus(linkId, STRIPE_CONFIRMED)
}

func (service *Service) findByLinkId(linkId string) (user, error) {
	return service.repository().findByLinkId(linkId)
}

func (service *Service) stripeLink(stripeAccId string, linkId string) (string, error) {
	refreshUrl := fmt.Sprintf("%srefresh/%s", service.Env.Host, linkId)
	returnUrl := fmt.Sprintf("%sconfirm/%s", service.Env.Host, linkId)
	par := &stripe.AccountLinkParams{
		Account:    stripe.String(stripeAccId),
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

func (service *Service) findByEmail(email string) (user, error) {
	return service.repository().findByEmail(email)
}
