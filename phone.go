package epgc

import (
	"database/sql"
	"log"
	"time"
)

// Phone - struct for phone
type Phone struct {
	ID        int64     `sql:"id" json:"id"`
	CompanyID int64     `sql:"company_id, pk, null" json:"company-id"`
	PeopleID  int64     `sql:"people_id, pk, null" json:"people-id"`
	Phone     int64     `sql:"phone, null" json:"phone"`
	Fax       bool      `sql:"fax, null" json:"fax"`
	CreatedAt time.Time `sql:"created_at" json:"created_at"`
	UpdatedAt time.Time `sql:"updated_at" json:"updated_at"`
}

func scanPhone(row *sql.Row) (Phone, error) {
	var (
		sID        sql.NullInt64
		sCompanyID sql.NullInt64
		sPeopleID  sql.NullInt64
		sPhone     sql.NullInt64
		sFax       sql.NullBool
		phone      Phone
	)
	err := row.Scan(&sID, &sCompanyID, &sPeopleID, &sPhone, &sFax)
	if err != nil {
		log.Println("scanPhone row.Scan ", err)
		return phone, err
	}
	phone.ID = n2i(sID)
	phone.CompanyID = n2i(sCompanyID)
	phone.PeopleID = n2i(sPeopleID)
	phone.Phone = n2i(sPhone)
	phone.Fax = n2b(sFax)
	return phone, nil
}

func scanPhones(rows *sql.Rows, opt string) ([]Phone, error) {
	var phones []Phone
	for rows.Next() {
		var (
			sID        sql.NullInt64
			sCompanyID sql.NullInt64
			sPeopleID  sql.NullInt64
			sPhone     sql.NullInt64
			sFax       sql.NullBool
			phone      Phone
		)
		switch opt {
		case "list":
			err := rows.Scan(&sID, &sCompanyID, &sPeopleID, &sPhone, &sFax)
			if err != nil {
				log.Println("scanPhones rows.Scan list ", err)
				return phones, err
			}
			phone.CompanyID = n2i(sCompanyID)
			phone.PeopleID = n2i(sPeopleID)
			phone.Fax = n2b(sFax)
		case "short":
			err := rows.Scan(&sID, &sPhone)
			if err != nil {
				log.Println("scanPhones rows.Scan short ", err)
				return phones, err
			}
		}
		phone.ID = n2i(sID)
		phone.Phone = n2i(sPhone)
		phones = append(phones, phone)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanPhones rows.Err ", err)
	}
	return phones, err
}

// GetPhone - get one phone by id
func (e *Edb) GetPhone(id int64) (Phone, error) {
	if id == 0 {
		return Phone{}, nil
	}
	stmt, err := e.db.Prepare(`SELECT id, company_id, people_id, phone, fax FROM phones WHERE id = $1`)
	if err != nil {
		log.Println("GetPhone e.db.Prepare", err)
		return Phone{}, err
	}
	row := stmt.QueryRow(id)
	phone, err := scanPhone(row)
	return phone, nil
}

// GetPhoneList - get all phones for list
func (e *Edb) GetPhoneList() ([]Phone, error) {
	rows, err := e.db.Query(`SELECT id, company_id, people_id, phone, fax FROM phones ORDER BY phone ASC`)
	if err != nil {
		log.Println("GetPhoneList e.db.Query ", err)
		return []Phone{}, err
	}
	phones, err := scanPhones(rows, "list")
	return phones, err
}

// GetCompanyPhones - get all phones by company id
func (e *Edb) GetCompanyPhones(id int64, fax bool) ([]Phone, error) {
	if id == 0 {
		return []Phone{}, nil
	}
	rows, err := e.db.Query(`SELECT id, phone FROM phones WHERE company_id = $1 AND fax = $2 ORDER BY phone ASC`, id, fax)
	if err != nil {
		log.Println("GetCompanyPhones e.db.Query ", err)
		return []Phone{}, err
	}
	phones, err := scanPhones(rows, "short")
	return phones, err
}

