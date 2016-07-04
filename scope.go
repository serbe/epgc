package epgc

import (
	"database/sql"
	"log"
)

// Scope - struct for scope
type Scope struct {
	ID   int64  `sql:"id" json:"id"`
	Name string `sql:"name" json:"name"`
	Note string `sql:"note, null" json:"note"`
}

func scanScope(row *sql.Row) (scope Scope, err error) {

}

// GetScope - get one scope by id
func (e *Edb) GetScope(id int64) (scope Scope, err error) {
	if id == 0 {
		return
	}
	row := e.db.QueryRow("SELECT id,name,note FROM scopes WHERE id = $1", id)
	scope, err = scanScope(row)
	if err != nil {
		log.Println("GetScope ", err)
	}
	return
}

// GetScopeAll - get all scope
func (e *Edb) GetScopeAll() (scopes []Scope, err error) {
	err = e.db.Model(&scopes).Order("name ASC").Select()
	if err != nil {
		log.Println("GetScopeAll ", err)
		return
	}
	return
}

// CreateScope - create new scope
func (e *Edb) CreateScope(s Scope) (id int64, err error) {
	err = db.QueryRow(`INSERT INTO scopes(name, note) VALUES($1, $2) RETURNING id`, qs(s.Name), qs(s.Note)).Scan(&id)
	if err != nil {
		log.Println("CreateScope ", err)
	}
	return
}

// UpdateScope - save scope changes
func (e *Edb) UpdateScope(scope Scope) (err error) {
	err = e.db.Update(&scope)
	if err != nil {
		log.Println("UpdateScope ", err)
	}
	return
}

// DeleteScope - delete scope by id
func (e *Edb) DeleteScope(id int64) (err error) {
	if id == 0 {
		return
	}
	_, err = e.db.Exec("DELETE FROM scopes WHERE id = ?", id)
	if err != nil {
		log.Println("DeleteScope ", err)
	}
	return
}

func (e *Edb) scopeCreateTable() (err error) {
	str := `CREATE TABLE IF NOT EXISTS scopes (id bigserial primary key, name text, note text)`
	_, err = e.db.Exec(str)
	if err != nil {
		log.Println("scopeCreateTable ", err)
	}
	return
}
