package collections

type UserRegisterInput struct {
	CredType  string `json:"credentialType" validate:"required,oneof=phone email"`
	CredValue string `json:"credentialValue" validate:"required"`
	Email     string
	Phone     string
	Name      string `json:"name" validate:"required,min=5,max=50"`
	Password  string `json:"password" validate:"required,min=5,max=15"`
	ImageUrl  string
}

type UserLoginInput struct {
	CredType  string `json:"credentialType" validate:"required,oneof=phone email"`
	CredValue string `json:"credentialValue" validate:"required"`
	Password  string `json:"password" validate:"required,min=5,max=15"`
}

type UserLinkEmail struct {
	UserID string
	Email  string `json:"email" validate:"required,email"`
}

type UserLinkPhone struct {
	UserID string
	Phone  string `json:"phone" validate:"required,startswith=+,min=7,max=13"`
}

type UserUpdateInput struct {
	UserID   string
	ImageUrl string `json:"imageUrl" validate:"required,url"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
}

type FriendInputParam struct {
	UserID        string
	Search        string `query:"search"`
	OnlyFriendStr string `query:"onlyFriend" validate:"required,oneof=true false"`
	OnlyFriend    bool
	OrderBy       string `query:"orderBy" validate:"required,oneof=asc desc"`
	SortBy        string `query:"sortBy" validate:"required,oneof=friendCount createdAt"`
	Limit         int    `query:"limit" validate:"min=0"`
	Offset        int    `query:"offset" validate:"min=0"`
}

type FriendInput struct {
	UserID   string
	FriendID string `json:"userId" validate:"required"`
}

type PostInput struct {
	UserID string
	Post   string   `json:"postInHtml" validate:"required,min=2,max=500"`
	Tags   []string `json:"tags" validate:"required"`
}

type PostInputParam struct {
	UserID string
	Tags   []string `query:"searchTag"`
	Search string   `query:"search"`
	Limit  int      `query:"limit" validate:"min=0"`
	Offset int      `query:"offset" validate:"min=0"`
}

type CommentInput struct {
	UserID  string
	PostID  string `json:"postId" validate:"required"`
	Comment string `json:"comment" validate:"required,min=2,max=500"`
}
