package payment

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	models "pay.me/v4/database-models"
)

type repository interface {
	save(payment, uuid.UUID) error
	findById(id uuid.UUID) (payment, error)
	findByLinkHash(id string) (payment, error)
	findPaymentByUserId(id uuid.UUID) (payment, error)
}

type RepositorySql struct {
	database *sql.DB
}

func CreateSqlRepo(db *sql.DB) repository {
	return &RepositorySql{database: db}
}

func (repo *RepositorySql) save(payment payment, userId uuid.UUID) error {
	dbPayment := models.Payment{
		ID:            payment.id.String(),
		LinkHash:      uuid.New().String(),
		ConfirmedHash: uuid.New().String(),
		CanceledHash:  uuid.New().String(),
		UserID:        userId.String(),
		Currency:      payment.currency,
		Amount:        payment.amount,
		Description:   payment.description,
		Status:        PAYMENT_CREATED,
	}
	return dbPayment.Insert(context.Background(), repo.database, boil.Infer())
}

func (repo *RepositorySql) findById(id uuid.UUID) (payment, error) {
	query := models.PaymentWhere.ID.EQ(id.String())
	return repo.findPaymentBy(query)
}

func (repo *RepositorySql) findByLinkHash(id string) (payment, error) {
	query := models.PaymentWhere.LinkHash.EQ(id)
	return repo.findPaymentBy(query)
}

func (repo *RepositorySql) findPaymentByUserId(userId uuid.UUID) (payment, error) {
	query := models.PaymentWhere.UserID.EQ(userId.String())
	return repo.findPaymentBy(query)
}

func (repo *RepositorySql) findPaymentBy(query qm.QueryMod) (payment, error) {
	dbPayment, err := models.Payments(query).One(context.Background(), repo.database)
	if err != nil {
		return payment{}, err
	}
	one, err := dbPayment.User().One(context.Background(), repo.database)
	if err != nil {
		return payment{}, err
	}
	payment := payment{}
	return payment.from(dbPayment, one), nil
}
