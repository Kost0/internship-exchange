package model

import "time"

type ListingStatus string
type ListingFormat string
type EmploymentType string

const (
	StatusDraft  ListingStatus = "draft"
	StatusActive ListingStatus = "active"
	StatusClosed ListingStatus = "closed"

	FormatOffice ListingFormat = "office"
	FormatRemote ListingFormat = "remote"
	FormatHybrid ListingFormat = "hybrid"

	EmploymentFullTime EmploymentType = "full_time"
	EmploymentPartTime EmploymentType = "part_time"
	EmploymentProject  EmploymentType = "project"
)

type Company struct {
	ID        string
	UserID    string
	Name      string
	LogoURL   string
	Industry  string
	City      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Listing struct {
	ID             string
	CompanyID      string
	Company        *Company
	Title          string
	Description    string
	Requirements   string
	WhatWeOffer    string
	City           string
	Format         ListingFormat
	EmploymentType EmploymentType
	SalaryFrom     int64
	SalaryTo       int64
	SalaryCurrency string
	Deadline       string
	Status         ListingStatus
	Skills         []ListingSkill
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type ListingSkill struct {
	ID         string
	ListingID  string
	Skill      string
	IsRequired bool
}

type ListingsFilter struct {
	Query          string
	Format         string
	EmploymentType string
	City           string
	Skill          string
	Page           int32
	Limit          int32
}
