package domain

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type ChangeDomainRequest struct {
	IsCustom  *bool  `json:"isCustom"`
	Subdomain string `json:"subdomain"`
	FullUrl   string `json:"fullUrl"`
}

type UpdateProfileRequest struct {
	Name     string `json:"name"`
	Bio      string `json:"bio"`
	Contact  string `json:"contact"`
	AttachId string `json:"attachId"`
}

type ChangeColorMode struct {
	Light Color `json:"light"`
	Dark  Color `json:"dark"`
}
