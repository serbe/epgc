package epdc

import "log"

// Scope - struct for scope
type Scope struct {
	TableName struct{} `sql:"scopes"`
	ID        int64    `sql:"id" json:"id"`
	Name      string   `sql:"name" json:"name"`
	Note      string   `sql:"note, null" json:"note"`
}

// GetScope - get one scope by id
func (e *EDc) GetScope(id int64) (scope Scope, err error) {
	if id == 0 {
		return
	}
	err = e.db.Model(&scope).Where("id = ?", id).Select()
	if err != nil {
		log.Println("GetScope ", err)
	}
	return
}

// GetScopeAll - get all scope
func (e *EDc) GetScopeAll() (scopes []Scope, err error) {
	err = e.db.Model(&scopes).Order("name ASC").Select()
	if err != nil {
		log.Println("GetScopeAll ", err)
		return
	}
	return
}

// CreateScope - create new scope
func (e *EDc) CreateScope(scope Scope) (err error) {
	err = e.db.Create(&scope)
	if err != nil {
		log.Println("CreateScope ", err)
	}
	return
}

// UpdateScope - save scope changes
func (e *EDc) UpdateScope(scope Scope) (err error) {
	err = e.db.Update(&scope)
	if err != nil {
		log.Println("UpdateScope ", err)
	}
	return
}

// DeleteScope - delete scope by id
func (e *EDc) DeleteScope(id int64) (err error) {
	if id == 0 {
		return
	}
	_, err = e.db.Exec("DELETE FROM scopes WHERE id = ?", id)
	if err != nil {
		log.Println("DeleteScope ", err)
	}
	return
}

func (e *EDc) scopeCreateTable() (err error) {
	str := `CREATE TABLE IF NOT EXISTS scopes (id bigserial primary key, name text, note text)`
	_, err = e.db.Exec(str)
	if err != nil {
		log.Println("scopeCreateTable ", err)
	}
	return
}
