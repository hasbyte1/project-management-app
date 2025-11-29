package dto

type OrganizationDTO struct {
	ID          string  `json:"id"`
	ParentID    *string `json:"parent_id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	LogoURL     *string `json:"logo_url"`
	Depth       int     `json:"depth"`
	Path        *string `json:"path"`
	Settings    []byte  `json:"settings"`
	CreatedBy   string  `json:"created_by"`
}
