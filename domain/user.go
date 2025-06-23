package domain

type CreateUserRequest struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	JobTitle         string `json:"jobTitle"`
	CompanyProductId string `json:"companyProductId"`
	Role             string `json:"role"`
}
