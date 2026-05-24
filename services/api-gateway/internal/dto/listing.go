package dto

type ListingSkillResponse struct {
	ID         string `json:"id"`
	ListingID  string `json:"listingId"`
	Skill      string `json:"skill"`
	IsRequired bool   `json:"isRequired"`
}

type CompanyInfoResponse struct {
	ID       string `json:"id"`
	UserID   string `json:"userId"`
	Name     string `json:"name"`
	LogoURL  string `json:"logoUrl"`
	Industry string `json:"industry"`
	City     string `json:"city"`
}

type ListingResponse struct {
	ID             string                 `json:"id"`
	CompanyID      string                 `json:"companyId"`
	Company        *CompanyInfoResponse   `json:"company"`
	Title          string                 `json:"title"`
	Description    string                 `json:"description"`
	Requirements   string                 `json:"requirements"`
	WhatWeOffer    string                 `json:"whatWeOffer"`
	City           string                 `json:"city"`
	Format         string                 `json:"format"`
	EmploymentType string                 `json:"employmentType"`
	SalaryFrom     int64                  `json:"salaryFrom"`
	SalaryTo       int64                  `json:"salaryTo"`
	SalaryCurrency string                 `json:"salaryCurrency"`
	Deadline       string                 `json:"deadline"`
	Status         string                 `json:"status"`
	Skills         []ListingSkillResponse `json:"skills"`
	CreatedAt      string                 `json:"createdAt"`
	UpdatedAt      string                 `json:"updatedAt"`
}

type GetListingsResponse struct {
	Items []ListingResponse `json:"items"`
	Total int64             `json:"total"`
	Page  int32             `json:"page"`
	Limit int32             `json:"limit"`
}
