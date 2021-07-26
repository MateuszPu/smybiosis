package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v71/account"
	"pay.me/v4/mail"
	"pay.me/v4/payment"
	"pay.me/v4/payprovider"
	"pay.me/v4/server"
)

type Service struct {
	BaseSever       *server.BaseSever
	Repository      *repository
	PaymentProvider *payprovider.PaymentProvider
	PaymentService  *payment.Service
	MailService     *mail.Service
}

func (service *Service) repository() repository {
	return *service.Repository
}

func (service *Service) paymentProvider() payprovider.PaymentProvider {
	return *service.PaymentProvider
}

func (service *Service) createUser(email string, usrAgent string) (*string, *uuid.UUID, error) {
	stripeAccId, err := service.paymentProvider().CreateUser(email)
	if err != nil {
		//todo:logger
		return nil, nil, err
	}

	usr, err := service.repository().create(email, stripeAccId, usrAgent)
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
	usr, err := service.repository().findByLinkId(linkId)
	if err != nil {
		return usr, err
	}
	if usr.status == STRIPE_CONFIRMED {
		return usr, errors.New("user is already registered in stripe")
	}

	go func(userId uuid.UUID) {
		service.BaseSever.Logger.Infof("Sending payment link for user %s", userId.String())
		service.PaymentService.GenerateFirstPaymentLink(userId)
	}(usr.id)
	return service.repository().updateUserStatus(linkId, STRIPE_CONFIRMED)
}

func (service *Service) findByLinkId(linkId string) (user, error) {
	return service.repository().findByLinkId(linkId)
}

func (service *Service) findByEmail(email string) (user, error) {
	return service.repository().findByEmail(email)
}

func (service *Service) sendCustomerEmail(from string, question string) {
	service.MailService.SendEmailFromCustomer(from, question)
}

func (service *Service) checkMissingParams(usr user) {
	acc, err := account.GetByID(usr.stripeId, nil)
	if err != nil {

	}
	if acc != nil {
		if acc.ChargesEnabled {
			//todo: send email with link
		} else {
			//todo: send email with info to fill profile
		}
		if !acc.PayoutsEnabled {
			//todo: send link with info
		}
		//todo: error handlong
		println("errr")
	}
}
