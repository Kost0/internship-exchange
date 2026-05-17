package model

import "time"

type ApplicationStatus string

const (
	StatusApplied   ApplicationStatus = "applied"
	StatusReviewing ApplicationStatus = "reviewing"
	StatusInterview ApplicationStatus = "interview"
	StatusAccepted  ApplicationStatus = "accepted"
	StatusRejected  ApplicationStatus = "rejected"
)

var validTransitions = map[ApplicationStatus][]ApplicationStatus{
	StatusApplied:   {StatusReviewing, StatusRejected},
	StatusReviewing: {StatusInterview, StatusRejected},
	StatusInterview: {StatusAccepted, StatusRejected},
}

func (s ApplicationStatus) CanTransitionTo(next ApplicationStatus) bool {
	allowed, ok := validTransitions[s]
	if !ok {
		return false
	}

	for _, a := range allowed {
		if a == next {
			return true
		}
	}

	return false
}

type Application struct {
	ID          string
	StudentID   string
	ListingID   string
	CoverLetter string
	Status      ApplicationStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Events      []ApplicationEvent
}

type ApplicationEvent struct {
	ID            string
	ApplicationID string
	OldStatus     ApplicationStatus
	NewStatus     ApplicationStatus
	Comment       string
	ChangedAt     time.Time
}
