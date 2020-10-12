package database

import (
	"database/sql"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"pay.me/v4/logging"
)

type SqlDatabase struct {
	Logger     *logging.StandardLogger
	DriverName string
	Name       string
	Host       string
	User       string
	Password   string
}

func (d SqlDatabase) CreateDb() *sql.DB {
	db, err := sql.Open(d.DriverName, fmt.Sprintf("dbname=%s host=%s user=%s password=%s sslmode=disable", d.Name, d.Host, d.User, d.Password))
	if err != nil {
		d.Logger.Error("Problem with open database connection %s", err.Error())
		panic("problem with connect to database")
	}
	pingErr := db.Ping()
	if pingErr != nil {
		d.Logger.Error("Problem with database connection %s", pingErr.Error())
		panic("problem with database connection")
	}
	boil.SetDB(db)
	boil.DebugMode = true

	return db
}
