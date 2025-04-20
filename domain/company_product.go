package domain

type CreateCompanyProductRequest struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	Code             string `json:"code"`
	SubscriptionType string `json:"subscriptionType"`
	CompanyId        string `json:"companyId"`
	LogoAttachId     string `json:"logoAttachId"`
}
