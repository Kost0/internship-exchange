package dto

type StudentProfileResponse struct {
	ID           string               `json:"id"`
	UserID       string               `json:"userId"`
	FirstName    string               `json:"firstName"`
	LastName     string               `json:"lastName"`
	Phone        string               `json:"phone"`
	City         string               `json:"city"`
	Bio          string               `json:"bio"`
	AvatarURL    string               `json:"avatarUrl"`
	ResumeURL    string               `json:"resumeUrl"`
	GithubURL    string               `json:"githubUrl"`
	LinkedinURL  string               `json:"linkedinUrl"`
	PortfolioURL string               `json:"portfolioUrl"`
	Educations   []EducationResponse  `json:"educations"`
	Experiences  []ExperienceResponse `json:"experiences"`
	Projects     []ProjectResponse    `json:"projects"`
	Skills       []SkillResponse      `json:"skills"`
	Languages    []LanguageResponse   `json:"languages"`
}

type EducationResponse struct {
	ID             string  `json:"id"`
	StudentID      string  `json:"studentId"`
	University     string  `json:"university"`
	Faculty        string  `json:"faculty"`
	Specialization string  `json:"specialization"`
	Degree         string  `json:"degree"`
	StartYear      int32   `json:"startYear"`
	EndYear        int32   `json:"endYear"`
	GPA            float64 `json:"gpa"`
	IsCurrent      bool    `json:"isCurrent"`
}

type ExperienceResponse struct {
	ID          string `json:"id"`
	StudentID   string `json:"studentId"`
	CompanyName string `json:"companyName"`
	Position    string `json:"position"`
	Description string `json:"description"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	IsCurrent   bool   `json:"isCurrent"`
	Format      string `json:"format"`
}

type ProjectResponse struct {
	ID          string   `json:"id"`
	StudentID   string   `json:"studentId"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Techs       []string `json:"techs"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
}

type SkillResponse struct {
	ID        string `json:"id"`
	StudentID string `json:"studentId"`
	Skill     string `json:"skill"`
	Level     string `json:"level"`
}

type LanguageResponse struct {
	ID        string `json:"id"`
	StudentID string `json:"studentId"`
	Language  string `json:"language"`
	Level     string `json:"level"`
}

type CompanyProfileResponse struct {
	ID               string `json:"id"`
	UserID           string `json:"userId"`
	Name             string `json:"name"`
	Tagline          string `json:"tagline"`
	Description      string `json:"description"`
	Industry         string `json:"industry"`
	Size             string `json:"size"`
	FoundedYear      int32  `json:"foundedYear"`
	Website          string `json:"website"`
	ContactEmail     string `json:"contactEmail"`
	City             string `json:"city"`
	Country          string `json:"country"`
	IsRemoteFriendly bool   `json:"isRemoteFriendly"`
	LogoURL          string `json:"logoUrl"`
}
