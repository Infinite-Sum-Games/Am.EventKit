package models

type TeamBookingRequest struct {
	TeamName    string       `json:"team_name" binding:"required"`
	TeamMembers []TeamMember `json:"team_members" binding:"required"`
}

type TeamMember struct {
	StudentEmail string `json:"student_email" binding:"required"`
	StudentRole  string `json:"student_role" binding:"required"`
}
