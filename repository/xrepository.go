package repository

import (
	"database/sql"
)

type Repository struct {
	User    User
	Friend  Friend
	Post    Post
	Comment Comment
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		User:    NewUserRepository(db),
		Friend:  NewFriendRepository(db),
		Post:    NewPostRepository(db),
		Comment: NewCommentRepository(db),
	}
}
