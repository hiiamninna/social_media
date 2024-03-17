package collections

type FileUpload struct {
	ImageUrl string `json:"imageUrl"`
}

type Meta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}
