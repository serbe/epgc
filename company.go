package epgc

import (
	"database/sql"
	"log"
)

// Company is struct for company
type Company struct {
	ID        int64           `sql:"id" json:"id"`
	Name      string          `sql:"name" json:"name"`
	Address   string          `sql:"address, null" json:"address"`
	Scope     Scope           `sql:"-"`
	ScopeID   int64           `sql:"scope_id, null" json:"scope_id"`
	Note      string          `sql:"note, null" json:"note"`
	Emails    []Email         `sql:"-"`
	Phones    []Phone         `sql:"-"`
	Faxes     []Phone         `sql:"-"`
	Practices []Practice      `sql:"-"`
	Peoples   []PeopleCompany `sql:"-"`
	CreatedAt string          `sql:"created_at" json:"created_at"`
	UpdatedAt string          `sql:"updated_at" json:"updated_at"`
}

// CompanyList is struct for list company
type CompanyList struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	Address   string   `json:"address"`
	ScopeName string   `json:"scope_name"`
	Emails    []string `json:"emails"`
	Phones    []string `json:"phones"`
	Faxes     []string `json:"faxes"`
	Practices []string `json:"practices"`
}

func scanCompany(row *sql.Row) (Company, error) {
	var (
		sID      sql.NullInt64
		sName    sql.NullString
		sAddress sql.NullString
		sScopeID sql.NullInt64
		sNote    sql.NullString
		sEmails  sql.NullString
		sPhones  sql.NullString
		sFaxes   sql.NullString
		company  Company
	)
	err := row.Scan(&sID, &sName, &sAddress, &sScopeID, &sNote, &sEmails, &sPhones, &sFaxes)
	if err != nil {
		log.Println("scanScope row.Scan ", err)
		return company, err
	}
	company.ID = n2i(sID)
	company.Name = n2s(sName)
	company.Address = n2s(sAddress)
	company.ScopeID = n2i(sScopeID)
	company.Note = n2s(sNote)
	company.Emails = n2emails(sEmails)
	company.Phones = n2phones(sPhones)
	company.Faxes = n2faxes(sFaxes)
	return company, err
}

func scanCompaniesList(rows *sql.Rows) ([]CompanyList, error) {
	var companies []CompanyList
	for rows.Next() {
		var (
			sID        sql.NullInt64
			sName      sql.NullString
			sAddress   sql.NullString
			sScopeName sql.NullString
			sEmails    sql.NullString
			sPhones    sql.NullString
			sFaxes     sql.NullString
			sPractices sql.NullString
			company    CompanyList
		)
		err := rows.Scan(&sID, &sName, &sAddress, &sScopeName, &sEmails, &sPhones, &sFaxes, &sPractices)
		if err != nil {
			log.Println("scanCompaniesList rows.Scan ", err)
			return companies, err
		}
		company.ID = n2i(sID)
		company.Name = n2s(sName)
		company.Address = n2s(sAddress)
		company.ScopeName = n2s(sScopeName)
		company.Emails = n2as(sEmails)
		company.Phones = n2as(sPhones)
		company.Faxes = n2as(sFaxes)
		company.Practices = n2ads(sPractices)
		companies = append(companies, company)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanCompaniesList rows.Err ", err)
	}
	return companies, err
}

func scanCompaniesSelect(rows *sql.Rows) ([]SelectItem, error) {
	var companies []SelectItem
	for rows.Next() {
		var (
			sID     sql.NullInt64
			sName   sql.NullString
			company SelectItem
		)
		err := rows.Scan(&sID, &sName)
		if err != nil {
			log.Println("scanCompaniesSelect rows.Scan ", err)
			return companies, err
		}
		company.ID = n2i(sID)
		company.Name = n2s(sName)
		// if len(company.Name) > 210 {
		// 	company.Name = company.Name[0:210]
		// }
		companies = append(companies, company)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanCompaniesSelect rows.Err ", err)
	}
	return companies, err
}

// GetCompany - get one company by id
func (e *Edb) GetCompany(id int64) (Company, error) {
	if id == 0 {
		return Company{}, nil
	}
	stmt, err := e.db.Prepare(`SELECT
			c.id,
			c.name,
			c.address,
			c.scope_id,
			c.note,
			array_to_string(array_agg(DISTINCT e.email),',') AS email,
			array_to_string(array_agg(DISTINCT p.phone),',') AS phone,
			array_to_string(array_agg(DISTINCT f.phone),',') AS fax
        FROM
			companies AS c
		LEFT JOIN emails AS e ON c.id = e.company_id
		LEFT JOIN phones AS p ON c.id = p.company_id AND p.fax = false
		LEFT JOIN phones AS f ON c.id = f.company_id AND f.fax = true
 		WHERE c.id = $1
		GROUP BY c.id`)
	if err != nil {
		log.Println("GetCompany e.db.Prepare ", err)
		return Company{}, err
	}
	row := stmt.QueryRow(id)
	company, err := scanCompany(row)
	if err != nil {
		log.Println("GetCompany scanCompany ", err)
		return Company{}, err
	}
	company.Practices, err = e.GetPracticeCompany(id)
	return company, err
}

