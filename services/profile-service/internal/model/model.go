package model

import "time"

type Student struct {
	ID           string
	UserID       string
	FirstName    string
	LastName     string
	Phone        string
	City         string
	Bio          string
	AvatarURL    string
	ResumeURL    string
	GithubURL    string
	LinkedinURL  string
	PortfolioURL string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Educations   []Education
	Experiences  []Experience
	Projects     []Project
	Skills       []StudentSkill
	Languages    []StudentLanguage
}

type Education struct {
	ID             string
	StudentID      string
	University     string
	Faculty        string
	Specialization string
	Degree         string
	StartYear      int32
	EndYear        int32
	GPA            float64
	IsCurrent      bool
}

type Experience struct {
	ID          string
	StudentID   string
	CompanyName string
	Position    string
	Description string
	StartDate   string
	EndDate     string
	IsCurrent   bool
	Format      string
}

type Project struct {
	ID          string
	StudentID   string
	Title       string
	Description string
	URL         string
	Techs       []string
	StartDate   string
	EndDate     string
}

type StudentSkill struct {
	ID        string
	StudentID string
	Skill     string
	Level     string
}

type StudentLanguage struct {
	ID        string
	StudentID string
	Language  string
	Level     string
}

type Company struct {
	ID               string
	UserID           string
	Name             string
	Tagline          string
	Description      string
	Industry         string
	Size             string
	FoundedYear      int32
	Website          string
	ContactEmail     string
	City             string
	Country          string
	IsRemoteFriendly bool
	LogoURL          string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
