package user

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	models "pay.me/v4/database-models"
)

type repository interface {
	create(email string, stripeId string) (user, error)
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

func (repo *RepositorySql) create(email string, stripeId string) (user, error) {
	var dbUser = models.User{
		StripeAccount:    stripeId,
		LinkRegistration: uuid.New().String(),
		Email:            email,
		Status:           ACCOUNT_CREATED,
		ID:               uuid.New().String(),
	}
	err := dbUser.Insert(context.Background(), repo.database, boil.Infer())
	if err != nil {
		return user{}, err
	}
	u := user{}
	return u.from(&dbUser), nil
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
