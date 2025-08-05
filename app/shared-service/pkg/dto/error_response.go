package dto

type ErrorResponse struct {
	StatusCode int    `json:"statusCode" example:"400"`
	Message    string `json:"message" example:"Some error message"`
}
