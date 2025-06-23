package domain

type LoginRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	AccessKey string `json:"accessKey"`
}

type RegisterRequest struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	AccessKey        string `json:"accessKey"`
	CompanyProductId string `json:"companyProductId"`
}

type VerifyRegisterRequest struct {
	Token string `json:"token"`
}

type EmailPasswordResetRequest struct {
	Email string `json:"email"`
}

type PasswordResetRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type SuperuserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SuperadminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AccountRequest struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	CompanyID        string `json:"companyId"`
	CompanyProductID string `json:"companyProductId"`
}
