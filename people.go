package epgc

import (
	"log"

	"database/sql"

	"github.com/lib/pq"
)

// People is struct for people
type People struct {
	ID         int64       `sql:"id" json:"id"`
	Name       string      `sql:"name" json:"name"`
	Company    Company     `sql:"-"`
	CompanyID  int64       `sql:"company_id, null" json:"company_id"`
	Post       Post        `sql:"-"`
	PostID     int64       `sql:"post_id, null" json:"post_id"`
	PostGO     Post        `sql:"-"`
	PostGOID   int64       `sql:"post_go_id, null" json:"post_go_id"`
	Rank       Rank        `sql:"-"`
	RankID     int64       `sql:"rank_id, null" json:"rank_id"`
	Birthday   string      `sql:"birthday, null" json:"birthday"`
	Note       string      `sql:"note, null" json:"note"`
	Emails     []Email     `sql:"-"`
	Phones     []Phone     `sql:"-"`
	Faxes      []Phone     `sql:"-"`
	Educations []Education `sql:"-"`
	CreatedAt  string      `sql:"created_at" json:"created_at"`
	UpdatedAt  string      `sql:"updated_at" json:"updated_at"`
}

// PeopleList is struct for people list
type PeopleList struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	CompanyName string   `json:"company_name"`
	PostName    string   `json:"post_name"`
	Phones      []string `json:"phones"`
	Faxes       []string `json:"faxes"`
}

// PeopleCompany is struct for company
type PeopleCompany struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	PostName   string `json:"post_name"`
	PostGOName string `json:"post_go_name"`
}

func scanPeople(row *sql.Row) (People, error) {
	var (
		sID        sql.NullInt64
		sName      sql.NullString
		sCompanyID sql.NullInt64
		sPostID    sql.NullInt64
		sPostGOID  sql.NullInt64
		sRankID    sql.NullInt64
		sBirthday  pq.NullTime
		sNote      sql.NullString
		sEmails    sql.NullString
		sPhones    sql.NullString
		sFaxes     sql.NullString
		// seducations sql.NullString
		people People
	)
	err := row.Scan(&sID, &sName, &sCompanyID, &sPostID, &sPostGOID, &sRankID, &sBirthday, &sNote, &sEmails, &sPhones, &sFaxes)
	if err != nil {
		log.Println("scanPeople row.Scan ", err)
		return People{}, err
	}
	people.ID = n2i(sID)
	people.Name = n2s(sName)
	people.CompanyID = n2i(sCompanyID)
	people.PostID = n2i(sPostID)
	people.PostGOID = n2i(sPostGOID)
	people.RankID = n2i(sRankID)
	people.Birthday = n2sd(sBirthday)
	people.Note = n2s(sNote)
	people.Emails = n2emails(sEmails)
	people.Phones = n2phones(sPhones)
	people.Faxes = n2faxes(sFaxes)
	// people.Practices = n2practices(spractices)
	return people, nil
}

func scanPeoplesList(rows *sql.Rows) ([]PeopleList, error) {
	var peoples []PeopleList
	for rows.Next() {
		var (
			sID          sql.NullInt64
			sName        sql.NullString
			sCompanyName sql.NullString
			sPostName    sql.NullString
			sPhones      sql.NullString
			sFaxes       sql.NullString
			people       PeopleList
		)
		err := rows.Scan(&sID, &sName, &sCompanyName, &sPostName, &sPhones, &sFaxes)
		if err != nil {
			log.Println("scanPeoplesList rows.Scan ", err)
			return peoples, err
		}
		people.ID = n2i(sID)
		people.Name = n2s(sName)
		people.CompanyName = n2s(sCompanyName)
		people.PostName = n2s(sPostName)
		people.Phones = n2as(sPhones)
		people.Faxes = n2as(sFaxes)
		peoples = append(peoples, people)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanPeoplesList rows.Err ", err)
	}
	return peoples, err
}

