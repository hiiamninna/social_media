package collections

import "time"

type UserRegisterByEmail struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

type UserRegisterByPhone struct {
	Phone       string `json:"phone"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

type UserLogin struct {
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

type User struct {
	Id       string
	Name     string
	Email    string
	Phone    string
	ImageUrl string
	Password string
}

type UserAsFriend struct {
	UserId      string `json:"userId"`
	Name        string `json:"name"`
	ImageUrl    string `json:"imageUrl"`
	FriendCount int    `json:"friendCount"`
	TCreatedAt  time.Time
	CreatedAt   string `json:"createdAt"`
}

type Friend struct {
	Id      int
	UserId  int
	AddedBy int
}

type FileUpload struct {
	ImageUrl string `json:"imageUrl"`
}

type Meta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type Post struct {
	UserID   string `json:"-"`
	ID       int    `json:"-"`
	PostID   string `json:"postId"`
	PostData struct {
		PostInHtml string    `json:"postInHtml"`
		Tags       []string  `json:"tags"`
		CreatedAt  time.Time `json:"createdAt"`
	} `json:"post"`
	Comments []Comment    `json:"comments"`
	Creator  UserAsFriend `json:"creator"`
}

type Comment struct {
	PostID  string       `json:"-"`
	Comment string       `json:"comment"`
	Creator UserAsFriend `json:"creator"`
}
