package dto

type PaginationResponse struct {
	Total int64 `json:"total"`
	Take  int   `json:"take"`
	Skip  int   `json:"skip"`
}
