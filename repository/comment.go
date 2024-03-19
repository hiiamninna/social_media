package repository

import (
	"database/sql"
	"fmt"
	"social_media/collections"
)

type Comment struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) Comment {
	return Comment{
		db: db,
	}
}

func (r Comment) Create(input collections.CommentInput) error {
	sql := `INSERT INTO comments (user_id, post_id, comment, created_at, updated_at) VALUES ($1, $2, $3, current_timestamp, current_timestamp);`
	_, err := r.db.Exec(sql, input.UserID, input.PostID, input.Comment)
	if err != nil {
		return fmt.Errorf("create : %s", err.Error())
	}

	return nil
}
