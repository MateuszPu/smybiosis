package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/account"
	"github.com/stripe/stripe-go/v71/accountlink"
	"time"
)

type Service struct {
	Repository Repository
}

func (service *Service) createUser(email string) (user, error) {
	stripeId, err := service.createUserInStripe(email)
	if err != nil {
		return user{}, err
	}
	createdUser := user{id: uuid.New(), email: email, stripeId: stripeId, createdAt: time.Now()}
	err = service.Repository.save(createdUser)
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

func (service *Service) stripeLink(stripeId string) (string, error) {
	par := &stripe.AccountLinkParams{
		Account:    stripe.String(stripeId),
		RefreshURL: stripe.String("https://example.com/reauth"), //TODO: change it to use env variables
		ReturnURL:  stripe.String("https://example.com/return"), //TODO: change it to use env variables
		Type:       stripe.String(string(stripe.AccountLinkTypeAccountOnboarding)),
	}

	newAcc, err := accountlink.New(par)
	if err != nil {
		return "", err
	}
	return newAcc.URL, nil
}

type Repository interface {
	save(user user) error
}

//for database service we will inject here mechanism to save in database
type RepositoryInMemory struct {
	inMemory map[string]user
}

func CreateInMemoryRepo() Repository {
	return &RepositoryInMemory{inMemory: make(map[string]user)}
}

func (service *RepositoryInMemory) save(user user) error {
	_, found := service.inMemory[user.email]
	if found {
		return errors.New("user already exist")
	}
	service.inMemory[user.email] = user
	return nil
}
