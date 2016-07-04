package epgc

import (
	"log"
)

// Company is struct for company
type Company struct {
	TableName struct{}   `sql:"companies"`
	ID        int64      `sql:"id" json:"id"`
	Name      string     `sql:"name" json:"name"`
	Address   string     `sql:"address, null" json:"address"`
	Scope     Scope      `sql:"-"`
	ScopeID   int64      `sql:"scope_id, null" json:"scope-id"`
	Note     string     `sql:"note, null" json:"note"`
	Emails    []Email    `sql:"-"`
	Phones    []Phone    `sql:"-"`
	Faxes     []Phone    `sql:"-"`
	Practices []Practice `sql:"-"`
}

// GetCompany - get one company by id
func (e *Edb) GetCompany(id int64) (company Company, err error) {
	if id == 0 {
		return
	}
	err = e.db.Model(&company).Where("id = ?", id).Select()
	if err != nil {
		log.Println("GetCompany Select ", err)
		return
	}
	company.Scope, _ = e.GetScope(company.ScopeID)
	company.Emails, _ = e.GetCompanyEmails(company.ID)
	company.Phones, _ = e.GetCompanyPhones(company.ID)
	company.Faxes, _ = e.GetCompanyFaxes(company.ID)
	company.Practices, _ = e.GetCompanyPractices(company.ID)
	return
}

// GetCompanyAll - get all companyes
func (e *Edb) GetCompanyAll() (companyes []Company, err error) {
	err = e.db.Model(&companyes).Order("name ASC").Select()
	if err != nil {
		log.Println("GetCompanyAll Select ", err)
		return
	}
	for i := range companyes {
		companyes[i].Scope, _ = e.GetScope(companyes[i].ScopeID)
		companyes[i].Emails, _ = e.GetCompanyEmails(companyes[i].ID)
		companyes[i].Phones, _ = e.GetCompanyPhones(companyes[i].ID)
		companyes[i].Faxes, _ = e.GetCompanyFaxes(companyes[i].ID)
		companyes[i].Practices, _ = e.GetCompanyPractices(companyes[i].ID)
		for j := range companyes[i].Practices {
			companyes[i].Practices[j].DateStr = companyes[i].Practices[j].DateOfPractice.Format("02.01.2006")
		}
	}
	return
}

// CreateCompany - create new company
func (e *Edb) CreateCompany(company Company) (err error) {
	err = e.db.Create(&company)
	if err != nil {
		log.Println("CreateCompany e.db.Create ", err)
		return
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
