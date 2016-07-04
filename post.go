package epdc

import "fmt"

// Post - struct for post
type Post struct {
	TableName struct{} `sql:"posts"`
	ID        int64    `sql:"id" json:"id"`
	Name      string   `sql:"name" json:"name"`
	GO        bool     `sql:"go" json:"go"`
	Note      string   `sql:"note, null" json:"note"`
}

// GetPost - get one post dy id
func (e *EDc) GetPost(id int64) (post Post, err error) {
	if id == 0 {
		return post, nil
	}
	err = e.db.Model(&post).Where("id = ?", id).Select()
	if err != nil {
		return post, fmt.Errorf("GetPost: %s", err)
	}
	return
}

// GetPostAll - get all post
func (e *EDc) GetPostAll() (posts []Post, err error) {
	err = e.db.Model(&posts).Order("name ASC").Select()
	if err != nil {
		return posts, fmt.Errorf("GetPostAll: %s", err)
	}
	return
}

// GetPostNoGOAll - get all post with no go
func (e *EDc) GetPostNoGOAll() (posts []Post, err error) {
	_, err = e.db.Query(&posts, "SELECT * FROM posts WHERE go = ?", false)
	if err != nil {
		return posts, fmt.Errorf("GetPostNoGOAll: %s", err)
	}
	return
}

// GetPostGOAll - get all post with go
func (e *EDc) GetPostGOAll() (posts []Post, err error) {
	_, err = e.db.Query(&posts, "SELECT * FROM posts WHERE go = ?", true)
	if err != nil {
		return posts, fmt.Errorf("GetPostGOAll: %s", err)
	}
	return
}

// CreatePost - create new post
func (e *EDc) CreatePost(post Post) (err error) {
	err = e.db.Create(&post)
	if err != nil {
		return fmt.Errorf("CreatePost: %s", err)
	}
	return
}

// UpdatePost - save post changes
func (e *EDc) UpdatePost(post Post) (err error) {
	err = e.db.Update(&post)
	if err != nil {
		return fmt.Errorf("UpdatePost: %s", err)
	}
	return
}

// DeletePost - delete post by id
func (e *EDc) DeletePost(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec("DELETE FROM posts WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("DeletePost: %s", err)
	}
	return nil
}

func (e *EDc) postCreateTable() (err error) {
	str := `CREATE TABLE IF NOT EXISTS posts (id BIGSERIAL PRIMARY KEY, name TEXT, go BOOL NOT NULL DEFAULT FALSE, note TEXT, UNIQUE (name, go))`
	_, err = e.db.Exec(str)
	if err != nil {
		return fmt.Errorf("postCreateTable: %s", err)
	}
	return
}
