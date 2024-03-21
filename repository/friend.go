package repository

import (
	"database/sql"
	"fmt"
	"social_media/collections"
	"strings"
	"time"
)

type Friend struct {
	db *sql.DB
}

func NewFriendRepository(db *sql.DB) Friend {
	return Friend{
		db: db,
	}
}

func (r Friend) Create(input collections.FriendInput) error {

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("init : %s", err.Error())
	}

	sql := `INSERT INTO friends (added_by, user_id, created_at, updated_at) VALUES ($1, $2, current_timestamp, current_timestamp);`
	_, err = tx.Exec(sql, input.UserID, input.FriendID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("create : %s", err.Error())
	}

	sql = `UPDATE users SET total_friend = total_friend + 1 WHERE id in ($1, $2)`
	_, err = tx.Exec(sql, input.UserID, input.FriendID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("add counter total_friend : %s", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("commit : %s", err.Error())
	}

	return nil
}

func (r Friend) List(input collections.FriendInputParam) ([]collections.UserAsFriend, error) {

	var filter, order string
	var values []interface{}

	valuesCount := 1

	sql := `SELECT u.id, u.name, u.image_url, u.total_friend, u.created_at u FROM users u [join] WHERE u.deleted_at is null [filter] [order] [limit] ;`

	if input.OnlyFriend {
		sql = strings.ReplaceAll(sql, "[join]", fmt.Sprintf(`INNER JOIN (
			SELECT user_id AS friend_id FROM friends WHERE added_by = $%d and deleted_at is null
			UNION
			SELECT added_by AS friend_id FROM friends WHERE user_id = $%d and deleted_at is null
		) AS f ON u.id = f.friend_id`, 1, 2))
		values = append(values, input.UserID, input.UserID)
		valuesCount = 3
	} else {
		sql = strings.ReplaceAll(sql, "[join]", "")
	}

	if input.Search != "" {
		filter += " AND u.name ~* $" + fmt.Sprint(valuesCount)
		values = append(values, input.Search)
		valuesCount++
	}

	switch input.SortBy {
	case "friendCount":
		order += " ORDER BY u.total_friend "
	default:
		order += " ORDER BY u.created_at "
	}

	switch input.OrderBy {
	case "asc":
		order += " ASC "
	default:
		order += " DESC "
	}

	sql = strings.ReplaceAll(sql, "[filter]", filter)
	sql = strings.ReplaceAll(sql, "[order]", order)

	if input.Limit != 0 && input.Offset != 0 {
		sql = strings.ReplaceAll(sql, "[limit]", fmt.Sprintf(" LIMIT %d OFFSET %d ", input.Limit, input.Offset))
	} else {
		sql = strings.ReplaceAll(sql, "[limit]", " LIMIT 5 OFFSET 0 ")
		input.Limit = 5
		input.Offset = 0
	}

	friends := []collections.UserAsFriend{}

	fmt.Println(sql)

	rows, err := r.db.Query(sql, values...)
	if err != nil {
		return friends, fmt.Errorf("exec query : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		f := collections.UserAsFriend{}

		err := rows.Scan(&f.UserId, &f.Name, &f.ImageUrl, &f.FriendCount, &f.CreatedAt)
		if err != nil {
			return friends, fmt.Errorf("rows scan : %w", err)
		}

		f.CreatedAtStr = f.CreatedAt.Format(time.RFC3339)

		friends = append(friends, f)
	}

	return friends, nil
}

func (r Friend) CountList(input collections.FriendInputParam) (int, error) {

	var filter string
	var values []interface{}

	valuesCount := 1

	sql := `SELECT COUNT(*) FROM users u [join] WHERE u.deleted_at is null [filter] ;`

	if input.OnlyFriend {
		sql = strings.ReplaceAll(sql, "[join]", fmt.Sprintf(`INNER JOIN (
			SELECT user_id AS friend_id FROM friends WHERE added_by = $%d and deleted_at is null
			UNION
			SELECT added_by AS friend_id FROM friends WHERE user_id = $%d and deleted_at is null
		) AS f ON u.id = f.friend_id`, 1, 2))
		values = append(values, input.UserID, input.UserID)
		valuesCount = 3
	} else {
		sql = strings.ReplaceAll(sql, "[join]", "")
	}

	if input.Search != "" {
		filter += " AND u.name ~* $" + fmt.Sprint(valuesCount)
		values = append(values, input.Search)
		valuesCount++
	}

	sql = strings.ReplaceAll(sql, "[filter]", filter)

	rows, err := r.db.Query(sql, values...)
	if err != nil {
		return 0, fmt.Errorf("exec query : %w", err)
	}
	defer rows.Close()

	var totalRow int
	for rows.Next() {
		err := rows.Scan(&totalRow)
		if err != nil {
			return 0, fmt.Errorf("rows scan : %w", err)
		}
	}

	return totalRow, nil
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

func (r Friend) Delete(id int, userID, friendID string) error {

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("init : %s", err.Error())
	}

	sql := `UPDATE friends SET deleted_at = current_timestamp WHERE id = $1 AND deleted_at IS NULL;`
	_, err = r.db.Exec(sql, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("delete : %s", err.Error())
	}

	sql = `UPDATE users SET total_friend = total_friend - 1 WHERE id in ($1, $2)`
	_, err = tx.Exec(sql, userID, friendID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("deduct counter total_friend : %s", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("commit : %s", err.Error())
	}

	return nil
}
