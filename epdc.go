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

// SelectItem - struct for select element
type SelectItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
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
	err = e.createAllTables()
	return e, err
}

func (e *Edb) createAllTables() error {
	err := e.educationCreateTable()
	if err != nil {
		return err
	}
	err = e.kindCreateTable()
	if err != nil {
		return err
	}
	err = e.emailCreateTable()
	if err != nil {
		return err
	}
	err = e.companyCreateTable()
	if err != nil {
		return err
	}
	err = e.peopleCreateTable()
	if err != nil {
		return err
	}
	err = e.postCreateTable()
	if err != nil {
		return err
	}
	err = e.rankCreateTable()
	if err != nil {
		return err
	}
	err = e.scopeCreateTable()
	if err != nil {
		return err
	}
	err = e.phoneCreateTable()
	if err != nil {
		return err
	}
	err = e.practiceCreateTable()
	if err != nil {
		return err
	}
	return nil
}
