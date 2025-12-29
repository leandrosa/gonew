package models

type Article struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

type ErrorResponse struct {
	Description string `json:"message"`
}

func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{
		Description: err.Error(),
	}
}
