package dto

type ApplicationEventResponse struct {
	ID            string `json:"id"`
	ApplicationID string `json:"applicationId"`
	OldStatus     string `json:"oldStatus"`
	NewStatus     string `json:"newStatus"`
	Comment       string `json:"comment"`
	ChangedAt     string `json:"changedAt"`
}

type ApplicationResponse struct {
	ID          string                     `json:"id"`
	StudentID   string                     `json:"studentId"`
	ListingID   string                     `json:"listingId"`
	CoverLetter string                     `json:"coverLetter"`
	Status      string                     `json:"status"`
	CreatedAt   string                     `json:"createdAt"`
	UpdatedAt   string                     `json:"updatedAt"`
	Student     *StudentProfileResponse    `json:"student,omitempty"`
	Listing     *ListingResponse           `json:"listing,omitempty"`
	Events      []ApplicationEventResponse `json:"events"`
}
