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
	var (
		id   sql.NullInt64
		name sql.NullString
		note sql.NullString
	)
	err = row.Scan(&id, &name, &note)
	if err != nil {
		log.Println("scanScope row.Scan ", err)
		return
	}
	scope.ID = n2i(id)
	scope.Name = n2s(name)
	scope.Note = n2s(note)
	return
}

func scanScopes(rows *sql.Rows, opt string) (scopes []Scope, err error) {
	for rows.Next() {
		var (
			id    sql.NullInt64
			name  sql.NullString
			note  sql.NullString
			scope Scope
		)
		switch opt {
		case "list":
			err = rows.Scan(&id, &name, &note)
		case "select":
			err = rows.Scan(&id, &name)
		}
		if err != nil {
			log.Println("scanScopes rows.Scan ", err)
			return
		}
		scope.ID = n2i(id)
		switch opt {
		case "list":
			scope.Name = n2s(name)
			scope.Note = n2s(note)
		case "select":
			scope.Name = n2s(name)
			if len(scope.Name) > 40 {
				scope.Name = scope.Name[0:40]
			}
		}
		scopes = append(scopes, scope)
	}
	err = rows.Err()
	if err != nil {
		log.Println("scanScopes rows.Err ", err)
	}
	return
}

// GetScope - get one scope by id
func (e *Edb) GetScope(id int64) (scope Scope, err error) {
	if id == 0 {
		return
	}
	row := e.db.QueryRow("SELECT id,name,note FROM scopes WHERE id = $1", id)
	scope, err = scanScope(row)
	return
}

// GetScopeList - get all scope for list
func (e *Edb) GetScopeList() (scopes []Scope, err error) {
	rows, err = e.db.Query("SELECT id,name,note FROM scopes ORDER BY name ASC")
	if err != nil {
		log.Println("GetScopeList e.db.Query ", err)
		return
	}
	scopes, err = scanScopes(rows, "list")
	return
}

// GetScopeSelect - get all scope for select
func (e *Edb) GetScopeSelect() (scopes []Scope, err error) {
	rows, err = e.db.Query("SELECT id,name FROM scopes ORDER BY name ASC")
	if err != nil {
		log.Println("GetScopeSelect e.db.Query ", err)
		return
	}
	scopes, err = scanScopes(rows, "select")
	return
}

// CreateScope - create new scope
func (e *Edb) CreateScope(scope Scope) (id int64, err error) {
	stmt, err := e.db.Prepare(`INSERT INTO scopes(name, note) VALUES($1, $2) RETURNING id`)
	if err != nil {
		log.Println("CreateScope e.db.Prepare ", err)
		return
	}
	err = stmt.QueryRow(s2n(scope.Name), s2n(scope.Note)).Scan(&scope.ID)
	if err != nil {
		log.Println("CreateScope db.QueryRow ", err)
	}
	return
}

// UpdateScope - save scope changes
func (e *Edb) UpdateScope(s Scope) (err error) {
	stmt, err := e.db.Prepare("UPDATE scopes SET name=$2,note=$3 WHERE id = $1")
	if err != nil {
		log.Println("UpdateScope e.db.Prepare ", err)
		return
	}
	_, err = stmt.Exec(i2n(s.ID), s2n(s.Name), s2n(s.Note))
	if err != nil {
		log.Println("UpdateScope stmt.Exec ", err)
	}
	return
}

// DeleteScope - delete scope by id
func (e *Edb) DeleteScope(id int64) (err error) {
	if id == 0 {
		return
	}
	_, err = e.db.Exec("DELETE FROM scopes WHERE id = $1", id)
	if err != nil {
		log.Println("DeleteScope e.db.Exec ", err)
	}
	return
}

func (e *Edb) scopeCreateTable() (err error) {
	str := `CREATE TABLE IF NOT EXISTS scopes (id bigserial primary key, name text, note text)`
	_, err = e.db.Exec(str)
	if err != nil {
		log.Println("scopeCreateTable e.db.Exec ", err)
	}
	return
}
