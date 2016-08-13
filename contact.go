package epgc

import (
	"log"

	"database/sql"

	"github.com/lib/pq"
)

// Contact is struct for contact
type Contact struct {
	ID           int64       `sql:"id" json:"id"`
	Name         string      `sql:"name" json:"name"`
	Company      Company     `sql:"-"`
	CompanyID    int64       `sql:"company_id, null" json:"company_id"`
	Department   Department  `sql:"-"`
	DepartmentID int64       `sql:"department_id, null" json:"department_id"`
	Post         Post        `sql:"-"`
	PostID       int64       `sql:"post_id, null" json:"post_id"`
	PostGO       Post        `sql:"-"`
	PostGOID     int64       `sql:"post_go_id, null" json:"post_go_id"`
	Rank         Rank        `sql:"-"`
	RankID       int64       `sql:"rank_id, null" json:"rank_id"`
	Birthday     string      `sql:"birthday, null" json:"birthday"`
	Note         string      `sql:"note, null" json:"note"`
	Emails       []Email     `sql:"-"`
	Phones       []Phone     `sql:"-"`
	Faxes        []Phone     `sql:"-"`
	Educations   []Education `sql:"-"`
	CreatedAt    string      `sql:"created_at" json:"created_at"`
	UpdatedAt    string      `sql:"updated_at" json:"updated_at"`
}

// ContactList is struct for contact list
type ContactList struct {
	ID             int64    `json:"id"`
	Name           string   `json:"name"`
	CompanyName    string   `json:"company_name"`
	DepartmentName string   `json:"department_name"`
	PostName       string   `json:"post_name"`
	Phones         []string `json:"phones"`
	Faxes          []string `json:"faxes"`
}

// ContactCompany is struct for company
type ContactCompany struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	DepartmentName string `json:"department_name"`
	PostName       string `json:"post_name"`
	PostGOName     string `json:"post_go_name"`
}

func scanContact(row *sql.Row) (Contact, error) {
	var (
		sID           sql.NullInt64
		sName         sql.NullString
		sCompanyID    sql.NullInt64
		sDepartmentID sql.NullInt64
		sPostID       sql.NullInt64
		sPostGOID     sql.NullInt64
		sRankID       sql.NullInt64
		sBirthday     pq.NullTime
		sNote         sql.NullString
		sEmails       sql.NullString
		sPhones       sql.NullString
		sFaxes        sql.NullString
		// seducations sql.NullString
		contact Contact
	)
	err := row.Scan(&sID, &sName, &sCompanyID, &sDepartmentID, &sPostID, &sPostGOID, &sRankID, &sBirthday, &sNote, &sEmails, &sPhones, &sFaxes)
	if err != nil {
		log.Println("scanContact row.Scan ", err)
		return Contact{}, err
	}
	contact.ID = n2i(sID)
	contact.Name = n2s(sName)
	contact.CompanyID = n2i(sCompanyID)
	contact.DepartmentID = n2i(sDepartmentID)
	contact.PostID = n2i(sPostID)
	contact.PostGOID = n2i(sPostGOID)
	contact.RankID = n2i(sRankID)
	contact.Birthday = n2sd(sBirthday)
	contact.Note = n2s(sNote)
	contact.Emails = n2emails(sEmails)
	contact.Phones = n2phones(sPhones)
	contact.Faxes = n2faxes(sFaxes)
	// contact.Practices = n2practices(spractices)
	return contact, nil
}

func scanContactsList(rows *sql.Rows) ([]ContactList, error) {
	var contacts []ContactList
	for rows.Next() {
		var (
			sID             sql.NullInt64
			sName           sql.NullString
			sCompanyName    sql.NullString
			sDepartmentName sql.NullString
			sPostName       sql.NullString
			sPhones         sql.NullString
			sFaxes          sql.NullString
			contact         ContactList
		)
		err := rows.Scan(&sID, &sName, &sCompanyName, &sDepartmentName, &sPostName, &sPhones, &sFaxes)
		if err != nil {
			log.Println("scanContactsList rows.Scan ", err)
			return contacts, err
		}
		contact.ID = n2i(sID)
		contact.Name = n2s(sName)
		contact.CompanyName = n2s(sCompanyName)
		contact.DepartmentName = n2s(sDepartmentName)
		contact.PostName = n2s(sPostName)
		contact.Phones = n2as(sPhones)
		contact.Faxes = n2as(sFaxes)
		contacts = append(contacts, contact)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanContactsList rows.Err ", err)
	}
	return contacts, err
}

