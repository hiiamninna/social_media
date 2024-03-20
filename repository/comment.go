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

func (r Comment) List(postIds []int) ([]collections.Comment, error) {

	comments := []collections.Comment{}
	sql := `SELECT c.comment, u.id, u.name, u.image_url, 1, c.created_at, c.post_id FROM comments c LEFT JOIN users u ON c.user_id = u.id WHERE c.deleted_at IS NULL and u.deleted_at IS NULL 
	`

	var values []interface{}
	if len(postIds) > 0 {
		sql += "AND c.post_id in("
		for i, id := range postIds {
			// and c.post_id in ([post_ids])
			sql += fmt.Sprintf("$%d", (i + 1))
			if i+1 != len(postIds) {
				sql += ","
			}
			values = append(values, id)
		}
		sql += ")"
	}

	sql += ";"

	rows, err := r.db.Query(sql, values...)
	if err != nil {
		return comments, fmt.Errorf("select comment list : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		c := collections.Comment{}

		err := rows.Scan(&c.Comment, &c.Creator.UserId, &c.Creator.Name, &c.Creator.ImageUrl, &c.Creator.FriendCount, &c.Creator.CreatedAt, &c.PostID)
		if err != nil {
			return comments, fmt.Errorf("rows scan : %w", err)
		}

		comments = append(comments, c)
	}

	return comments, nil
}
