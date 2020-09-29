package user

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/account"
	"github.com/stripe/stripe-go/v71/accountlink"
	"pay.me/server"
	"strings"
	"time"
)

type Service struct {
	Env        *server.Env
	Repository repository
}

func (service *Service) createUser(email string) (user, error) {
	stripeId, err := service.createUserInStripe(email)
	if err != nil {
		return user{}, err
	}
	createdUser := user{id: uuid.New(),
		email:     email,
		stripeId:  stripeId,
		linkId:    strings.ReplaceAll(uuid.New().String(), "-", ""),
		status:    "ACCOUNT_CREATED",
		createdAt: time.Now()}
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

func (service *Service) finishedStripeRegistration(linkId string) (user, error) {
	service.Repository.updateUserStatus(linkId, "STRIPE_ACCOUNT_CREATED")
	return service.Repository.findByLinkId(linkId)
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

func (service *Service) findByEmail(email string) user {
	return service.Repository.findByEmail(email)
}

type repository interface {
	save(user user) error
	findByEmail(email string) user
	findByLinkId(linkId string) (user, error)
	updateUserStatus(stripeId string, status string) error
}

//for database service we will inject here mechanism to save in database
type RepositoryInMemory struct {
	inMemory map[string]user
}

func CreateInMemoryRepo() repository {
	return &RepositoryInMemory{inMemory: make(map[string]user)}
}

func (repo *RepositoryInMemory) findByLinkId(linkId string) (user, error) {
	for _, user := range repo.inMemory {
		if user.linkId == linkId {
			return user, nil
		}
	}
	return user{}, errors.New("user does not exist")
}

func (repo *RepositoryInMemory) save(user user) error {
	_, found := repo.inMemory[user.email]
	if found {
		return errors.New("user already exist")
	}
	repo.inMemory[user.email] = user
	return nil
}

func (repo *RepositoryInMemory) findByEmail(email string) user {
	usr, _ := repo.inMemory[email]
	return usr
}

func (repo *RepositoryInMemory) updateUserStatus(linkId string, status string) error {
	for _, user := range repo.inMemory {
		if user.linkId == linkId {
			user.status = status
			repo.inMemory[user.email] = user
		}
	}
	return errors.New("user not exist")
}
