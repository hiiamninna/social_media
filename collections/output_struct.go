package collections

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

type FileUpload struct {
	ImageUrl string `json:"imageUrl"`
}

type Meta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}
