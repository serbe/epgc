package epgc

import (
	"database/sql"
	"log"
	"time"
)

// Kind - struct for kind
type Kind struct {
	ID        int64     `sql:"id" json:"id"`
	Name      string    `sql:"name" json:"name"`
	Note      string    `sql:"note, null" json:"note"`
	CreatedAt time.Time `sql:"created_at" json:"created_at"`
	UpdatedAt time.Time `sql:"updated_at" json:"updated_at"`
}

func scanKind(row *sql.Row) (Kind, error) {
	var (
		sid   sql.NullInt64
		sname sql.NullString
		snote sql.NullString
		kind  Kind
	)
	err := row.Scan(&sid, &sname, &snote)
	if err != nil {
		log.Println("scanKind row.Scan ", err)
		return kind, err
	}
	kind.ID = n2i(sid)
	kind.Name = n2s(sname)
	kind.Note = n2s(snote)
	return kind, nil
}

func scanKinds(rows *sql.Rows, opt string) ([]Kind, error) {
	var kinds []Kind
	for rows.Next() {
		var (
			sid   sql.NullInt64
			sname sql.NullString
			snote sql.NullString
			kind  Kind
		)
		switch opt {
		case "list":
			err := rows.Scan(&sid, &sname, &snote)
			if err != nil {
				log.Println("scanKinds list rows.Scan ", err)
				return kinds, err
			}
			kind.Name = n2s(sname)
			kind.Note = n2s(snote)
		case "select":
			err := rows.Scan(&sid, &sname)
			if err != nil {
				log.Println("scanKinds select rows.Scan ", err)
				return kinds, err
			}
			kind.Name = n2s(sname)
			// if len(kind.Name) > 210 {
			// 	kind.Name = kind.Name[0:210]
			// }
		}
		kind.ID = n2i(sid)
		kinds = append(kinds, kind)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanKinds rows.Err ", err)
	}
	return kinds, err
}

// GetKind - get one kind by id
func (e *Edb) GetKind(id int64) (Kind, error) {
	if id == 0 {
		return Kind{}, nil
	}
	row := e.db.QueryRow(`SELECT id, name, note FROM kinds WHERE id = $1`, id)
	kind, err := scanKind(row)
	return kind, err
}

// GetKindList - get all kind for list
func (e *Edb) GetKindList() ([]Kind, error) {
	rows, err := e.db.Query(`SELECT id, name, note FROM kinds ORDER BY name ASC`)
	if err != nil {
		log.Println("GetKindList e.db.Query ", err)
		return []Kind{}, err
	}
	kinds, err := scanKinds(rows, "list")
	return kinds, err
}

// GetKindSelect - get all kind for select
func (e *Edb) GetKindSelect() ([]Kind, error) {
	rows, err := e.db.Query(`SELECT id, name FROM kinds ORDER BY name ASC`)
	if err != nil {
		log.Println("GetKindSelect e.db.Query ", err)
		return []Kind{}, err
	}
	kinds, err := scanKinds(rows, "select")
	return kinds, err
}

// CreateKind - create new kind
func (e *Edb) CreateKind(kind Kind) (int64, error) {
	stmt, err := e.db.Prepare(`INSERT INTO kinds(name, note, created_at) VALUES($1, $2, now()) RETURNING id`)
	if err != nil {
		log.Println("CreateKind e.db.Prepare ", err)
		return 0, err
	}
	err = stmt.QueryRow(s2n(kind.Name), s2n(kind.Note)).Scan(&kind.ID)
	if err != nil {
		log.Println("CreateKind db.QueryRow ", err)
		return 0, err
	}
	return kind.ID, nil
}

// UpdateKind - save kind changes
func (e *Edb) UpdateKind(s Kind) error {
	stmt, err := e.db.Prepare(`UPDATE kinds SET name=$2, note=$3, updated_at = now() WHERE id = $1`)
	if err != nil {
		log.Println("UpdateKind e.db.Prepare ", err)
		return err
	}
	_, err = stmt.Exec(i2n(s.ID), s2n(s.Name), s2n(s.Note))
	if err != nil {
		log.Println("UpdateKind stmt.Exec ", err)
	}
	return err
}

// DeleteKind - delete kind by id
func (e *Edb) DeleteKind(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`DELETE FROM kinds WHERE id = $1`, id)
	if err != nil {
		log.Println("DeleteKind e.db.Exec ", err)
	}
	return err
}

func (e *Edb) kindCreateTable() error {
	str := `CREATE TABLE IF NOT EXISTS kinds (id bigserial primary key, name text, note text, created_at TIMESTAMP without time zone, updated_at TIMESTAMP without time zone)`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("kindCreateTable e.db.Exec ", err)
	}
	return err
}
