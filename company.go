package epgc

import (
	"database/sql"
	"log"
)

// Company is struct for company
type Company struct {
	ID        int64      `sql:"id" json:"id"`
	Name      string     `sql:"name" json:"name"`
	Address   string     `sql:"address, null" json:"address"`
	Scope     Scope      `sql:"-"`
	ScopeID   int64      `sql:"scope_id, null" json:"scope-id"`
	Note      string     `sql:"note, null" json:"note"`
	Emails    []Email    `sql:"-"`
	Phones    []Phone    `sql:"-"`
	Faxes     []Phone    `sql:"-"`
	Practices []Practice `sql:"-"`
}

func scanCompany(row *sql.Row) (company Company, err error) {
	var (
		id        sql.NullInt64
		name      sql.NullString
		address   sql.NullString
		scopeID   sql.NullInt64
		note      sql.NullString
		scope     sql.NullString
		emails    sql.NullString
		phones    sql.NullString
		faxes     sql.NullString
		practices sql.NullString
	)
	err = row.Scan(&id, &name, &address, &scopeID, &note, &scope, &emails, &phones, &faxes, &practices)
	if err != nil {
		log.Println("scanScope row.Scan ", err)
		return
	}
	company.ID = n2i(id)
	company.Name = n2s(name)
	company.Address = n2s(address)
	company.ScopeID = n2i(scopeID)
	company.Note = n2s(note)
	company.Scope.Name = n2s(scope)
	company.Emails = n2emails(emails)
	company.Phones = n2phones(phones)
	company.Faxes = n2faxes(faxes)
	company.Practices = n2practices(practices)
	return
}

func scanCompanies(rows *sql.Rows, opt string) (companies []Company, err error) {
	for rows.Next() {
		var (
			id        sql.NullInt64
			name      sql.NullString
			address   sql.NullString
			scopeID   sql.NullInt64
			note      sql.NullString
			scope     sql.NullString
			emails    sql.NullString
			phones    sql.NullString
			faxes     sql.NullString
			practices sql.NullString
			company   Company
		)
		switch opt {
		case "list":
			err = rows.Scan(&id, &name, &address, &scope, &phones, &faxes, &practices)
		case "select":
			err = rows.Scan(&id, &name)
		}
		if err != nil {
			log.Println("scanCompanies rows.Scan ", err)
			return
		}
		switch opt {
		case "list":
			company.Name = n2s(name)
			company.Address = n2s(address)
			company.Scope.Name = n2s(scope)
			company.Phones = n2phones(phones)
			company.Faxes = n2faxes(faxes)
			company.Practices = n2practices(practices)
		case "select":
			company.Name = n2s(name)
			if len(company.Name) > 40 {
				company.Name = company.Name[0:40]
			}
		}
		companies = append(companies, company)
	}
	err = rows.Err()
	if err != nil {
		log.Println("scanScopes rows.Err ", err)
	}
	return

}

// GetCompany - get one company by id
func (e *Edb) GetCompany(id int64) (company Company, err error) {
	if id == 0 {
		return
	}
	stmt, err := e.db.Prepare(`SELECT
			c.id,
			c.name,
			c.address,
			c.scope_id,
			c.note,
			s.name AS scope_name,
			array_to_string(array_agg(DISTINCT e.email),',') AS email,
			array_to_string(array_agg(DISTINCT p.phone),',') AS phone,
			array_to_string(array_agg(DISTINCT f.phone),',') AS fax,
			array_to_string(array_agg(DISTINCT pr.topic),',') AS practice 
        FROM
			companies AS c 
		LEFT JOIN scopes AS s ON c.scope_id = s.id
		LEFT JOIN emails AS e ON c.id = e.company_id
		LEFT JOIN phones AS p ON c.id = p.company_id AND p.fax = false
		LEFT JOIN phones AS f ON c.id = f.company_id AND f.fax = true
		LEFT JOIN practices AS pr ON c.id = pr.company_id
		GROUP BY c.id, s.name
 		WHERE id = $1`)
	if err != nil {
		log.Println("GetCompany e.db.Prepare ", err)
		return
	}
	row := stmt.QueryRow(id)
	company, err = scanCompany(row)
	return
}

// GetCompanyList - get all companyes for list
func (e *Edb) GetCompanyList() (companies []Company, err error) {
	rows, err := e.db.Query(`SELECT
			c.id,
			c.name,
			c.address,
			s.name AS scope_name,
			array_to_string(array_agg(DISTINCT p.phone),',') AS phone,
			array_to_string(array_agg(DISTINCT f.phone),',') AS fax,
			array_to_string(array_agg(DISTINCT pr.topic),',') AS practice 
        FROM
			companies AS c 
		LEFT JOIN scopes AS s ON c.scope_id = s.id
		LEFT JOIN phones AS p ON c.id = p.company_id AND p.fax = false
		LEFT JOIN phones AS f ON c.id = f.company_id AND f.fax = true
		LEFT JOIN practices AS pr ON c.id = pr.company_id
		GROUP BY c.id, s.name
		ORDER BY c.name ASC`)
	if err != nil {
		log.Println("GetCompanyList e.db.Query ", err)
		return
	}
	companies, err = scanCompanies(rows, "list")
	return
}

// GetCompanySelect - get all companyes for select
func (e *Edb) GetCompanySelect() (companies []Company, err error) {
	rows, err := e.db.Query(`SELECT
			c.id,
			c.name
        FROM
			companies AS c 
		ORDER BY c.name ASC`)
	if err != nil {
		log.Println("GetCompanyList e.db.Query ", err)
		return
	}
	companies, err = scanCompanies(rows, "select")
	return
}

// CreateCompany - create new company
func (e *Edb) CreateCompany(company Company) (err error) {
	stmt, err := e.db.Prepare(`INSERT INTO companies(name, address, scope_id, note) VALUES($1, $2, $3, $4) RETURNING id`)
	if err != nil {
		log.Println("CreateCompany e.db.Prepare ", err)
		return
	}
	err = stmt.QueryRow(s2n(company.Name), s2n(company.Address), i2n(company.ScopeID), s2n(company.Note)).Scan(&company.ID)
	if err != nil {
		log.Println("CreateScope db.QueryRow ", err)
	}
	e.CreateCompanyEmails(company)
	e.CreateCompanyPhones(company)
	e.CreateCompanyFaxes(company)
	return
}

// UpdateCompany - save company changes
func (e *Edb) UpdateCompany(company Company) (err error) {
	err = e.db.Update(&company)
	if err != nil {
		log.Println("UpdateCompany e.db.Update ", err)
		return
	}
	e.CreateCompanyEmails(company)
	e.CreateCompanyPhones(company)
	e.CreateCompanyFaxes(company)
	// CreateCompanyPractices(c)
	return
}

// DeleteCompany - delete company by id
func (e *Edb) DeleteCompany(id int64) (err error) {
	if id == 0 {
		return
	}
	e.DeleteAllCompanyPhones(id)
	_, err = e.db.Exec("DELETE FROM companies WHERE id=?", id)
	if err != nil {
		log.Println("DeleteCompany e.db.Exec ", err)
	}
	return
}

func (e *Edb) companyCreateTable() (err error) {
	str := `CREATE TABLE IF NOT EXISTS companies (id BIGSERIAL PRIMARY KEY, name TEXT, address TEXT, scope_id BIGINT, note TEXT, UNIQUE(name, scope_id))`
	_, err = e.db.Exec(str)
	if err != nil {
		log.Println("companyCreateTable e.db.Exec ", err)
	}
	return
}
