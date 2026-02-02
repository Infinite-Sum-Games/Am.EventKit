package messagequeue

import (
	"encoding/json"
	"fmt"
	"strings"

	db "github.com/Infinite-Sum-Games/Am.EventKit/db/gen"
)

// Payload builder for Hackathon events.
type HackathonPayload struct {
	TeamName          string                `json:"team_name"`
	LeaderName        string                `json:"leader_name"`
	LeaderEmail       string                `json:"leader_email"`
	LeaderPhoneNumber string                `json:"leader_phone_number"`
	LeaderCollegeName string                `json:"leader_college_name"`
	ProblemStatement  string                `json:"problem_statement"`
	TeamMembers       []HackathonTeamMember `json:"team_members"`
}

type HackathonTeamMember struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	CollegeName string `json:"college_name"`
}

// CreateHackathonPayload builds the hackathon payload for metadata / MQ.
func CreateHackathonPayload(payload HackathonPayload) ([]byte, error) {
	return json.Marshal(payload)
}

func BuildHackathonTeamMembers(
	students []db.Student,
	leaderEmail string,
) ([]HackathonTeamMember, error) {

	members := make([]HackathonTeamMember, 0, len(students))

	for _, s := range students {
		// Skip leader
		if strings.EqualFold(s.Email, leaderEmail) {
			continue
		}

		if s.Name == "" || s.Email == "" {
			return nil, fmt.Errorf("invalid student data for hackathon payload: %s", s.Email)
		}

		members = append(members, HackathonTeamMember{
			Name:        s.Name,
			Email:       s.Email,
			PhoneNumber: s.PhoneNumber,
			CollegeName: s.CollegeName,
		})
	}

	return members, nil
}
