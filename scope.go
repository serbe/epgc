package epgc

import (
	"database/sql"
	"log"
)

// Scope - struct for scope
type Scope struct {
	ID        int64  `sql:"id" json:"id"`
	Name      string `sql:"name" json:"name"`
	Note      string `sql:"note, null" json:"note"`
	CreatedAt string `sql:"created_at" json:"created_at"`
	UpdatedAt string `sql:"updated_at" json:"updated_at"`
}

func scanScope(row *sql.Row) (Scope, error) {
	var (
		sID   sql.NullInt64
		sName sql.NullString
		sNote sql.NullString
		scope Scope
	)
	err := row.Scan(&sID, &sName, &sNote)
	if err != nil {
		log.Println("scanScope row.Scan ", err)
		return scope, err
	}
	scope.ID = n2i(sID)
	scope.Name = n2s(sName)
	scope.Note = n2s(sNote)
	return scope, nil
}

func scanScopes(rows *sql.Rows, opt string) ([]Scope, error) {
	var scopes []Scope
	for rows.Next() {
		var (
			sID   sql.NullInt64
			sName sql.NullString
			sNote sql.NullString
			scope Scope
		)
		switch opt {
		case "list":
			err := rows.Scan(&sID, &sName, &sNote)
			if err != nil {
				log.Println("scanScopes rows.Scan list ", err)
				return []Scope{}, err
			}
			scope.Name = n2s(sName)
			scope.Note = n2s(sNote)
		case "select":
			err := rows.Scan(&sID, &sName)
			if err != nil {
				log.Println("scanScopes rows.Scan select ", err)
				return []Scope{}, err
			}
			scope.Name = n2s(sName)
			// if len(scope.Name) > 210 {
			// 	scope.Name = scope.Name[0:210]
			// }
		}
		scope.ID = n2i(sID)
		scopes = append(scopes, scope)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanScopes rows.Err ", err)
	}
	return scopes, err
}

// GetScope - get one scope by id
func (e *Edb) GetScope(id int64) (Scope, error) {
	if id == 0 {
		return Scope{}, nil
	}
	row := e.db.QueryRow(`SELECT id, name, note FROM scopes WHERE id = $1`, id)
	scope, err := scanScope(row)
	return scope, err
}

// GetScopeList - get all scope for list
func (e *Edb) GetScopeList() ([]Scope, error) {
	rows, err := e.db.Query(`SELECT id, name, note FROM scopes ORDER BY name ASC`)
	if err != nil {
		log.Println("GetScopeList e.db.Query ", err)
		return []Scope{}, err
	}
	scopes, err := scanScopes(rows, "list")
	return scopes, err
}

// GetScopeSelect - get all scope for select
func (e *Edb) GetScopeSelect() ([]Scope, error) {
	rows, err := e.db.Query(`SELECT id, name FROM scopes ORDER BY name ASC`)
	if err != nil {
		log.Println("GetScopeSelect e.db.Query ", err)
		return []Scope{}, err
	}
	scopes, err := scanScopes(rows, "select")
	return scopes, err
}

// CreateScope - create new scope
func (e *Edb) CreateScope(scope Scope) (int64, error) {
	stmt, err := e.db.Prepare(`INSERT INTO scopes(name, note, created_at) VALUES($1, $2, now()) RETURNING id`)
	if err != nil {
		log.Println("CreateScope e.db.Prepare ", err)
		return 0, err
	}
	err = stmt.QueryRow(s2n(scope.Name), s2n(scope.Note)).Scan(&scope.ID)
	if err != nil {
		log.Println("CreateScope db.QueryRow ", err)
	}
	return scope.ID, err
}

// UpdateScope - save scope changes
func (e *Edb) UpdateScope(s Scope) error {
	stmt, err := e.db.Prepare(`UPDATE scopes SET name=$2, note=$3, updated_at = now() WHERE id = $1`)
	if err != nil {
		log.Println("UpdateScope e.db.Prepare ", err)
		return err
	}
	_, err = stmt.Exec(i2n(s.ID), s2n(s.Name), s2n(s.Note))
	if err != nil {
		log.Println("UpdateScope stmt.Exec ", err)
	}
	return err
}

// DeleteScope - delete scope by id
func (e *Edb) DeleteScope(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`DELETE FROM scopes WHERE id = $1`, id)
	if err != nil {
		log.Println("DeleteScope e.db.Exec ", id, err)
	}
	return err
}

func (e *Edb) scopeCreateTable() error {
	str := `CREATE TABLE IF NOT EXISTS scopes (id bigserial primary key, name text, note text, created_at TIMESTAMP without time zone, updated_at TIMESTAMP without time zone)`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("scopeCreateTable e.db.Exec ", err)
	}
	return err
}
