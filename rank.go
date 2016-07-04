package epgc

import "fmt"

// Rank - struct for rank
type Rank struct {
	TableName struct{} `sql:"ranks"`
	ID        int64    `sql:"id" json:"id"`
	Name      string   `sql:"name" json:"name"`
	Note      string   `sql:"note, null" json:"note"`
}

// GetRank - get one rank dy id
func (e *Edb) GetRank(id int64) (rank Rank, err error) {
	if id == 0 {
		return rank, nil
	}
	err = e.db.Model(&rank).Where("id = ?", id).Select()
	if err != nil {
		return rank, fmt.Errorf("GetRank: %s", err)
	}
	return
}

// GetRankAll - get all rank
func (e *Edb) GetRankAll() (ranks []Rank, err error) {
	err = e.db.Model(&ranks).Order("name ASC").Select()
	if err != nil {
		return ranks, fmt.Errorf("GetRankAll: %s", err)
	}
	return
}

// CreateRank - create new rank
func (e *Edb) CreateRank(rank Rank) (err error) {
	err = e.db.Create(&rank)
	if err != nil {
		return fmt.Errorf("CreateRank: %s", err)
	}
	return
}

// UpdateRank - save rank changes
func (e *Edb) UpdateRank(rank Rank) (err error) {
	err = e.db.Update(&rank)
	if err != nil {
		return fmt.Errorf("UpdateRank: %s", err)
	}
	return
}

// DeleteRank - delete rank by id
func (e *Edb) DeleteRank(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec("DELETE FROM ranks WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("DeleteRank: %s", err)
	}
	return nil
}

func (e *Edb) rankCreateTable() (err error) {
	str := `CREATE TABLE IF NOT EXISTS ranks (id bigserial primary key, name text, note text)`
	_, err = e.db.Exec(str)
	if err != nil {
		return fmt.Errorf("rankCreateTable: %s", err)
	}
	return
}
