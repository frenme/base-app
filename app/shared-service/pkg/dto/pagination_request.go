package dto

type PaginationRequest struct {
	Take int `form:"take" example:"10"`
	Skip int `form:"skip" example:"0"`
}
