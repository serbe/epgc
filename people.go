package epgc

import (
	"log"
	"time"

	"database/sql"

	"github.com/lib/pq"
)

// People is struct for people
type People struct {
	ID        int64      `sql:"id" json:"id"`
	Name      string     `sql:"name" json:"name"`
	Company   Company    `sql:"-"`
	CompanyID int64      `sql:"company_id, null" json:"company-id"`
	Post      Post       `sql:"-"`
	PostID    int64      `sql:"post_id, null" json:"post-id"`
	PostGO    Post       `sql:"-"`
	PostGOID  int64      `sql:"post_go_id, null" json:"post-go-id"`
	Rank      Rank       `sql:"-"`
	RankID    int64      `sql:"rank_id, null" json:"rank-id"`
	Birthday  time.Time  `sql:"birthday, null" json:"birthday"`
	Note      string     `sql:"note, null" json:"note"`
	Emails    []Email    `sql:"-"`
	Phones    []Phone    `sql:"-"`
	Faxes     []Phone    `sql:"-"`
	Trainings []Training `sql:"-"`
}

func scanPeople(row *sql.Row) (People, error) {
	var (
		sid        sql.NullInt64
		sname      sql.NullString
		scompanyID sql.NullInt64
		spostID    sql.NullInt64
		spostGOID  sql.NullInt64
		srankID    sql.NullInt64
		sbirthday  pq.NullTime
		snote      sql.NullString
		semails    sql.NullString
		sphones    sql.NullString
		sfaxes     sql.NullString
		// strainings sql.NullString
		people People
	)
	err := row.Scan(&sid, &sname, &scompanyID, &spostID, &spostGOID, &srankID, &snote, &semails, &sphones, &sfaxes)
	if err != nil {
		log.Println("scanPeople row.Scan ", err)
		return People{}, err
	}
	people.ID = n2i(sid)
	people.Name = n2s(sname)
	people.CompanyID = n2i(scompanyID)
	people.PostID = n2i(spostID)
	people.PostGOID = n2i(spostGOID)
	people.RankID = n2i(srankID)
	people.Note = n2s(snote)
	people.Emails = n2emails(semails)
	people.Phones = n2phones(sphones)
	people.Faxes = n2faxes(sfaxes)
	// people.Practices = n2practices(spractices)
	return people, nil
}

func scanPeoples(rows *sql.Rows, opt string) ([]People, error) {
	var (
		peoples []People
		err     error
	)
	for rows.Next() {
		var (
			sid          sql.NullInt64
			sname        sql.NullString
			scompanyName sql.NullString
			spostName    sql.NullString
			sphones      sql.NullString
			sfaxes       sql.NullString
			people       People
		)
		switch opt {
		case "list":
			err := rows.Scan(&sid, &sname, &scompanyName, &spostName, &sphones, &sfaxes)
			if err != nil {
				log.Println("scanPeople rows.Scan list ", err)
				return peoples, err
			}
		case "select":
			err := rows.Scan(&sid, &sname)
			if err != nil {
				log.Println("scanPeople rows.Scan select ", err)
				return peoples, err
			}
		}
		people.ID = n2i(sid)
		switch opt {
		case "list":
			people.Name = n2s(sname)
			people.Company.Name = n2s(scompanyName)
			people.Post.Name = n2s(spostName)
			people.Phones = n2phones(sphones)
			people.Faxes = n2faxes(sfaxes)
		case "select":
			people.Name = n2s(sname)
			if len(people.Name) > 40 {
				people.Name = people.Name[0:40]
			}
		}
		peoples = append(peoples, people)
	}
	err = rows.Err()
	if err != nil {
		log.Println("scanPeoples rows.Err ", err)
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
	GROUP BY p.id
	WHERE id = $1`)
	if err != nil {
		log.Println("GetPeople e.db.Prepare ", err)
		return People{}, err
	}
	row := stmt.QueryRow(id)
	people, err := scanPeople(row)
	// people.Trainings = GetPeopleTrainings(people.ID)
	return people, err
}

// GetPeopleList - get all peoples for list
func (e *Edb) GetPeopleList() ([]People, error) {
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
		return []People{}, err
	}
	peoples, err := scanPeoples(rows, "list")
	return peoples, err
}

// GetPeopleSelect - get all peoples for select
func (e *Edb) GetPeopleSelect() ([]People, error) {
	rows, err := e.db.Query(`SELECT
		p.id,
		p.name
	FROM
		peoples AS p
	ORDER BY name ASC`)
	if err != nil {
		log.Println("GetPeopleSelect e.db.Query ", err)
		return []People{}, err
	}
	peoples, err := scanPeoples(rows, "select")
	return peoples, err
}

// CreatePeople - create new people
func (e *Edb) CreatePeople(people People) (int64, error) {
	stmt, err := e.db.Prepare(`INSERT INTO peoples(name, company_id, post_id, post_go_id, rank_id, birthday, note) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`)
	if err != nil {
		log.Println("CreatePeople e.db.Prepare ", err)
		return 0, err
	}
	err = stmt.QueryRow(s2n(people.Name), i2n(people.CompanyID), i2n(people.PostID), i2n(people.PostGOID), i2n(people.RankID), d2n(people.Birthday), s2n(people.Note)).Scan(&people.ID)
	if err != nil {
		log.Println("CreatePeople db.QueryRow ", err)
		return 0, err
	}
	_ = e.CreatePeopleEmails(people)
	_ = e.CreatePeoplePhones(people, false)
	_ = e.CreatePeoplePhones(people, true)
	// CreatePeopleTrainings(people)
	return people.ID, nil
}

// UpdatePeople - save people changes
func (e *Edb) UpdatePeople(people People) error {
	stmt, err := e.db.Prepare(`UPDATE peoples SET name=$2, company_id=$3, post_id=$4, post_go_id=$5, rank_id=$6, birthday=$7, note=$8 WHERE id = $1`)
	if err != nil {
		log.Println("UpdatePeople e.db.Prepare ", err)
		return err
	}
	_, err = stmt.Exec(i2n(people.ID), s2n(people.Name), i2n(people.CompanyID), i2n(people.PostID), i2n(people.PostGOID), i2n(people.RankID), d2n(people.Birthday), s2n(people.Note))
	if err != nil {
		log.Println("UpdatePeople stmt.Exec ", err)
		return err
	}
	_ = e.CreatePeopleEmails(people)
	_ = e.CreatePeoplePhones(people, false)
	_ = e.CreatePeoplePhones(people, true)
	// CreatePeopleTrainings(people)
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
	e.db.Exec("DELETE FROM peoples WHERE id = $1", id)
	if err != nil {
		log.Println("DeletePeople e.db.Exec ", err)
	}
	return err
}

func (e *Edb) peopleCreateTable() error {
	str := `CREATE TABLE IF NOT EXISTS peoples (id bigserial primary key, name text, company_id bigint, post_id bigint, post_go_id bigint, rank_id bigint, birthday date, note text)`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("peopleCreateTable ", err)
	}
	return err
}
