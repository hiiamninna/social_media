package repository

import (
	"database/sql"
	"fmt"
	"social_media/collections"
	"time"
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

func (r Comment) List(postIds []int) (map[int][]collections.Comment, error) {

	result := make(map[int][]collections.Comment)

	if len(postIds) > 0 {
		for _, p := range postIds {
			sql := fmt.Sprintf(`SELECT c.comment, u.id, u.name, u.image_url, u.total_friend, c.created_at, c.post_id FROM comments c 
					LEFT JOIN users u ON c.user_id = u.id WHERE c.deleted_at IS NULL AND c.post_id = %d ORDER BY c.created_at DESC;`, p)
			rows, err := r.db.Query(sql)
			if err != nil {
				return result, fmt.Errorf("select comment list : %w", err)
			}
			defer rows.Close()

			comments := []collections.Comment{}
			for rows.Next() {
				c := collections.Comment{}

				err := rows.Scan(&c.Comment, &c.Creator.UserId, &c.Creator.Name, &c.Creator.ImageUrl, &c.Creator.FriendCount, &c.CreatedAt, &c.PostID)
				if err != nil {
					return result, fmt.Errorf("rows scan : %w", err)
				}

				// TODO : check time format ISO 8601
				c.CreatedAtStr = c.CreatedAt.Format(time.RFC3339)
				comments = append(comments, c)
			}

			result[p] = comments
		}
	}

	return result, nil
}
