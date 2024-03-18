package repository

import (
	"database/sql"
	"fmt"
	"social_media/collections"
)

type Friend struct {
	db *sql.DB
}

func NewFriendRepository(db *sql.DB) User {
	return User{
		db: db,
	}
}

func (r Friend) Create(input collections.FriendInput) error {
	sql := `INSERT INTO friends (added_by, user_id, created_at, updated_at) VALUES ($1, $2, current_timestamp, current_timestamp);`
	_, err := r.db.Exec(sql, input.UserID, input.FriendID)
	if err != nil {
		return fmt.Errorf("create : %s", err.Error())
	}

	// TODO : should we add friend count to the users so, to get we should not always count it ?

	return nil
}

func (r Friend) List(input collections.FriendInputParam) ([]collections.UserAsFriend, error) {

	// TODO : params, query

	friends := []collections.UserAsFriend{}
	sql := `SELECT DISTINCT(f.user_id), u.name, u.image_url, 1, f.created_at
			FROM friends f LEFT JOIN users u ON f.user_id = u.id
			WHERE f.added_by = $1 AND f.deleted_at IS NULL;`

	rows, err := r.db.Query(sql, input.UserID)
	if err != nil {
		return friends, fmt.Errorf("added by logged user : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		f := collections.UserAsFriend{}

		err := rows.Scan(&f.UserId, &f.Name, &f.ImageUrl, &f.FriendCount, &f.CreatedAt)
		if err != nil {
			return friends, fmt.Errorf("rows scan : %w", err)
		}

		friends = append(friends, f)
	}

	sql = `SELECT DISTINCT(f.added_by), u.name, u.image_url, 1, f.created_at
			FROM friends f LEFT JOIN users u ON f.added_by = u.id
			WHERE f.user_id = $1 AND f.deleted_at IS NULL;`

	rows, err = r.db.Query(sql, input.UserID)
	if err != nil {
		return friends, fmt.Errorf("logged user being added : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		f := collections.UserAsFriend{}

		err := rows.Scan(&f.UserId, &f.Name, &f.ImageUrl, &f.FriendCount, &f.CreatedAt)
		if err != nil {
			return friends, fmt.Errorf("rows scan : %w", err)
		}

		friends = append(friends, f)
	}

	return friends, nil
}

func (r Friend) GetByUser(input collections.FriendInput) (collections.Friend, error) {

	friend := collections.Friend{}
	sql := `SELECT id, user_id, added_by FROM friends WHERE ((added_by = $1 AND user_id = $2) OR (user_id = $3 AND added_by = $4)) AND deleted_at IS NULL;`
	err := r.db.QueryRow(sql, input.UserID, input.FriendID, input.UserID, input.FriendID).Scan(&friend.Id, &friend.UserId, &friend.AddedBy)
	if err != nil {
		return friend, fmt.Errorf("select by user : %s", err.Error())
	}

	return friend, nil
}

func (r Friend) Delete(id int) error {
	sql := `UPDATE friends SET deleted_at = current_timestamp WHERE id = $1 AND deleted_at IS NULL;`
	_, err := r.db.Exec(sql, id)
	if err != nil {
		return fmt.Errorf("delete : %s", err.Error())
	}

	return nil
}
