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

func (r Post) List() ([]collections.Post, error) {

	posts := []collections.Post{}
	sql := `SELECT id, post, tags, created_at, user_id FROM posts WHERE deleted_at IS NULL;`
	rows, err := r.db.Query(sql)
	if err != nil {
		return posts, fmt.Errorf("select post list : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		p := collections.Post{}

		err := rows.Scan(&p.PostID, &p.PostData.PostInHtml, pq.Array(&p.PostData.Tags), &p.PostData.CreatedAt, &p.UserID)
		if err != nil {
			return posts, fmt.Errorf("rows scan : %w", err)
		}

		posts = append(posts, p)
	}

	return posts, nil
}
