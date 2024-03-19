package collections

type UserRegisterInput struct {
	CredType  string `json:"credentialType"`
	CredValue string `json:"credentialValue"`
	Email     string
	Phone     string
	Name      string `json:"name"`
	Password  string `json:"password"`
	ImageUrl  string
}

type UserLoginInput struct {
	CredType  string `json:"credentialType"`
	CredValue string `json:"credentialValue"`
	Password  string `json:"password"`
}

type UserLinkInput struct {
	UserID string
	Email  string `json:"email"`
	Phone  string `json:"phone"`
}

type UserUpdateInput struct {
	UserID   string
	ImageUrl string `json:"imageUrl"`
	Name     string `json:"name"`
}

type FriendInputParam struct {
	UserID     string
	Search     string `query:"search"`
	OnlyFriend bool   `query:"onlyFriend"`
	OrderBy    string `query:"orderBy"`
	SortBy     string `query:"sortBy"`
	Limit      int    `query:"limit"`
	Offset     int    `query:"offset"`
}

type FriendInput struct {
	UserID   string
	FriendID string `json:"userId"`
}

type PostInput struct {
	UserID string
	Post   string   `json:"postInHtml"`
	Tags   []string `json:"tags"`
}

type PostInputParam struct {
	UserID string
	Tags   []string `query:"searchTag"`
	Search string   `query:"search"`
	Limit  int      `query:"limit"`
	Offset int      `query:"offset"`
}

type CommentInput struct {
	UserID  string
	PostID  string `json:"postId"`
	Comment string `json:"comment"`
}