func scanPeoplesSelect(rows *sql.Rows) ([]SelectItem, error) {
	var peoples []SelectItem
	for rows.Next() {
		var (
			sID    sql.NullInt64
			sName  sql.NullString
			people SelectItem
		)
		err := rows.Scan(&sID, &sName)
		if err != nil {
			log.Println("scanPeoplesSelect rows.Scan ", err)
			return peoples, err
		}
		people.ID = n2i(sID)
		people.Name = n2s(sName)
		peoples = append(peoples, people)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanPeoplesSelect rows.Err ", err)
	}
	return peoples, err
}

func scanPeoplesCompany(rows *sql.Rows) ([]PeopleCompany, error) {
	var peoples []PeopleCompany
	for rows.Next() {
		var (
			sID         sql.NullInt64
			sName       sql.NullString
			sPostName   sql.NullString
			sPostGOName sql.NullString
			people      PeopleCompany
		)
		err := rows.Scan(&sID, &sName, &sPostName, &sPostGOName)
		if err != nil {
			log.Println("scanPeoplesCompany rows.Scan ", err)
			return peoples, err
		}
		people.ID = n2i(sID)
		people.Name = n2s(sName)
		people.PostName = n2s(sPostName)
		people.PostGOName = n2s(sPostGOName)
		peoples = append(peoples, people)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanPeoplesCompany rows.Err ", err)
	}
	return peoples, err
}

// GetPeople - get one people by id
func (e *Edb) GetPeople(id int64) (People, error) {
	if id == 0 {
		return People{}, nil
	}
	stmt, err := e.db.Prepare(`SELECT
		p.id,
		p.name,
		p.company_id,
		p.post_id,
		p.post_go_id,
		p.rank_id,
		p.birthday,
		p.note,
		array_to_string(array_agg(DISTINCT e.email),',') AS email,
		array_to_string(array_agg(DISTINCT ph.phone),',') AS phone,
		array_to_string(array_agg(DISTINCT f.phone),',') AS fax
	FROM
		peoples AS p
	LEFT JOIN emails AS e ON p.id = e.people_id
	LEFT JOIN phones AS ph ON p.id = ph.people_id AND ph.fax = false
	LEFT JOIN phones AS f ON p.id = f.people_id AND f.fax = true
	WHERE p.id = $1
	GROUP BY p.id`)
	if err != nil {
		log.Println("GetPeople e.db.Prepare ", err)
		return People{}, err
	}
	row := stmt.QueryRow(id)
	people, err := scanPeople(row)
	// people.Educations = GetPeopleEducations(people.ID)
	return people, err
}

// GetPeopleList - get all peoples for list
func (e *Edb) GetPeopleList() ([]PeopleList, error) {
	rows, err := e.db.Query(`SELECT
		p.id,
		p.name,
		c.name AS company_name,
		po.name AS post_name,
		array_to_string(array_agg(DISTINCT ph.phone),',') AS phone,
		array_to_string(array_agg(DISTINCT f.phone),',') AS fax
	FROM
		peoples AS p
	LEFT JOIN companies AS c ON p.company_id = c.id
	LEFT JOIN posts AS po ON p.post_id = po.id
	LEFT JOIN phones AS ph ON p.id = ph.people_id AND ph.fax = false
	LEFT JOIN phones AS f ON p.id = f.people_id AND f.fax = true
	GROUP BY p.id, c.name, po.name
	ORDER BY name ASC`)
	if err != nil {
		log.Println("GetPeopleList e.db.Query ", err)
		return []PeopleList{}, err
	}
	peoples, err := scanPeoplesList(rows)
	return peoples, err
}

// GetPeopleSelect - get all peoples for select
func (e *Edb) GetPeopleSelect() ([]SelectItem, error) {
	rows, err := e.db.Query(`SELECT
		p.id,
		p.name
	FROM
		peoples AS p
	ORDER BY name ASC`)
	if err != nil {
		log.Println("GetPeopleSelect e.db.Query ", err)
		return []SelectItem{}, err
	}
	peoples, err := scanPeoplesSelect(rows)
	return peoples, err
}

