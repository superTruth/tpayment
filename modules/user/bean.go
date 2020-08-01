package user

type LoginRequest struct {
	Email     string `json:"email"`
	Pwd       string `json:"pwd"`
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
	Role  string `json:"role,omitempty"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AddUserRequest struct {
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
	Role  string `json:"role"`
	Name  string `json:"name"`
}
