package epgc

import (
	"database/sql"
	"log"
)

// Post - struct for post
type Post struct {
	ID   int64  `sql:"id" json:"id"`
	Name string `sql:"name" json:"name"`
	GO   bool   `sql:"go" json:"go"`
	Note string `sql:"note, null" json:"note"`
}

func scanPost(row *sql.Row) (Post, error) {
	var (
		sid   sql.NullInt64
		sname sql.NullString
		sgo   sql.NullBool
		snote sql.NullString
		post  Post
	)
	err := row.Scan(&sid, &sname, &sgo, &snote)
	if err != nil {
		log.Println("scanPost row.Scan ", err)
		return post, err
	}
	post.ID = n2i(sid)
	post.Name = n2s(sname)
	post.GO = n2b(sgo)
	post.Note = n2s(snote)
	return post, nil
}

func scanPosts(rows *sql.Rows, opt string) ([]Post, error) {
	var posts []Post
	for rows.Next() {
		var (
			sid   sql.NullInt64
			sname sql.NullString
			sgo   sql.NullBool
			snote sql.NullString
			post  Post
		)
		switch opt {
		case "list":
			err := rows.Scan(&sid, &sname, &sgo, &snote)
			if err != nil {
				log.Println("scanPosts rows.Scan list ", err)
				return posts, err
			}
			post.Name = n2s(sname)
			post.GO = n2b(sgo)
			post.Note = n2s(snote)
		case "select":
			err := rows.Scan(&sid, &sname)
			if err != nil {
				log.Println("scanPosts rows.Scan select ", err)
				return posts, err
			}
			post.Name = n2s(sname)
			// if len(post.Name) > 210 {
			// 	post.Name = post.Name[0:210]
			// }
		}
		post.ID = n2i(sid)
		posts = append(posts, post)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanPosts rows.Err ", err)
	}
	return posts, err
}

// GetPost - get one post by id
func (e *Edb) GetPost(id int64) (Post, error) {
	if id == 0 {
		return Post{}, nil
	}
	row := e.db.QueryRow("SELECT id,name,go,note FROM posts WHERE id = $1", id)
	post, err := scanPost(row)
	return post, err
}

// GetPostList - get all post for list
func (e *Edb) GetPostList() ([]Post, error) {
	rows, err := e.db.Query("SELECT id,name,go,note FROM posts ORDER BY name ASC")
	if err != nil {
		log.Println("GetPostList e.db.Query ", err)
		return []Post{}, err
	}
	posts, err := scanPosts(rows, "list")
	return posts, err
}

// GetPostSelect - get all post for select
func (e *Edb) GetPostSelect(g bool) ([]Post, error) {
	rows, err := e.db.Query("SELECT id,name FROM posts WHERE go=$1 ORDER BY name ASC", g)
	if err != nil {
		log.Println("GetPostSelect e.db.Query ", err)
		return []Post{}, err
	}
	posts, err := scanPosts(rows, "select")
	return posts, err
}

// CreatePost - create new post
func (e *Edb) CreatePost(post Post) (int64, error) {
	stmt, err := e.db.Prepare(`INSERT INTO posts(name, go, note) VALUES($1, $2, $3) RETURNING id`)
	if err != nil {
		log.Println("CreatePost e.db.Prepare ", err)
		return 0, err
	}
	err = stmt.QueryRow(s2n(post.Name), post.GO, s2n(post.Note)).Scan(&post.ID)
	if err != nil {
		log.Println("CreatePost db.QueryRow ", err)
	}
	return post.ID, err
}

// UpdatePost - save post changes
func (e *Edb) UpdatePost(s Post) error {
	stmt, err := e.db.Prepare("UPDATE posts SET name=$2,note=$3 WHERE id = $1")
	if err != nil {
		log.Println("UpdatePost e.db.Prepare ", err)
		return err
	}
	_, err = stmt.Exec(i2n(s.ID), s2n(s.Name), s2n(s.Note))
	if err != nil {
		log.Println("UpdatePost stmt.Exec ", err)
	}
	return err
}

// DeletePost - delete post by id
func (e *Edb) DeletePost(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec("DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		log.Println("DeletePost e.db.Exec ", err)
	}
	return err
}

func (e *Edb) postCreateTable() error {
	str := `CREATE TABLE IF NOT EXISTS posts (id BIGSERIAL PRIMARY KEY, name TEXT, go BOOL NOT NULL DEFAULT FALSE, note TEXT, UNIQUE (name, go))`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("postCreateTable e.db.Exec ", err)
	}
	return err
}
