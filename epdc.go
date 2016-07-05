package epgc

import (
	"fmt"

	"database/sql"
	// need to sql dialect
	_ "github.com/lib/pq"
)

// Edb struct to store *DB
type Edb struct {
	db  *sql.DB
	log bool
}

// InitDB initialize database
func InitDB(dbname string, user string, password string, sslmode string, logsql bool) (*Edb, error) {
	e := new(Edb)
	// sslmode=disable
	opt := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", user, password, dbname, sslmode)
	db, err := sql.Open("postgres", opt)
	if err != nil {
		return e, err
	}
	err = db.Ping()
	if err != nil {
		return e, err
	}
	e.db = db
	e.log = logsql
	err = e.createTables()
	return e, err
}

func (e *Edb) createTables() (err error) {
	err = e.trainingCreateTable()
	if err != nil {
		return
	}
	err = e.kindCreateTable()
	if err != nil {
		return
	}
	err = e.emailCreateTable()
	if err != nil {
		return
	}
	err = e.companyCreateTable()
	if err != nil {
		return
	}
	err = e.peopleCreateTable()
	if err != nil {
		return
	}
	err = e.postCreateTable()
	if err != nil {
		return
	}
	err = e.rankCreateTable()
	if err != nil {
		return
	}
	err = e.scopeCreateTable()
	if err != nil {
		return
	}
	err = e.phoneCreateTable()
	if err != nil {
		return
	}
	err = e.practiceCreateTable()
	if err != nil {
		return
	}
	return
}