// GetCompanyList - get all companyes for list
func (e *Edb) GetCompanyList() ([]CompanyList, error) {
	rows, err := e.db.Query(`SELECT
			c.id,
			c.name,
			c.address,
			s.name AS scope_name,
			array_to_string(array_agg(DISTINCT e.email),',') AS email,
			array_to_string(array_agg(DISTINCT p.phone),',') AS phone,
			array_to_string(array_agg(DISTINCT f.phone),',') AS fax,
			array_to_string(array_agg(DISTINCT pr.date_of_practice),',') AS practice
        FROM
			companies AS c
		LEFT JOIN scopes AS s ON c.scope_id = s.id
		LEFT JOIN emails AS e ON c.id = e.company_id
		LEFT JOIN phones AS p ON c.id = p.company_id AND p.fax = false
		LEFT JOIN phones AS f ON c.id = f.company_id AND f.fax = true
		LEFT JOIN practices AS pr ON c.id = pr.company_id
		GROUP BY c.id, s.name
		ORDER BY c.name ASC`)
	if err != nil {
		log.Println("GetCompanyList e.db.Query ", err)
		return []CompanyList{}, err
	}
	companies, err := scanCompaniesList(rows)
	return companies, err
}

// GetCompanySelect - get all companyes for select
func (e *Edb) GetCompanySelect() ([]SelectItem, error) {
	rows, err := e.db.Query(`SELECT
			c.id,
			c.name
        FROM
			companies AS c
		ORDER BY c.name ASC`)
	if err != nil {
		log.Println("GetCompanyList e.db.Query ", err)
		return []SelectItem{}, err
	}
	companies, err := scanCompaniesSelect(rows)
	return companies, err
}

// CreateCompany - create new company
func (e *Edb) CreateCompany(company Company) (int64, error) {
	stmt, err := e.db.Prepare(`INSERT INTO companies(name, address, scope_id, note, created_at) VALUES($1, $2, $3, $4, now()) RETURNING id`)
	if err != nil {
		log.Println("CreateCompany e.db.Prepare ", err)
		return 0, err
	}
	err = stmt.QueryRow(s2n(company.Name), s2n(company.Address), i2n(company.ScopeID), s2n(company.Note)).Scan(&company.ID)
	if err != nil {
		log.Println("CreateScope db.QueryRow ", err)
		return 0, err
	}
	_ = e.CreateCompanyEmails(company)
	_ = e.CreateCompanyPhones(company, false)
	_ = e.CreateCompanyPhones(company, true)
	return company.ID, nil
}

// UpdateCompany - save company changes
func (e *Edb) UpdateCompany(company Company) error {
	stmt, err := e.db.Prepare(`UPDATE companies SET name=$2, address=$3, scope_id=$4, note=$5, updated_at = now() WHERE id=$1`)
	if err != nil {
		log.Println("UpdateCompany e.db.Prepare ", err)
		return err
	}
	_, err = stmt.Exec(i2n(company.ID), s2n(company.Name), s2n(company.Address), i2n(company.ScopeID), s2n(company.Note))
	if err != nil {
		log.Println("UpdateCompany stmt.Exec ", err)
		return err
	}
	_ = e.CreateCompanyEmails(company)
	_ = e.CreateCompanyPhones(company, false)
	_ = e.CreateCompanyPhones(company, true)
	return nil
}

// DeleteCompany - delete company by id
func (e *Edb) DeleteCompany(id int64) error {
	if id == 0 {
		return nil
	}
	e.DeleteAllCompanyPhones(id)
	_, err := e.db.Exec(`DELETE FROM companies WHERE id = $1`, id)
	if err != nil {
		log.Println("DeleteCompany e.db.Exec ", id, err)
	}
	return err
}

func (e *Edb) companyCreateTable() error {
	str := `CREATE TABLE IF NOT EXISTS companies (id BIGSERIAL PRIMARY KEY, name TEXT, address TEXT, scope_id BIGINT, note TEXT, created_at timestamp without time zone, updated_at timestamp without time zone, UNIQUE(name, scope_id))`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("companyCreateTable e.db.Exec ", err)
	}
	return err
}
