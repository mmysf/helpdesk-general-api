package domain

type CreateCompanyRequest struct {
	Name         string        `json:"name"`
	Email        string        `json:"email"`
	LogoAttachId string        `json:"logoAttachId"`
	ColorMode    ColorMode     `json:"colorMode"`
	Domain       CompanyDomain `json:"domain"`
}

type ColorMode struct {
	Light Color `json:"light"`
	Dark  Color `json:"dark"`
}

type Color struct {
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
}

type CompanyDomain struct {
	IsCustom  *bool  `json:"isCustom"`
	Subdomain string `json:"subdomain"`
	FullUrl   string `json:"fullUrl"`
}
