package user

import (
	"github.com/google/uuid"
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

func (service *Service) createUser(email string, stripeId string) (user, error) {
	createdUser := user{
		id: uuid.New(),
		email:     email,
		stripeId:  stripeId,
		linkId:    strings.ReplaceAll(uuid.New().String(), "-", "")}
	err := service.repository().save(createdUser)
	if err != nil {
		return user{}, err
	}
	return createdUser, nil
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
