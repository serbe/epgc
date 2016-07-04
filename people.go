package epdc

import (
	"log"
	"time"
)

// People is struct for people
type People struct {
	TableName struct{}   `sql:"peoples"`
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

// GetPeople - get one people by id
func (e *EDc) GetPeople(id int64) (people People, err error) {
	if id == 0 {
		return
	}
	err = e.db.Model(&people).Where("id = ?", id).Select()
	if err != nil {
		log.Println("GetPeople ", err)
		return
	}
	people.Company, _ = e.GetCompany(people.CompanyID)
	people.Post, _ = e.GetPost(people.PostID)
	people.PostGO, _ = e.GetPost(people.PostGOID)
	people.Rank, _ = e.GetRank(people.RankID)
	people.Emails, _ = e.GetPeopleEmails(people.ID)
	people.Phones, _ = e.GetPeoplePhones(people.ID)
	people.Faxes, _ = e.GetPeopleFaxes(people.ID)
	// people.Trainings = GetPeopleTrainings(people.ID)
	return
}

// GetPeopleAll - get all peoples
func (e *EDc) GetPeopleAll() (peoples []People, err error) {
	err = e.db.Model(&peoples).Order("name ASC").Select()
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
func (e *EDc) CreatePeople(people People) (err error) {
	err = e.db.Create(&people)
	if err != nil {
		log.Println("CreatePeople ", err)
		return
	}
	_ = e.CreatePeopleEmails(people)
	_ = e.CreatePeoplePhones(people)
	_ = e.CreatePeopleFaxes(people)
	// CreatePeopleTrainings(people)
	return
}

// UpdatePeople - save people changes
func (e *EDc) UpdatePeople(people People) (err error) {
	err = e.db.Update(&people)
	if err != nil {
		log.Println("UpdatePeople ", err)
		return
	}
	_ = e.CreatePeopleEmails(people)
	_ = e.CreatePeoplePhones(people)
	_ = e.CreatePeopleFaxes(people)
	// CreatePeopleTrainings(people)
	return
}

// DeletePeople - delete people by id
func (e *EDc) DeletePeople(id int64) (err error) {
	if id == 0 {
		return
	}
	err = e.DeleteAllPeoplePhones(id)
	if err != nil {
		log.Println("DeletePeople DeleteAllPeoplePhones ", err)
		return
	}
	e.db.Exec("DELETE FROM peoples WHERE id = ?", id)
	if err != nil {
		log.Println("DeletePeople ", err)
	}
	return
}

func (e *EDc) peopleCreateTable() (err error) {
	str := `CREATE TABLE IF NOT EXISTS peoples (id bigserial primary key, name text, company_id bigint, post_id bigint, post_go_id bigint, rank_id bigint, birthday date, note text)`
	_, err = e.db.Exec(str)
	if err != nil {
		log.Println("peopleCreateTable ", err)
	}
	return
}
