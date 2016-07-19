package epgc

import (
	"database/sql"
	"log"

	"github.com/lib/pq"
)

// Education - struct for education
type Education struct {
	ID        int64  `sql:"id" json:"id" `
	StartDate string `sql:"start_date" json:"start_date"`
	EndDate   string `sql:"end_date" json:"end_date"`
	StartStr  string `sql:"-" json:"start_str"`
	EndStr    string `sql:"-" json:"end_str"`
	Note      string `sql:"note, null" json:"note"`
	CreatedAt string `sql:"created_at" json:"created_at"`
	UpdatedAt string `sql:"updated_at" json:"updated_at"`
}

func scanEducation(row *sql.Row) (Education, error) {
	var (
		sID        sql.NullInt64
		sStartDate pq.NullTime
		sEndDate   pq.NullTime
		sNote      sql.NullString
		education  Education
	)
	err := row.Scan(&sID, &sStartDate, &sEndDate, &sNote)
	if err != nil {
		log.Println("scanEducation row.Scan ", err)
		return education, err
	}
	education.ID = n2i(sID)
	education.StartDate = n2sd(sStartDate)
	education.EndDate = n2sd(sEndDate)
	education.Note = n2s(sNote)
	return education, nil
}

func scanEducations(rows *sql.Rows, opt string) ([]Education, error) {
	var educations []Education
	for rows.Next() {
		var (
			sID        sql.NullInt64
			sStartDate pq.NullTime
			sEndDate   pq.NullTime
			sNote      sql.NullString
			education  Education
		)
		switch opt {
		case "list":
			err := rows.Scan(&sID, &sStartDate, &sEndDate, &sNote)
			if err != nil {
				log.Println("scanEducations rows.Scan list ", err)
				return educations, err
			}
		}
		education.ID = n2i(sID)
		education.StartDate = n2sd(sStartDate)
		education.EndDate = n2sd(sEndDate)
		education.Note = n2s(sNote)
		educations = append(educations, education)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanEducations rows.Err ", err)
	}
	return educations, err
}

// GetEducation - get education by id
func (e *Edb) GetEducation(id int64) (Education, error) {
	if id == 0 {
		return Education{}, nil
	}
	stmt, err := e.db.Prepare(`SELECT
		id,
		start_date,
		end_date,
		note
	FROM
		educations
	WHERE id = $1
	ORDER BY
		start_date`)
	if err != nil {
		log.Println("GetEducation e.db.Prepare ", err)
		return Education{}, err
	}
	row := stmt.QueryRow(id)
	education, err := scanEducation(row)
	return education, err
}

// GetEducationList - get all education for list
func (e *Edb) GetEducationList() ([]Education, error) {
	rows, err := e.db.Query(`SELECT
		id,
		start_date,
		end_date,
		note
	FROM
		educations
	ORDER BY
		start_date`)
	if err != nil {
		log.Println("GetEducationList e.db.Query ", err)
		return []Education{}, err
	}
	educations, err := scanEducations(rows, "list")
	if err != nil {
		log.Println("GetEducationList scanEducations ", err)
		return []Education{}, err
	}
	for i := range educations {
		educations[i].StartStr = setStrMonth(educations[i].StartDate)
		educations[i].EndStr = setStrMonth(educations[i].EndDate)
	}
	return educations, err
}

// CreateEducation - create new education
func (e *Edb) CreateEducation(education Education) (int64, error) {
	stmt, err := e.db.Prepare(`INSERT INTO educations(start_date, end_date, note, created_at) VALUES($1, $2, $3, now()) RETURNING id`)
	if err != nil {
		log.Println("CreateEducation e.db.Prepare ", err)
		return 0, err
	}
	err = stmt.QueryRow(sd2n(education.StartDate), sd2n(education.EndDate), s2n(education.Note)).Scan(&education.ID)
	return education.ID, err
}

// UpdateEducation - save changes to education
func (e *Edb) UpdateEducation(education Education) error {
	stmt, err := e.db.Prepare(`UPDATE educations SET start_date = $2, end_date = $3, note = $4, updated_at = now() WHERE id = $1`)
	if err != nil {
		log.Println("UpdateEducation e.db.Prepare ", err)
		return err
	}
	_, err = stmt.Exec(education.ID, sd2n(education.StartDate), sd2n(education.EndDate), s2n(education.Note))
	return err
}

// DeleteEducation - delete education by id
func (e *Edb) DeleteEducation(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`DELETE FROM educations WHERE id = $1`, id)
	if err != nil {
		log.Println("DeleteEducation ", id, err)
	}
	return err
}

func (e *Edb) educationCreateTable() error {
	str := `CREATE TABLE IF NOT EXISTS educations (id bigserial primary key, start_date date, end_date date, note text, created_at TIMESTAMP without time zone, updated_at TIMESTAMP without time zone)`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("educationCreateTable ", err)
	}
	return err
}