// GetPeoplePhones - get all phones by people id
func (e *Edb) GetPeoplePhones(id int64, fax bool) ([]Phone, error) {
	if id == 0 {
		return []Phone{}, nil
	}
	rows, err := e.db.Query(`SELECT id, phone FROM phones ORDER BY phone ASC WHERE people_id = $1 AND fax = $2`, id, fax)
	if err != nil {
		log.Println("GetPeoplePhones e.db.Query ", err)
		return []Phone{}, err
	}
	phones, err := scanPhones(rows, "short")
	return phones, err
}

// GetCompanyPhonesAll - get all faxes or phones by company id and isfax
func (e *Edb) GetCompanyPhonesAll(id int64, fax bool) ([]Phone, error) {
	if id == 0 {
		return []Phone{}, nil
	}
	rows, err := e.db.Query(`SELECT id, phone FROM phones WHERE company_id = $1 and fax = $2 ORDER BY phone ASC`, id, fax)
	if err != nil {
		log.Println("GetCompanyPhonesAll e.db.Query ", err)
		return []Phone{}, nil
	}
	phones, err := scanPhones(rows, "short")
	return phones, err
}

// GetPeoplePhonesAll - get all phones and faxes by people id
func (e *Edb) GetPeoplePhonesAll(id int64, fax bool) ([]Phone, error) {
	if id == 0 {
		return []Phone{}, nil
	}
	rows, err := e.db.Query(`SELECT id, phone FROM phones WHERE people_id = $1 and fax = $2 ORDER BY phone ASC`, id, fax)
	if err != nil {
		log.Println("GetPeoplePhonesAll e.db.Query ", err)
		return []Phone{}, nil
	}
	phones, err := scanPhones(rows, "short")
	return phones, err
}

// CreatePhone - create new phone
func (e *Edb) CreatePhone(phone Phone) (int64, error) {
	stmt, err := e.db.Prepare(`INSERT INTO phones(company_id, people_id, phone, fax, created_at) VALUES($1, $2, $3, $4, now()) RETURNING id`)
	if err != nil {
		log.Println("CreatePhone e.db.Prepare ", err)
		return 0, err
	}
	err = stmt.QueryRow(i2n(phone.CompanyID), i2n(phone.PeopleID), i2n(phone.Phone), phone.Fax).Scan(&phone.ID)
	if err != nil {
		log.Println("CreatePhone db.QueryRow ", err)
		return 0, err
	}
	return phone.ID, nil
}

// CreateCompanyPhones - create new phones to company
func (e *Edb) CreateCompanyPhones(company Company, fax bool) error {
	err := e.CleanCompanyPhones(company, fax)
	if err != nil {
		log.Println("CreateCompanyPhones CleanCompanyPhones ", err)
		return err
	}
	for _, value := range company.Phones {
		var id int64
		phone := Phone{}
		err = e.db.QueryRow(`SELECT id FROM phones WHERE company_id = $1 and phone = $2 and fax = $3 RETURNING id`, company.ID, value.Phone, fax).Scan(&id)
		if phone.ID == 0 {
			value.CompanyID = company.ID
			value.Fax = fax
			_, err = e.CreatePhone(value)
			if err != nil {
				log.Println("CreateCompanyPhones CreatePhone ", err)
				return err
			}
		}
	}
	return nil
}

// CreatePeoplePhones - create new phones to people
func (e *Edb) CreatePeoplePhones(people People, fax bool) error {
	err := e.CleanPeoplePhones(people, fax)
	if err != nil {
		log.Println("CreatePeoplePhones CleanPeoplePhones ", err)
		return err
	}
	var allPhones []Phone
	if fax {
		allPhones = people.Faxes
	} else {
		allPhones = people.Phones
	}
	for _, value := range allPhones {
		phone := Phone{}
		err = e.db.QueryRow(`SELECT id FROM phones WHERE people_id = $1 and phone = $2 and fax = $3 RETURNING id`, people.ID, value.Phone, fax).Scan(&phone.ID)
		if phone.ID == 0 {
			value.PeopleID = people.ID
			value.Fax = fax
			_, err = e.CreatePhone(value)
			if err != nil {
				log.Println("CreatePeoplePhones CreatePhone ", err)
				return err
			}
		}
	}
	return nil
}

