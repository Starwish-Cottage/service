package admin

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	FullName     string `json:"full_name"`
	SessionToken string `json:"session_token"`
}
