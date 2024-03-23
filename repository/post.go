package repository

import (
	"database/sql"
	"fmt"
	"social_media/collections"
	"strings"
	"time"

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

func (r Post) List(input collections.PostInputParam) ([]collections.Post, []int, collections.Meta, error) {

	posts := []collections.Post{}
	meta := collections.Meta{
		Limit:  input.Limit,
		Offset: input.Offset,
		Total:  0,
	}
	var ids []int
	var values []interface{}
	counter, query := 1, ""

	sql := `SELECT p.id, TEXT(p.id), p.post, p.tags, p.created_at, TEXT(p.user_id), TEXT(u.id), u.name, u.image_url, u.total_friend, u.created_at FROM posts p 
			LEFT JOIN users u ON u.id = p.user_id
			WHERE (p.user_id = ` + input.UserID + ` OR p.user_id IN (SELECT user_id AS id from friends WHERE deleted_at is null AND added_by = ` + input.UserID + ` UNION SELECT added_by AS id from friends WHERE deleted_at is null AND user_id = ` + input.UserID + `)) AND p.deleted_at is null [query]
			ORDER BY p.created_at DESC [pagination];`

	if input.Search != "" {
		query += fmt.Sprintf("AND p.post ~* $%d ", counter)
		values = append(values, input.Search)
		counter += 1
	}

	if len(input.Tags) > 0 {
		query += " AND ("
		for i, v := range input.Tags {
			query += "$" + fmt.Sprint(counter) + " = ANY(p.tags)"
			values = append(values, v)
			counter += 1
			if i+1 < len(input.Tags) {
				query += " OR "
			}
		}
		query += ")"
	}

	sql = strings.ReplaceAll(sql, "[query]", query)
	sql = strings.ReplaceAll(sql, "[pagination]", fmt.Sprintf("limit %d offset %d", input.Limit, input.Offset))

	rows, err := r.db.Query(sql, values...)
	if err != nil {
		return posts, ids, meta, fmt.Errorf("select post list : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		p := collections.Post{}

		err := rows.Scan(&p.ID, &p.PostID, &p.PostData.PostInHtml, pq.Array(&p.PostData.Tags), &p.PostData.CreatedAt, &p.UserID, &p.Creator.UserId, &p.Creator.Name, &p.Creator.ImageUrl, &p.Creator.FriendCount, &p.Creator.CreatedAt)
		if err != nil {
			return posts, ids, meta, fmt.Errorf("row scan : %w", err)
		}

		p.PostData.CreatedAtStr = p.PostData.CreatedAt.Format(time.RFC3339)
		p.Creator.CreatedAtStr = p.Creator.CreatedAt.Format(time.RFC3339)

		posts = append(posts, p)
		ids = append(ids, p.ID)
	}

	meta.Total, err = r.ListCount(input)
	if err != nil {
		fmt.Println(time.Now().Format("2006-01-02 15:01:02 "), "list count : "+err.Error())
	}

	return posts, ids, meta, nil
}

func (r Post) ListCount(input collections.PostInputParam) (int, error) {

	var total int
	var values []interface{}
	counter, query := 1, ""

	sql := `SELECT COUNT(*) FROM posts p WHERE (p.user_id = ` + input.UserID + ` OR p.user_id IN (SELECT user_id AS id from friends WHERE deleted_at is null AND added_by = ` + input.UserID + ` UNION SELECT added_by AS id from friends WHERE deleted_at is null AND user_id = ` + input.UserID + `)) AND p.deleted_at is null [query];`

	if input.Search != "" {
		query += fmt.Sprintf("AND p.post ~* $%d ", counter)
		values = append(values, input.Search)
		counter += 1
	}

	if len(input.Tags) > 0 {
		query += " AND ("
		for i, v := range input.Tags {
			query += "$" + fmt.Sprint(counter) + " = ANY(p.tags)"
			values = append(values, v)
			counter += 1
			if i+1 < len(input.Tags) {
				query += " OR "
			}
		}
		query += ")"
	}

	sql = strings.ReplaceAll(sql, "[query]", query)
	rows, err := r.db.Query(sql, values...)
	if err != nil {
		return 0, fmt.Errorf("select post list count : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&total)
		if err != nil {
			return total, fmt.Errorf("rows scan : %w", err)
		}
	}

	return total, nil
}

func (r Post) GetById(id string) (collections.PostData, error) {

	post := collections.PostData{}
	sql := `SELECT TEXT(id), post, tags, created_at, user_id FROM posts WHERE TEXT(id) = $1 AND deleted_at IS NULL;`
	err := r.db.QueryRow(sql, id).Scan(&post.ID, &post.Post, pq.Array(&post.Tags), &post.CreatedAt, &post.UserID)
	if err != nil {
		return post, fmt.Errorf("get post by id : %w", err)
	}

	return post, nil
}