// CleanCompanyPhones - delete all unnecessary phones by company id
func (e *Edb) CleanCompanyPhones(company Company, fax bool) error {
	var (
		phones    []int64
		allPhones []Phone
	)
	if fax {
		allPhones = company.Faxes
	} else {
		allPhones = company.Phones
	}
	for _, value := range allPhones {
		phones = append(phones, value.Phone)
	}
	if len(phones) == 0 {
		_, err := e.db.Exec(`DELETE FROM phones WHERE company_id = $1 and fax = $2`, company.ID, fax)
		if err != nil {
			log.Println("CleanCompanyPhones e.db.Exec ", err)
			return err
		}
	} else {
		rows, err := e.db.Query(`SELECT id, phone FROM phones WHERE company_id = $1 and fax = $2`, company.ID, fax)
		if err != nil {
			log.Println("CleanCompanyPhones e.db.Query ", err)
			return err
		}
		companyPhones, err := scanPhones(rows, "short")
		if err != nil {
			log.Println("CleanCompanyPhones scanPhones ", err)
			return err
		}
		for _, value := range companyPhones {
			if int64InSlice(value.Phone, phones) == false {
				_, err = e.db.Exec(`DELETE FROM phones WHERE company_id = $1 and phone = $2 and fax = $3`, company.ID, value.Phone, fax)
				if err != nil {
					log.Println("CleanCompanyPhones e.db.Exec ", err)
					return err
				}
			}
		}
	}
	return nil
}

// CleanPeoplePhones - delete all unnecessary phones by people id
func (e *Edb) CleanPeoplePhones(people People, fax bool) error {
	var (
		phones    []int64
		allPhones []Phone
	)
	if fax {
		allPhones = people.Faxes
	} else {
		allPhones = people.Phones
	}
	for _, value := range allPhones {
		phones = append(phones, value.Phone)
	}
	if len(phones) == 0 {
		_, err := e.db.Exec(`DELETE FROM phones WHERE people_id = $1 and fax = $2`, people.ID, fax)
		if err != nil {
			log.Println("CleanPeoplePhones e.db.Exec ", err)
			return err
		}
	} else {
		rows, err := e.db.Query(`SELECT id, phone FROM phones WHERE people_id = $1 and fax = $2`, people.ID, fax)
		if err != nil {
			log.Println("CleanPeoplePhones e.db.Query ", err)
			return err
		}
		peoplePhones, err := scanPhones(rows, "short")
		if err != nil {
			log.Println("CleanPeoplePhones scanPhones ", err)
			return err
		}
		for _, value := range peoplePhones {
			if int64InSlice(value.Phone, phones) == false {
				_, err = e.db.Exec(`DELETE FROM phones WHERE people_id = $1 and phone = $2 and fax = $3`, people.ID, value.Phone, fax)
				if err != nil {
					log.Println("CleanPeoplePhones e.db.Exec ", err)
					return err
				}
			}
		}
	}
	return nil
}

// DeleteAllCompanyPhones - delete all phones and faxes by company id
func (e *Edb) DeleteAllCompanyPhones(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`DELETE FROM phones WHERE company_id = $1`, id)
	if err != nil {
		log.Println("DeleteAllCompanyPhones e.db.Exec ", id, err)
	}
	return err
}

// DeleteAllPeoplePhones - delete all phones and faxes by people id
func (e *Edb) DeleteAllPeoplePhones(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`DELETE FROM phones WHERE people_id = $1`, id)
	if err != nil {
		log.Println("DeleteAllPeoplePhones e.db.Exec ", id, err)
	}
	return err
}

func (e *Edb) phoneCreateTable() error {
	str := `CREATE TABLE IF NOT EXISTS phones (id bigserial primary key, people_id bigint, company_id bigint, phone bigint, fax bool NOT NULL DEFAULT false, created_at TIMESTAMP without time zone, updated_at TIMESTAMP without time zone)`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("phoneCreateTable e.db.Exec ", err)
	}
	return err
}
