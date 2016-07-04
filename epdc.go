package epgc

import (
	"fmt"
	"strings"
	"time"

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
		return
	}
	err = db.Ping
	if err != nil {
		return
	}
	e.db = pg.Connect(&opt)
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

func toInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func int64InSlice(a int64, list []int64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func setStrMonth(d time.Time) (output string) {
	str := d.Format("02.01.2006")
	spl := strings.Split(str, ".")
	month := map[string]string{"01": "января", "02": "февраля", "03": "марта", "04": "апреля", "05": "мая", "06": "июня", "07": "июля", "08": "августа", "09": "сентября", "10": "октября", "11": "ноября", "12": "декабря "}
	output = spl[0] + " " + month[spl[1]] + " " + spl[2] + " года"
	return
}
