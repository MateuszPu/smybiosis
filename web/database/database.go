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

	//dbURI := fmt.Sprintf("host=/cloudsql/%s user=%s password=%s dbname=%s sslmode=disable", d.Host, d.User, d.Password, d.Name)
	dbURI := fmt.Sprintf("dbname=%s host=%s user=%s password=%s sslmode=disable", d.Name, d.Host, d.User, d.Password)
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		//d.Logger.Error("Problem with open database connection %s", err.Error())
		//panic("problem with connect to database " + dbURI)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		//d.Logger.Error("Problem with database connection %s", pingErr.Error())
		//println("problem with connect to database " + dbURI)
		//panic("problem with database connection" + dbURI)
	}
	boil.SetDB(db)
	boil.DebugMode = false

	return db
}
