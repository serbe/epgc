package epdc

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
)

// EDc struct to store *pg.DB
type EDc struct {
	db *pg.DB
}

// InitDB initialize database
func InitDB(dbname string, user string, password string, sslmode bool, logsql bool) (*EDc, error) {
	e := new(EDc)
	opt := pg.Options{
		User:     user,
		Password: password,
		Database: dbname,
		SSL:      sslmode,
	}
	if logsql == true {
		pg.SetQueryLogger(log.New(os.Stdout, "", log.LstdFlags))
	}
	e.db = pg.Connect(&opt)
	err := e.createTables()
	return e, err
}

func (e *EDc) createTables() (err error) {
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
