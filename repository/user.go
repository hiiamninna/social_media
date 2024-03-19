package repository

import (
	"database/sql"
	"fmt"
	"social_media/collections"
)

type User struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) User {
	return User{
		db: db,
	}
}

func (r User) Create(input collections.UserRegisterInput) (int, error) {

	var id int
	sql := `INSERT INTO users (email, phone, name, password, image_url, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, current_timestamp, current_timestamp) RETURNING id;`
	rows, err := r.db.Query(sql, input.Email, input.Phone, input.Name, input.Password, input.ImageUrl)

	for rows.Next() {
		rows.Scan(&id)
	}
	defer rows.Close()

	return id, fmt.Errorf("create : %s", err.Error())
}

func (r User) GetByEmail(email string) (collections.User, error) {

	user := collections.User{}

	sql := `SELECT TEXT(id), name, email, phone, image_url, password FROM users WHERE UPPER(email) = UPPER($1) and deleted_at is null;`
	err := r.db.QueryRow(sql, email).Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.ImageUrl, &user.Password)
	if err != nil {
		return user, fmt.Errorf("get by email : %s", err.Error())
	}

	return user, nil
}

func (r User) GetByPhone(phone string) (collections.User, error) {

	user := collections.User{}

	sql := `SELECT TEXT(id), name, email, phone, image_url, password FROM users WHERE UPPER(phone) = UPPER($1) and deleted_at is null;`
	err := r.db.QueryRow(sql, phone).Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.ImageUrl, &user.Password)
	if err != nil {
		return user, fmt.Errorf("get by phone : %s", err.Error())
	}

	return user, nil
}

func (r User) GetByID(id string) (collections.User, error) {

	user := collections.User{}

	sql := `SELECT TEXT(id), name, email, phone, image_url, password FROM users WHERE id = $1 and deleted_at is null;`
	err := r.db.QueryRow(sql, id).Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.ImageUrl, &user.Password)
	if err != nil {
		return user, fmt.Errorf("get by id : %s", err.Error())
	}

	return user, nil
}

func (r User) UpdateProfile(input collections.UserUpdateInput) error {

	sql := `UPDATE users SET name = $1, image_url = $2, updated_at = current_timestamp WHERE id = $3 AND deleted_at IS NULL;`
	_, err := r.db.Exec(sql, input.Name, input.ImageUrl, input.UserID)
	if err != nil {
		return fmt.Errorf("update user : %s", err.Error())
	}

	return nil
}

func (r User) UpdateEmail(input collections.UserLinkInput) error {

	sql := `UPDATE users SET email = $1, updated_at = current_timestamp WHERE id = $2 AND deleted_at IS NULL;`
	_, err := r.db.Exec(sql, input.Email, input.UserID)
	if err != nil {
		return fmt.Errorf("update email : %s", err.Error())
	}

	return nil
}

func (r User) UpdatePhone(input collections.UserLinkInput) error {

	sql := `UPDATE users SET phone = $1, updated_at = current_timestamp WHERE id = $2 AND deleted_at IS NULL;`
	_, err := r.db.Exec(sql, input.Phone, input.UserID)
	if err != nil {
		return fmt.Errorf("update phone : %s", err.Error())
	}

	return nil
}

func (r User) List() ([]collections.UserAsFriend, error) {

	users := []collections.UserAsFriend{}
	sql := `SELECT id, name, image_url, 1, created_at FROM users WHERE deleted_at IS NULL;`
	rows, err := r.db.Query(sql)
	if err != nil {
		return users, fmt.Errorf("select user list : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		u := collections.UserAsFriend{}

		err := rows.Scan(&u.UserId, &u.Name, &u.ImageUrl, &u.FriendCount, &u.CreatedAt)
		if err != nil {
			return users, fmt.Errorf("rows scan : %w", err)
		}

		users = append(users, u)
	}

	return users, nil
}
