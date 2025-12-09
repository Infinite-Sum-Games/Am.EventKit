package pkg

import (
	"encoding/json"
	"fmt"
	"strings"

	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
)

type WocStudentMeta struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type WocPayload struct {
	Queue    string           `json:"queue"`
	Students []WocStudentMeta `json:"students"`
}

func splitName(full string) (first string, last string) {
	parts := strings.Fields(full)
	if len(parts) == 0 {
		return "", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.Join(parts[1:], " ")
}

func BuildSpecialTags(tags []string, students []db.Student) ([]byte, error) {
	// Basic validation
	if len(tags) == 0 {
		return nil, nil
	}
	if len(students) == 0 {
		return nil, fmt.Errorf("no students provided for metadata generation")
	}

	// Normalize tags
	normalized := make([]string, 0, len(tags))
	for _, t := range tags {
		tag := strings.TrimSpace(strings.ToLower(t))
		if tag == "" {
			continue
		}
		normalized = append(normalized, tag)
	}

	if len(normalized) == 0 {
		return nil, nil
	}

	// Handle tags
	for _, t := range normalized {
		switch t {
		case "!woc":
			queueName := strings.TrimPrefix(t, "!") // TODO: Assuming that tagName is queueName
			meta, err := buildWocMetadata(students, queueName)
			if err != nil {
				return nil, fmt.Errorf("[BOOKING-ERROR]: woc metadata build failed: %w", err)
			}
			return meta, nil

		default:
			Log.Warn("[BOOKING-WARN] Unknown special tag encountered: ")
		}
	}

	return nil, nil
}

func buildWocMetadata(students []db.Student, queueName string) ([]byte, error) {
	if len(students) == 0 {
		return nil, fmt.Errorf("[BOOKING-ERROR]: cannot build woc metadata: empty student list")
	}

	res := make([]WocStudentMeta, 0, len(students))

	for _, s := range students {
		if s.Email == "" {
			return nil, fmt.Errorf("[BOOKING-ERROR]: student email missing for woc metadata")
		}
		if s.Name == "" {
			return nil, fmt.Errorf("[BOOKING-ERROR]: student name missing for woc metadata: email=%s", s.Email)
		}

		fn, ln := splitName(s.Name)
		res = append(res, WocStudentMeta{
			FirstName: fn,
			LastName:  ln,
			Email:     s.Email,
		})
	}

	return json.Marshal(map[string]any{
		"queue":    queueName,
		"students": res,
	})
}