func scanContactsSelect(rows *sql.Rows) ([]SelectItem, error) {
	var contacts []SelectItem
	for rows.Next() {
		var (
			sID     sql.NullInt64
			sName   sql.NullString
			contact SelectItem
		)
		err := rows.Scan(&sID, &sName)
		if err != nil {
			log.Println("scanContactsSelect rows.Scan ", err)
			return contacts, err
		}
		contact.ID = n2i(sID)
		contact.Name = n2s(sName)
		contacts = append(contacts, contact)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanContactsSelect rows.Err ", err)
	}
	return contacts, err
}

func scanContactsCompany(rows *sql.Rows) ([]ContactCompany, error) {
	var contacts []ContactCompany
	for rows.Next() {
		var (
			sID         sql.NullInt64
			sName       sql.NullString
			sPostName   sql.NullString
			sPostGOName sql.NullString
			contact     ContactCompany
		)
		err := rows.Scan(&sID, &sName, &sPostName, &sPostGOName)
		if err != nil {
			log.Println("scanContactsCompany rows.Scan ", err)
			return contacts, err
		}
		contact.ID = n2i(sID)
		contact.Name = n2s(sName)
		contact.PostName = n2s(sPostName)
		contact.PostGOName = n2s(sPostGOName)
		contacts = append(contacts, contact)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanContactsCompany rows.Err ", err)
	}
	return contacts, err
}

// GetContact - get one contact by id
func (e *Edb) GetContact(id int64) (Contact, error) {
	if id == 0 {
		return Contact{}, nil
	}
	stmt, err := e.db.Prepare(`SELECT
		c.id,
		c.name,
		c.company_id,
    c.department_id,
		c.post_id,
		c.post_go_id,
		c.rank_id,
		c.birthday,
		c.note,
		array_to_string(array_agg(DISTINCT e.email),',') AS email,
		array_to_string(array_agg(DISTINCT p.phone),',') AS phone,
		array_to_string(array_agg(DISTINCT f.phone),',') AS fax
	FROM
		contacts AS c
	LEFT JOIN emails AS e ON c.id = e.contact_id
	LEFT JOIN phones AS p ON c.id = p.contact_id AND p.fax = false
	LEFT JOIN phones AS f ON c.id = f.contact_id AND f.fax = true
	WHERE c.id = $1
	GROUP BY c.id`)
	if err != nil {
		log.Println("GetContact e.db.Prepare ", err)
		return Contact{}, err
	}
	row := stmt.QueryRow(id)
	contact, err := scanContact(row)
	// contact.Educations = GetContactEducationscontacte.ID)
	return contact, err
}

// GetContactList - get all contacts for list
func (e *Edb) GetContactList() ([]ContactList, error) {
	rows, err := e.db.Query(`SELECT
		c.id,
		c.name,
		co.name AS company_name,
    d.name AS department_name,
		po.name AS post_name,
		array_to_string(array_agg(DISTINCT ph.phone),',') AS phone,
		array_to_string(array_agg(DISTINCT f.phone),',') AS fax
	FROM
		contacts AS c
	LEFT JOIN companies AS co ON c.company_id = co.id
  LEFT JOIN departments AS d ON c.department_id = d.id
	LEFT JOIN posts AS po ON c.post_id = po.id
	LEFT JOIN phones AS ph ON c.id = ph.contact_id AND ph.fax = false
	LEFT JOIN phones AS f ON c.id = f.contact_id AND f.fax = true
	GROUP BY c.id, co.name, po.name
	ORDER BY name ASC`)
	if err != nil {
		log.Println("GetContactList e.db.Query ", err)
		return []ContactList{}, err
	}
	contacts, err := scanContactsList(rows)
	return contacts, err
}

