package epgc

import (
	"database/sql"
	"log"
)

// Department - struct for department
type Department struct {
	ID        int64  `sql:"id" json:"id"`
	Name      string `sql:"name" json:"name"`
	Note      string `sql:"note, null" json:"note"`
	CreatedAt string `sql:"created_at" json:"created_at"`
	UpdatedAt string `sql:"updated_at" json:"updated_at"`
}

func scanDepartment(row *sql.Row) (Department, error) {
	var (
		sID        sql.NullInt64
		sName      sql.NullString
		sNote      sql.NullString
		department Department
	)
	err := row.Scan(&sID, &sName, &sNote)
	if err != nil {
		log.Println("scanDepartment row.Scan ", err)
		return department, err
	}
	department.ID = n2i(sID)
	department.Name = n2s(sName)
	department.Note = n2s(sNote)
	return department, nil
}

func scanDepartments(rows *sql.Rows) ([]Department, error) {
	var departments []Department
	for rows.Next() {
		var (
			sID        sql.NullInt64
			sName      sql.NullString
			sNote      sql.NullString
			department Department
		)
		err := rows.Scan(&sID, &sName, &sNote)
		if err != nil {
			log.Println("scanDepartments list rows.Scan ", err)
			return departments, err
		}
		department.ID = n2i(sID)
		department.Name = n2s(sName)
		department.Note = n2s(sNote)
		departments = append(departments, department)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanDepartments rows.Err ", err)
	}
	return departments, err
}

func scanDepartmentsSelect(rows *sql.Rows) ([]SelectItem, error) {
	var departments []SelectItem
	for rows.Next() {
		var (
			sID        sql.NullInt64
			sName      sql.NullString
			department SelectItem
		)
		err := rows.Scan(&sID, &sName)
		if err != nil {
			log.Println("scanDepartmentsSelect rows.Scan ", err)
			return departments, err
		}
		department.ID = n2i(sID)
		department.Name = n2s(sName)
		departments = append(departments, department)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanDepartmentsSelect rows.Err ", err)
	}
	return departments, err
}

// GetDepartment - get one department by id
func (e *Edb) GetDepartment(id int64) (Department, error) {
	if id == 0 {
		return Department{}, nil
	}
	row := e.db.QueryRow(`
		SELECT
			id,
			name,
			note
		FROM
			departments
		WHERE
			id = $1
	`, id)
	department, err := scanDepartment(row)
	return department, err
}

// GetDepartmentList - get all department for list
func (e *Edb) GetDepartmentList() ([]Department, error) {
	rows, err := e.db.Query(`
		SELECT
			id,
			name,
			note
		FROM
			departments
		ORDER BY
			name
		ASC
	`)
	if err != nil {
		log.Println("GetDepartmentList e.db.Query ", err)
		return []Department{}, err
	}
	departments, err := scanDepartments(rows)
	return departments, err
}

// GetDepartmentSelect - get all department for select
func (e *Edb) GetDepartmentSelect() ([]SelectItem, error) {
	rows, err := e.db.Query(`
		SELECT
			id,
			name
		FROM
			departments
		ORDER BY
			name ASC
	`)
	if err != nil {
		log.Println("GetDepartmentSelect e.db.Query ", err)
		return []SelectItem{}, err
	}
	departments, err := scanDepartmentsSelect(rows)
	return departments, err
}

// CreateDepartment - create new department
func (e *Edb) CreateDepartment(department Department) (int64, error) {
	stmt, err := e.db.Prepare(`
		INSERT INTO
			departments (
				name,
				note,
				created_at
			)
		VALUES (
			$1,
			$2,
			now()
		)
		RETURNING
			id
	`)
	if err != nil {
		log.Println("CreateDepartment e.db.Prepare ", err)
		return 0, err
	}
	err = stmt.QueryRow(s2n(department.Name), s2n(department.Note)).Scan(&department.ID)
	if err != nil {
		log.Println("CreateDepartment db.QueryRow ", err)
		return 0, err
	}
	return department.ID, nil
}

// UpdateDepartment - save department changes
func (e *Edb) UpdateDepartment(s Department) error {
	stmt, err := e.db.Prepare(`
		UPDATE
			departments
		SET
			name=$2,
			note=$3,
			updated_at = now()
		WHERE
			id = $1
	`)
	if err != nil {
		log.Println("UpdateDepartment e.db.Prepare ", err)
		return err
	}
	_, err = stmt.Exec(i2n(s.ID), s2n(s.Name), s2n(s.Note))
	if err != nil {
		log.Println("UpdateDepartment stmt.Exec ", err)
	}
	return err
}

// DeleteDepartment - delete department by id
func (e *Edb) DeleteDepartment(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`
		DELETE FROM
			departments
		WHERE
			id = $1
	`, id)
	if err != nil {
		log.Println("DeleteDepartment e.db.Exec ", id, err)
	}
	return err
}

func (e *Edb) departmentCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			departments (
				id bigserial primary key,
				name text,
				note text,
				created_at TIMESTAMP without time zone,
				updated_at TIMESTAMP without time zone
			)
	`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("departmentCreateTable e.db.Exec ", err)
	}
	return err
}
