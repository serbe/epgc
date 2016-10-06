package epgc

import (
	"database/sql"
	"log"
)

// SirenType - struct for sirenType
type SirenType struct {
	ID        int64  `sql:"id" json:"id"`
	Name      string `sql:"name" json:"name"`
	Radius    int64  `sql:"radius" json:"radius"`
	Note      string `sql:"note, null" json:"note"`
	CreatedAt string `sql:"created_at" json:"created_at"`
	UpdatedAt string `sql:"updated_at" json:"updated_at"`
}

func scanSirenType(row *sql.Row) (SirenType, error) {
	var (
		sID       sql.NullInt64
		sName     sql.NullString
		sRadius   sql.NullInt64
		sNote     sql.NullString
		sirenType SirenType
	)
	err := row.Scan(&sID, &sName, &sRadius, &sNote)
	if err != nil {
		log.Println("scanSirenType row.Scan ", err)
		return sirenType, err
	}
	sirenType.ID = n2i(sID)
	sirenType.Name = n2s(sName)
	sirenType.Radius = n2i(sRadius)
	sirenType.Note = n2s(sNote)
	return sirenType, nil
}

func scanSirenTypes(rows *sql.Rows) ([]SirenType, error) {
	var sirenTypes []SirenType
	for rows.Next() {
		var (
			sID       sql.NullInt64
			sName     sql.NullString
			sRadius   sql.NullInt64
			sNote     sql.NullString
			sirenType SirenType
		)
		err := rows.Scan(&sID, &sName, &sRadius, &sNote)
		if err != nil {
			log.Println("scanSirenTypes list rows.Scan ", err)
			return sirenTypes, err
		}
		sirenType.Name = n2s(sName)
		sirenType.Radius = n2i(sRadius)
		sirenType.Note = n2s(sNote)
		sirenType.ID = n2i(sID)
		sirenTypes = append(sirenTypes, sirenType)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanSirenTypes rows.Err ", err)
	}
	return sirenTypes, err
}

func scanSirenTypesSelect(rows *sql.Rows) ([]SelectItem, error) {
	var sirenTypes []SelectItem
	for rows.Next() {
		var (
			sID       sql.NullInt64
			sName     sql.NullString
			sirenType SelectItem
		)
		err := rows.Scan(&sID, &sName)
		if err != nil {
			log.Println("scanSirenTypes select rows.Scan ", err)
			return sirenTypes, err
		}
		sirenType.Name = n2s(sName)
		sirenType.ID = n2i(sID)
		sirenTypes = append(sirenTypes, sirenType)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanSirenTypes rows.Err ", err)
	}
	return sirenTypes, err
}

// GetSirenType - get one sirenType by id
func (e *Edb) GetSirenType(id int64) (SirenType, error) {
	if id == 0 {
		return SirenType{}, nil
	}
	row := e.db.QueryRow(`
		SELECT
			id,
			name,
			radius,
			note
		FROM
			sirenTypes
		WHERE
			id = $1
	`, id)
	sirenType, err := scanSirenType(row)
	return sirenType, err
}

// GetSirenTypeList - get all sirenType for list
func (e *Edb) GetSirenTypeList() ([]SirenType, error) {
	rows, err := e.db.Query(`
		SELECT
			id,
			name,
			radius,
			note
		FROM
			sirenTypes
		ORDER BY
			name ASC
	`)
	if err != nil {
		log.Println("GetSirenTypeList e.db.Query ", err)
		return []SirenType{}, err
	}
	sirenTypes, err := scanSirenTypes(rows)
	return sirenTypes, err
}

// GetSirenTypeSelect - get all sirenType for select
func (e *Edb) GetSirenTypeSelect() ([]SelectItem, error) {
	rows, err := e.db.Query(`
		SELECT
			id,
			name
		FROM
			sirenTypes
		ORDER BY
			name ASC`)
	if err != nil {
		log.Println("GetSirenTypeSelect e.db.Query ", err)
		return []SelectItem{}, err
	}
	sirenTypes, err := scanSirenTypesSelect(rows)
	return sirenTypes, err
}

// CreateSirenType - create new sirenType
func (e *Edb) CreateSirenType(sirenType SirenType) (int64, error) {
	stmt, err := e.db.Prepare(`
		INSERT INTO
			sirenTypes (
				name,
				radius,
				note,
				created_at
			) VALUES (
				$1,
				$2,
				$3,
				now()
			)
		RETURNING
			id
	`)
	if err != nil {
		log.Println("CreateSirenType e.db.Prepare ", err)
		return 0, err
	}
	err = stmt.QueryRow(s2n(sirenType.Name), s2n(sirenType.Note)).Scan(&sirenType.ID)
	if err != nil {
		log.Println("CreateSirenType db.QueryRow ", err)
		return 0, err
	}
	return sirenType.ID, nil
}

// UpdateSirenType - save sirenType changes
func (e *Edb) UpdateSirenType(s SirenType) error {
	stmt, err := e.db.Prepare(`
		UPDATE
			sirenTypes
		SET
			name=$2,
			radius=$3,
			note=$4,
			updated_at = now()
		WHERE
			id = $1`)
	if err != nil {
		log.Println("UpdateSirenType e.db.Prepare ", err)
		return err
	}
	_, err = stmt.Exec(i2n(s.ID), s2n(s.Name), s2n(s.Note))
	if err != nil {
		log.Println("UpdateSirenType stmt.Exec ", err)
	}
	return err
}

// DeleteSirenType - delete sirenType by id
func (e *Edb) DeleteSirenType(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`
		DELETE FROM
			sirenTypes
		WHERE
			id = $1
	`, id)
	if err != nil {
		log.Println("DeleteSirenType e.db.Exec ", id, err)
	}
	return err
}

func (e *Edb) sirenTypeCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			sirenTypes (
				id         bigserial primary key,
				name       text,
				radius     bigint,
				note       text,
				created_at TIMESTAMP without time zone,
				updated_at TIMESTAMP without time zone,
				UNIQUE(name, radius)
			);`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("sirenCreateTable e.db.Exec ", err)
	}
	return err
}