// GetPeopleCompany - get all peoples from company
func (e *Edb) GetPeopleCompany(id int64) ([]PeopleCompany, error) {
	stmt, err := e.db.Prepare(`SELECT
		p.id,
		p.name,
		po.name AS post_name,
		pog.name AS post_go_name
	FROM
		peoples AS p
	LEFT JOIN posts AS po ON p.post_id = po.id
	LEFT JOIN posts AS pog ON p.post_go_id = pog.id
	WHERE p.company_id = $1
	ORDER BY name ASC`)
	if err != nil {
		log.Println("GetPeopleCompany e.db.Prepare ", err)
		return []PeopleCompany{}, err
	}
	rows, err := stmt.Query(id)
	if err != nil {
		log.Println("GetPeopleCompany e.db.Query ", err)
		return []PeopleCompany{}, err
	}
	peoples, err := scanPeoplesCompany(rows)
	return peoples, err
}

// CreatePeople - create new people
func (e *Edb) CreatePeople(people People) (int64, error) {
	stmt, err := e.db.Prepare(`INSERT INTO peoples(name, company_id, post_id, post_go_id, rank_id, birthday, note, created_at) VALUES($1, $2, $3, $4, $5, $6, $7, now()) RETURNING id`)
	if err != nil {
		log.Println("CreatePeople e.db.Prepare ", err)
		return 0, err
	}
	err = stmt.QueryRow(s2n(people.Name), i2n(people.CompanyID), i2n(people.PostID), i2n(people.PostGOID), i2n(people.RankID), sd2n(people.Birthday), s2n(people.Note)).Scan(&people.ID)
	if err != nil {
		log.Println("CreatePeople db.QueryRow ", err)
		return 0, err
	}
	_ = e.CreatePeopleEmails(people)
	_ = e.CreatePeoplePhones(people, false)
	_ = e.CreatePeoplePhones(people, true)
	// CreatePeopleEducations(people)
	return people.ID, nil
}

// UpdatePeople - save people changes
func (e *Edb) UpdatePeople(people People) error {
	stmt, err := e.db.Prepare(`UPDATE peoples SET name=$2, company_id=$3, post_id=$4, post_go_id=$5, rank_id=$6, birthday=$7, note=$8, updated_at = now() WHERE id = $1`)
	if err != nil {
		log.Println("UpdatePeople e.db.Prepare ", err)
		return err
	}
	_, err = stmt.Exec(i2n(people.ID), s2n(people.Name), i2n(people.CompanyID), i2n(people.PostID), i2n(people.PostGOID), i2n(people.RankID), sd2n(people.Birthday), s2n(people.Note))
	if err != nil {
		log.Println("UpdatePeople stmt.Exec ", err)
		return err
	}
	_ = e.CreatePeopleEmails(people)
	_ = e.CreatePeoplePhones(people, false)
	_ = e.CreatePeoplePhones(people, true)
	// CreatePeopleEducations(people)
	return nil
}

// DeletePeople - delete people by id
func (e *Edb) DeletePeople(id int64) error {
	if id == 0 {
		return nil
	}
	err := e.DeleteAllPeoplePhones(id)
	if err != nil {
		log.Println("DeletePeople DeleteAllPeoplePhones ", err)
		return err
	}
	e.db.Exec(`DELETE FROM peoples WHERE id = $1`, id)
	if err != nil {
		log.Println("DeletePeople e.db.Exec ", id, err)
	}
	return err
}

func (e *Edb) peopleCreateTable() error {
	str := `CREATE TABLE IF NOT EXISTS peoples (id bigserial primary key, name text, company_id bigint, post_id bigint, post_go_id bigint, rank_id bigint, birthday date, note text, created_at TIMESTAMP without time zone, updated_at TIMESTAMP without time zone)`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("peopleCreateTable ", err)
	}
	return err
}