// GetContactSelect - get all contacts for select
func (e *Edb) GetContactSelect() ([]SelectItem, error) {
	rows, err := e.db.Query(`SELECT
		c.id,
		c.name
	FROM
		contacts AS c
	ORDER BY name ASC`)
	if err != nil {
		log.Println("GetContactSelect e.db.Query ", err)
		return []SelectItem{}, err
	}
	contacts, err := scanContactsSelect(rows)
	return contacts, err
}

// GetContactCompany - get all contacts from company
func (e *Edb) GetContactCompany(id int64) ([]ContactCompany, error) {
	stmt, err := e.db.Prepare(`SELECT
		c.id,
		c.name,
		po.name AS post_name,
		pog.name AS post_go_name
	FROM
		contacts AS c
	LEFT JOIN posts AS po ON c.post_id = po.id
	LEFT JOIN posts AS pog ON c.post_go_id = pog.id
	WHERE c.company_id = $1
	ORDER BY name ASC`)
	if err != nil {
		log.Println("GetContactCompany e.db.Prepare ", err)
		return []ContactCompany{}, err
	}
	rows, err := stmt.Query(id)
	if err != nil {
		log.Println("GetContactCompany e.db.Query ", err)
		return []ContactCompany{}, err
	}
	contacts, err := scanContactsCompany(rows)
	return contacts, err
}

// CreateContact - create new contact
func (e *Edb) CreateContact(contact Contact) (int64, error) {
	stmt, err := e.db.Prepare(`INSERT INTO contacts(name, company_id, department_id, post_id, post_go_id, rank_id, birthday, note, created_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, now()) RETURNING id`)
	if err != nil {
		log.Println("CreateContact e.db.Prepare ", err)
		return 0, err
	}
	err = stmt.QueryRow(s2n(contact.Name), i2n(contact.CompanyID), i2n(contact.DepartmentID), i2n(contact.PostID), i2n(contact.PostGOID), i2n(contact.RankID), sd2n(contact.Birthday), s2n(contact.Note)).Scan(&contact.ID)
	if err != nil {
		log.Println("CreateContact db.QueryRow ", err)
		return 0, err
	}
	_ = e.CreateContactEmails(contact)
	_ = e.CreateContactPhones(contact, false)
	_ = e.CreateContactPhones(contact, true)
	// CreateContactEducations(contact)
	return contact.ID, nil
}

// UpdateContact - save contact changes
func (e *Edb) UpdateContact(contact Contact) error {
	stmt, err := e.db.Prepare(`UPDATE contacts SET name=$2, company_id=$3, department_id=$4, post_id=$5, post_go_id=$5, rank_id=$6, birthday=$8, note=$9, updated_at = now() WHERE id = $1`)
	if err != nil {
		log.Println("UpdateContact e.db.Prepare ", err)
		return err
	}
	_, err = stmt.Exec(i2n(contact.ID), s2n(contact.Name), i2n(contact.CompanyID), i2n(contact.DepartmentID), i2n(contact.PostID), i2n(contact.PostGOID), i2n(contact.RankID), sd2n(contact.Birthday), s2n(contact.Note))
	if err != nil {
		log.Println("UpdateContact stmt.Exec ", err)
		return err
	}
	_ = e.CreateContactEmails(contact)
	_ = e.CreateContactPhones(contact, false)
	_ = e.CreateContactPhones(contact, true)
	// CreateContactEducations(contact)
	return nil
}

// DeleteContact - delete contact by id
func (e *Edb) DeleteContact(id int64) error {
	if id == 0 {
		return nil
	}
	err := e.DeleteAllContactPhones(id)
	if err != nil {
		log.Println("DeleteContact DeleteAllContactPhones ", err)
		return err
	}
	e.db.Exec(`DELETE FROM contacts WHERE id = $1`, id)
	if err != nil {
		log.Println("DeleteContact e.db.Exec ", id, err)
	}
	return err
}

func (e *Edb) contactCreateTable() error {
	str := `CREATE TABLE IF NOT EXISTS contacts (id bigserial primary key, name text, company_id bigint, post_id bigint, post_go_id bigint, rank_id bigint, birthday date, note text, created_at TIMESTAMP without time zone, updated_at TIMESTAMP without time zone, department_id bigint)`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("contactCreateTable ", err)
	}
	return err
}
