package epgc

import (
	"database/sql"
	"log"
	"time"
)

// Email - struct for email
type Email struct {
	ID        int64     `sql:"id" json:"id"`
	CompanyID int64     `sql:"company_id, pk, null" json:"company-id"`
	PeopleID  int64     `sql:"people_id, pk, null" json:"people-id"`
	Email     string    `sql:"email, null" json:"email"`
	CreatedAt time.Time `sql:"created_at" json:"created_at"`
	UpdatedAt time.Time `sql:"updated_at" json:"updated_at"`
}

func scanEmail(row *sql.Row) (Email, error) {
	var (
		sid        sql.NullInt64
		scompanyID sql.NullInt64
		speopleID  sql.NullInt64
		semail     sql.NullString
		email      Email
	)
	err := row.Scan(&sid, &scompanyID, &speopleID, &semail)
	if err != nil {
		log.Println("scanEmail row.Scan ", err)
		return email, err
	}
	email.ID = n2i(sid)
	email.CompanyID = n2i(scompanyID)
	email.PeopleID = n2i(speopleID)
	email.Email = n2s(semail)
	return email, nil
}

func scanEmails(rows *sql.Rows, opt string) ([]Email, error) {
	var emails []Email
	for rows.Next() {
		var (
			sid    sql.NullInt64
			semail sql.NullString
			email  Email
		)
		switch opt {
		case "list":
			err := rows.Scan(&sid, &semail)
			if err != nil {
				log.Println("scanEmails rows.Scan ", err)
				return emails, err
			}
			email.Email = n2s(semail)
		case "select":
			err := rows.Scan(&sid, &semail)
			if err != nil {
				log.Println("scanEmails rows.Scan ", err)
				return emails, err
			}
			email.Email = n2s(semail)
			// if len(email.Email) > 210 {
			// 	email.Email = email.Email[0:210]
			// }
		}
		email.ID = n2i(sid)
		emails = append(emails, email)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanEmails rows.Err ", err)
	}
	return emails, err
}

// GetEmail - get one email by id
func (e *Edb) GetEmail(id int64) (Email, error) {
	if id == 0 {
		return Email{}, nil
	}
	stmt, err := e.db.Prepare(`SELECT id, company_id, people_id, email FROM emails WHERE id = $1`)
	if err != nil {
		log.Println("GetEmail e.db.Prepare", err)
		return Email{}, err
	}
	row := stmt.QueryRow(id)
	email, err := scanEmail(row)
	return email, nil
}

// GetEmailList - get all emails for list
func (e *Edb) GetEmailList() ([]Email, error) {
	rows, err := e.db.Query(`SELECT id, email FROM emails ORDER BY name ASC`)
	if err != nil {
		log.Println("GetEmailList e.db.Query ", err)
		return []Email{}, err
	}
	emails, err := scanEmails(rows, "list")
	return emails, err
}

// GetCompanyEmails - get all emails by company id
func (e *Edb) GetCompanyEmails(id int64) ([]Email, error) {
	if id == 0 {
		return []Email{}, nil
	}
	rows, err := e.db.Query(`SELECT id, email FROM emails ORDER BY name ASC WHERE company_id = $1`, id)
	if err != nil {
		log.Println("GetCompanyEmails e.db.Query ", err)
		return []Email{}, err
	}
	emails, err := scanEmails(rows, "list")
	return emails, err
}

// GetPeopleEmails - get all emails by people id
func (e *Edb) GetPeopleEmails(id int64) ([]Email, error) {
	if id == 0 {
		return []Email{}, nil
	}
	rows, err := e.db.Query(`SELECT id, email FROM emails ORDER BY name ASC WHERE people_id = $1`, id)
	if err != nil {
		log.Println("GetPeopleEmails e.db.Query ", err)
		return []Email{}, err
	}
	emails, err := scanEmails(rows, "list")
	return emails, err
}

// CreateEmail - create new email
func (e *Edb) CreateEmail(email Email) (int64, error) {
	stmt, err := e.db.Prepare(`INSERT INTO emails(company_id, people_id, email, created_at) VALUES($1, $2, $3, now()) RETURNING id`)
	if err != nil {
		log.Println("CreateEmail e.db.Prepare ", err)
		return 0, err
	}
	err = stmt.QueryRow(i2n(email.CompanyID), i2n(email.PeopleID), s2n(email.Email)).Scan(&email.ID)
	if err != nil {
		log.Println("CreateEmail db.QueryRow ", err)
		return 0, err
	}
	return email.ID, nil
}

// CreateCompanyEmails - create new company email
func (e *Edb) CreateCompanyEmails(company Company) error {
	err := e.DeleteCompanyEmails(company.ID)
	if err != nil {
		log.Println("CreateCompanyEmails DeleteCompanyEmails ", err)
		return err
	}
	for _, email := range company.Emails {
		email.CompanyID = company.ID
		_, err = e.CreateEmail(email)
		if err != nil {
			log.Println("CreateCompanyEmails CreateEmail ", err)
			return err
		}
	}
	return nil
}

// CreatePeopleEmails - create new people email
func (e *Edb) CreatePeopleEmails(people People) error {
	err := e.DeletePeopleEmails(people.ID)
	if err != nil {
		log.Println("CreatePeopleEmails DeletePeopleEmails ", err)
		return err
	}
	for _, email := range people.Emails {
		email.PeopleID = people.ID
		_, err = e.CreateEmail(email)
		if err != nil {
			log.Println("CreatePeopleEmails CreateEmail ", err)
			return err
		}
	}
	return nil
}

// UpdateEmail - save email changes
func (e *Edb) UpdateEmail(email Email) error {
	stmt, err := e.db.Prepare(`UPDATE emails SET company_id=$2, people_id=$3, email=$4, updated_at = now() WHERE id=$1`)
	if err != nil {
		log.Println("UpdateEmail e.db.Prepare ", err)
		return err
	}
	_, err = stmt.Exec(i2n(email.ID), i2n(email.CompanyID), i2n(email.PeopleID), s2n(email.Email))
	if err != nil {
		log.Println("UpdateEmail stmt.Exec ", err)
	}
	return err
}

// DeleteEmail - delete email by id
func (e *Edb) DeleteEmail(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`DELETE FROM emails WHERE id = $1`, id)
	if err != nil {
		log.Println("DeleteEmail ", err)
	}
	return err
}

// DeleteCompanyEmails - delete all emails by company id
func (e *Edb) DeleteCompanyEmails(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`DELETE FROM emails WHERE company_id = $1`, id)
	if err != nil {
		log.Println("DeleteCompanyEmails ", err)
	}
	return err
}

// DeletePeopleEmails - delete all emails by people id
func (e *Edb) DeletePeopleEmails(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`DELETE FROM emails WHERE people_id = $1`, id)
	if err != nil {
		log.Println("DeletePeopleEmails ", err)
	}
	return err
}

func (e *Edb) emailCreateTable() error {
	str := `CREATE TABLE IF NOT EXISTS emails (id bigserial primary key, company_id bigint, people_id bigint, email text, created_at timestamp without time zone, updated_at timestamp without time zone)`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("emailCreateTable ", err)
	}
	return err
}
