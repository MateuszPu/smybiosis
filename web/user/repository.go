package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	models "pay.me/v4/database-models"
)

type repository interface {
	save(user user) error
	findByEmail(email string) (user, error)
	findByLinkId(linkId string) (user, error)
	updateUserStatus(stripeId string, status string) (user, error)
}

type RepositorySql struct {
	database *sql.DB
}

func CreateSqlRepo(db *sql.DB) repository {
	return &RepositorySql{database: db}
}

func (repo *RepositorySql) save(user user) error {
	var dbUser = models.User{
		StripeAccount: user.stripeId,
		LinkRegistration:   user.linkId,
		Email:    user.email,
		Status:   ACCOUNT_CREATED,
		ID:       user.id.String(),
	}
	return dbUser.Insert(context.Background(), repo.database, boil.Infer())
}

func (repo *RepositorySql) findByLinkId(linkId string) (user, error) {
	dbUser, err := models.Users(models.UserWhere.LinkRegistration.EQ(linkId)).One(context.Background(), repo.database)
	if err != nil {
		return user{}, err
	}
	u := user{}
	return u.from(dbUser), nil
}

func (repo *RepositorySql) findByEmail(email string) (user, error) {
	dbUser, err := models.Users(qm.Where("LOWER(email)=?", email)).One(context.Background(), repo.database)
	if err != nil {
		return user{}, err
	}
	u := user{}
	return u.from(dbUser), nil
}

func (repo *RepositorySql) updateUserStatus(linkId string, status string) (user, error) {
	dbUser, err := models.Users(models.UserWhere.LinkRegistration.EQ(linkId)).One(context.Background(), repo.database)
	if err != nil {
		return user{}, err
	}
	dbUser.Status = status
	_, err = dbUser.Update(context.Background(), repo.database, boil.Whitelist("status"))
	u := user{}
	return u.from(dbUser), nil
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

func (repo *RepositoryInMemory) findByEmail(email string) (user, error) {
	usr, _ := repo.inMemory[email]
	return usr, nil
}

func (repo *RepositoryInMemory) updateUserStatus(linkId string, status string) (user, error) {
	for _, user := range repo.inMemory {
		if user.linkId == linkId {
			user.status = status
			repo.inMemory[user.email] = user
			return user, nil
		}
	}
	return user{}, errors.New("user not exist")
}
