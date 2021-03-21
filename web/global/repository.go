package global

import (
	"context"
	"database/sql"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	models "pay.me/v4/database-models"
)

type userDetails struct {
	Email string `boil:"email"`
	Amount float64 `boil:"amount"`
	Title string `boil:"title"`
	Currency string `boil:"currency"`
}

type  repository interface {
	findUserDetailsBy(userCookieId string) (userDetails, error)
}

type RepositorySql struct {
	database *sql.DB
}

func CreateSqlRepo(db *sql.DB) repository {
	return &RepositorySql{database: db}
}

func (repo *RepositorySql) findUserDetailsBy(userCookieId string) (userDetails, error) {
	//models.P
	var userDetails userDetails
	err := models.NewQuery(
		qm.Select("users.email as email", "payments.amount as amount", "payments.description as title", "payments.currency as currency"),
		qm.From("users"),
		qm.InnerJoin("payments on payments.user_id = users.id"),
		qm.Where("users.cookie_id = ?", userCookieId),
		qm.OrderBy("payments.created_at DESC"),
		qm.Limit(1),
		).Bind(context.Background(), repo.database, &userDetails)
	return userDetails, err
}
