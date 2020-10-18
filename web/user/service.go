package user

import (
	"github.com/google/uuid"
	"pay.me/v4/payprovider"
	"pay.me/v4/server"
)

type Service struct {
	Env        *server.Env
	Repository *repository
	PaymentProvider *payprovider.PaymentProvider
}

func (service *Service) repository() repository {
	return *service.Repository
}

func (service *Service) paymentProvider() payprovider.PaymentProvider {
	return *service.PaymentProvider
}

func (service *Service) createUser(email string) (*string, *uuid.UUID, error) {
	stripeAccId, err := service.paymentProvider().CreateUser(email)
	if err != nil {
		//todo:logger
		return nil, nil, err
	}

	usr, err := service.repository().create(email, stripeAccId)
	if err != nil {
		//todo:logger
		return nil, nil, err
	}

	link, err := service.registrationLink(stripeAccId, usr.linkId)
	if err != nil {
		//todo:logger
		return nil, nil, err
	}

	return &link, &usr.id, nil
}

func (service *Service) registrationLink(stripeAccId string, linkId string) (string, error) {
	return service.paymentProvider().RegistrationLink(stripeAccId, linkId)
}

func (service *Service) finishedStripeRegistration(linkId string) (user, error) {
	return service.repository().updateUserStatus(linkId, STRIPE_CONFIRMED)
}

func (service *Service) findByLinkId(linkId string) (user, error) {
	return service.repository().findByLinkId(linkId)
}

func (service *Service) findByEmail(email string) (user, error) {
	return service.repository().findByEmail(email)
}
