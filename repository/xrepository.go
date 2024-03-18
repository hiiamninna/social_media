package repository

import (
	"database/sql"
)

type Repository struct {
	User   User
	Friend Friend
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		User:   NewUserRepository(db),
		Friend: Friend(NewFriendRepository(db)),
	}
}
