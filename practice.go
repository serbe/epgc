package epgc

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
)

// Practice - struct for practice
type Practice struct {
	ID             int64     `sql:"id" json:"id"`
	Company        Company   `sql:"-"`
	CompanyID      int64     `sql:"company_id, null" json:"company-id"`
	Kind           Kind      `sql:"-"`
	KindID         int64     `sql:"kind_id, null" json:"kind-id"`
	Topic          string    `sql:"topic, null" json:"topic"`
	DateOfPractice time.Time `sql:"date_of_practice, null" json:"date-of-practice"`
	DateStr        string    `sql:"-" json:"date-str"`
	Note           string    `sql:"note, null" json:"note"`
	CreatedAt      time.Time `sql:"created_at" json:"created_at"`
	UpdatedAt      time.Time `sql:"updated_at" json:"updated_at"`
}

func scanPractice(row *sql.Row) (Practice, error) {
	var (
		sID             sql.NullInt64
		sCompanyID      sql.NullInt64
		sKindID         sql.NullInt64
		sTopic          sql.NullString
		sDateOfPractice pq.NullTime
		sNote           sql.NullString
		practice        Practice
	)
	err := row.Scan(&sID, &sCompanyID, &sKindID, &sTopic, &sDateOfPractice, &sNote)
	if err != nil {
		log.Println("scanPractice row.Scan ", err)
		return practice, err
	}
	practice.ID = n2i(sID)
	practice.CompanyID = n2i(sCompanyID)
	practice.KindID = n2i(sKindID)
	practice.Topic = n2s(sTopic)
	practice.DateOfPractice = n2d(sDateOfPractice)
	practice.Note = n2s(sNote)
	return practice, nil
}

func scanPractices(rows *sql.Rows, opt string) ([]Practice, error) {
	var practices []Practice
	for rows.Next() {
		var (
			sID sql.NullInt64
			// sCompanyID      sql.NullInt64
			sCompanyName sql.NullString
			// sKindID         sql.NullInt64
			sKindName       sql.NullString
			sTopic          sql.NullString
			sDateOfPractice pq.NullTime
			sNote           sql.NullString
			practice        Practice
		)
		switch opt {
		case "list":
			err := rows.Scan(&sID, &sCompanyName, &sKindName, &sTopic, &sDateOfPractice, &sNote)
			if err != nil {
				log.Println("scanPractices rows.Scan list ", err)
				return practices, err
			}
			practice.Company.Name = n2s(sCompanyName)
			practice.Kind.Name = n2s(sKindName)
			practice.Topic = n2s(sTopic)
			practice.Note = n2s(sNote)
		case "company":
			err := rows.Scan(&sID, &sKindName, &sTopic, &sDateOfPractice)
			if err != nil {
				log.Println("scanPractices rows.Scan company ", err)
				return practices, err
			}
			practice.Kind.Name = n2s(sKindName)
			// if len(practice.Kind.Name) > 210 {
			// 	practice.Kind.Name = practice.Kind.Name[0:210]
			// }
			practice.Topic = n2s(sTopic)
			// if len(practice.Topic) > 210 {
			// 	practice.Topic = practice.Topic[0:210]
			// }
		}
		practice.ID = n2i(sID)
		practice.DateOfPractice = n2d(sDateOfPractice)
		practice.DateStr = setStrMonth(practice.DateOfPractice)
		practices = append(practices, practice)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanPractices rows.Err ", err)
	}
	return practices, err
}

// GetPractice - get one practice by id
func (e *Edb) GetPractice(id int64) (Practice, error) {
	if id == 0 {
		return Practice{}, nil
	}
	stmt, err := e.db.Prepare(`SELECT
		id,
		company_id,
		kind_id,
		topic,
		date_of_practice,
		note
	FROM
		practices
	WHERE id = $1`)
	if err != nil {
		log.Println("GetPractice e.db.Prepare ", err)
		return Practice{}, err
	}
	row := stmt.QueryRow(id)
	practice, err := scanPractice(row)
	return practice, err
}

// GetPracticeList - get all practices for list
func (e *Edb) GetPracticeList() ([]Practice, error) {
	rows, err := e.db.Query(`SELECT
		p.id,
		c.name AS company_name,
		k.name AS kind_name,
		p.topic,
		p.date_of_practice,
		p.note
	FROM
		practices AS p
	LEFT JOIN
		companies AS c ON c.id = p.company_id
	LEFT JOIN
		kinds AS k ON k.id = p.kind_id
	ORDER BY
		date_of_practice DESC`)
	if err != nil {
		log.Println("GetPracticeAll e.db.Query ", err)
		return []Practice{}, err
	}
	practices, err := scanPractices(rows, "list")
	return practices, err
}

// GetPracticeCompany - get all practices of company
func (e *Edb) GetPracticeCompany(id int64) ([]Practice, error) {
	if id == 0 {
		return []Practice{}, nil
	}
	stmt, err := e.db.Prepare(`SELECT
		p.id,
		k.name AS kind_name,
		p.topic,
		p.date_of_practice
	FROM
		practices AS p
	LEFT JOIN
		kinds AS k ON k.id = p.kind_id
	WHERE
	    p.company_id = $1
	ORDER BY
		date_of_practice`)
	if err != nil {
		log.Println("GetPracticeCompany e.db.Prepare ", err)
		return []Practice{}, err
	}
	rows, err := stmt.Query(id)
	if err != nil {
		log.Println("GetPracticeCompany stmt.Query ", err)
		return []Practice{}, err
	}
	practices, err := scanPractices(rows, "company")
	return practices, err
}

// CreatePractice - create new practice
func (e *Edb) CreatePractice(practice Practice) (int64, error) {
	stmt, err := e.db.Prepare(`INSERT INTO practices(company_id, kind_id, topic, date_of_practice, note, created_at) VALUES($1, $2, $3, $4, $5, now()) RETURNING id`)
	if err != nil {
		log.Println("CreatePractice e.db.Prepare ", err)
		return 0, err
	}
	err = stmt.QueryRow(i2n(practice.CompanyID), i2n(practice.KindID), s2n(practice.Topic), d2n(practice.DateOfPractice), s2n(practice.Note)).Scan(&practice.ID)
	return practice.ID, err
}

// UpdatePractice - save practice changes
func (e *Edb) UpdatePractice(practice Practice) error {
	stmt, err := e.db.Prepare(`UPDATE practices SET company_id=$2, kind_id=$3, topic=$4, date_of_practice=$5, note=$6, updated_at = now() WHERE id=$1`)
	if err != nil {
		log.Println("UpdatePractice e.db.Prepare ", err)
		return err
	}
	_, err = stmt.Exec(practice.ID, i2n(practice.CompanyID), i2n(practice.KindID), s2n(practice.Topic), d2n(practice.DateOfPractice), s2n(practice.Note))
	return err
}

// DeletePractice - delete practice by id
func (e *Edb) DeletePractice(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`DELETE FROM practices WHERE id = $1`, id)
	if err != nil {
		log.Println("DeletePractice e.db.Exec: ", id, err)
		return fmt.Errorf("DeletePractice e.db.Exec: %s", err)
	}
	return err
}

func (e *Edb) practiceCreateTable() error {
	str := `CREATE TABLE IF NOT EXISTS practices (id bigserial primary key, company_id bigint, kind_id bigint, topic text, date_of_practice date, note text, created_at TIMESTAMP without time zone, updated_at TIMESTAMP without time zone)`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("practiceCreateTable e.db.Exec: ", err)
		return fmt.Errorf("practiceCreateTable e.db.Exec: %s", err)
	}
	return err
}
