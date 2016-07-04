package epdc

import "log"

// Kind - struct for kind
type Kind struct {
	TableName struct{} `sql:"kinds"`
	ID        int64    `sql:"id" json:"id"`
	Name      string   `sql:"name" json:"name"`
	Note      string   `sql:"note, null" json:"note"`
}

// GetKind - get one kind by id
func (e *EDc) GetKind(id int64) (kind Kind, err error) {
	if id == 0 {
		return
	}
	err = e.db.Model(&kind).Where("id = ?", id).Select()
	if err != nil {
		log.Println("GetKind ", err)
	}
	return
}

// GetKindAll - get all kinds
func (e *EDc) GetKindAll() (kinds []Kind, err error) {
	err = e.db.Model(&kinds).Order("name ASC").Select()
	if err != nil {
		log.Println("GetKindAll ", err)
		return
	}
	return
}

// CreateKind - create new kind
func (e *EDc) CreateKind(kind Kind) (err error) {
	err = e.db.Create(&kind)
	if err != nil {
		log.Println("CreateKind ", err)
	}
	return
}

// UpdateKind - save kind changes
func (e *EDc) UpdateKind(kind Kind) (err error) {
	err = e.db.Update(&kind)
	if err != nil {
		log.Println("UpdateKind ", err)
	}
	return
}

// DeleteKind - delete kind by id
func (e *EDc) DeleteKind(id int64) (err error) {
	if id == 0 {
		return
	}
	_, err = e.db.Exec("DELETE FROM kinds WHERE id = ?", id)
	if err != nil {
		log.Println("DeleteKind ", err)
	}
	return
}

func (e *EDc) kindCreateTable() (err error) {
	str := `CREATE TABLE IF NOT EXISTS kinds (id bigserial primary key, name text, note text)`
	_, err = e.db.Exec(str)
	if err != nil {
		log.Println("kindCreateTable ", err)
	}
	return
}
