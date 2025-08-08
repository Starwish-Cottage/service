package admin

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Success      bool   `json:"success" binding:"required"`
	FullName     string `json:"full_name" binding:"required"`
	SessionToken string `json:"session_token" binding:"required"`
	Message      string `json:"message" biding:"required"`
}

type VerifySessionRequest struct {
	SessionToken string `json:"session_token" binding:"required"`
}

type VerifySessionResponse struct {
	Valid   bool   `json:"valid" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type UploadImageResponse struct {
	Success   bool     `json:"success" binding:"required"`
	ImageUrls []string `json:"image_urls" binding:"required"`
	Message   string   `json:"message" binding:"required"`
}
