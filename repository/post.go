package repository

import (
	"database/sql"
	"fmt"
	"social_media/collections"

	"github.com/lib/pq"
)

type Post struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) Post {
	return Post{
		db: db,
	}
}

func (r Post) Create(input collections.PostInput) error {
	sql := `INSERT INTO posts (user_id, post, tags, created_at, updated_at) VALUES ($1, $2, $3, current_timestamp, current_timestamp);`
	_, err := r.db.Exec(sql, input.UserID, input.Post, pq.Array(input.Tags))
	if err != nil {
		return fmt.Errorf("create : %s", err.Error())
	}

	return nil
}
