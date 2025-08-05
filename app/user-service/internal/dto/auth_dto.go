package dto

type AuthRequestDTO struct {
	Username string `json:"username" binding:"required" example:"john_doe"`
	Password string `json:"password" binding:"required" example:"password"`
}

type RefreshTokenRequestDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"refresh_token"`
}

type TokenResponseDTO struct {
	AccessToken  string `json:"access_token" example:"access_token"`
	RefreshToken string `json:"refresh_token" example:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in" example:"3600"`
}
