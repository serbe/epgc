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
		pid        sql.NullInt64
		pname      sql.NullString
		pcompanyID sql.NullInt64
		ppostID    sql.NullInt64
		ppostGOID  sql.NullInt64
		prankID    sql.NullInt64
		pbirthday  pq.NullTime
		pnote      sql.NullString
		pemails    sql.NullString
		pphones    sql.NullString
		pfaxes     sql.NullString
		// strainings sql.NullString
		people People
	)
	err := row.Scan(&pid, &pname, &pcompanyID, &ppostID, &ppostGOID, &prankID, &pnote, &pemails, &pphones, &pfaxes)
	if err != nil {
		log.Println("scanPeople row.Scan ", err)
		return People{}, err
	}
	people.ID = n2i(pid)
	people.Name = n2s(pname)
	people.CompanyID = n2i(pcompanyID)
	people.PostID = n2i(ppostID)
	people.PostGOID = n2i(ppostGOID)
	people.RankID = n2i(prankID)
	people.Note = n2s(pnote)
	people.Emails = n2emails(pemails)
	people.Phones = n2phones(pphones)
	people.Faxes = n2faxes(pfaxes)
	// people.Practices = n2practices(spractices)
	return people, nil
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
	people, err = scanPeople(row)
	// people.Trainings = GetPeopleTrainings(people.ID)
	return people, err
}

// GetPeopleList - get all peoples for list
func (e *Edb) GetPeopleList() ([]People, error) {
	rows, err = e.db.Query(&peoples).Order("name ASC").Select()
	if err != nil {
		log.Println("GetPeopleAll ", err)
		return
	}
	for i := range peoples {
		peoples[i].Company, _ = e.GetCompany(peoples[i].CompanyID)
		peoples[i].Post, _ = e.GetPost(peoples[i].PostID)
		peoples[i].PostGO, _ = e.GetPost(peoples[i].PostGOID)
		peoples[i].Rank, _ = e.GetRank(peoples[i].RankID)
		peoples[i].Emails, _ = e.GetPeopleEmails(peoples[i].ID)
		peoples[i].Phones, _ = e.GetPeoplePhones(peoples[i].ID)
		peoples[i].Faxes, _ = e.GetPeopleFaxes(peoples[i].ID)
		// people[i].Trainings = GetPeopleTrainings(people[i].ID)
	}
	return
}

// CreatePeople - create new people
func (e *Edb) CreatePeople(people People) (err error) {

	err = e.db.Create(&people)
	if err != nil {
		log.Println("CreatePeople ", err)
		return err
	}
	_ = e.CreatePeopleEmails(people)
	_ = e.CreatePeoplePhones(people)
	_ = e.CreatePeopleFaxes(people)
	// CreatePeopleTrainings(people)
	return err
}

// UpdatePeople - save people changes
func (e *Edb) UpdatePeople(people People) error {
	err = e.db.Update(&people)
	if err != nil {
		log.Println("UpdatePeople ", err)
		return err
	}
	_ = e.CreatePeopleEmails(people)
	_ = e.CreatePeoplePhones(people)
	_ = e.CreatePeopleFaxes(people)
	// CreatePeopleTrainings(people)
	return err
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
		log.Println("DeletePeople ", err)
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
