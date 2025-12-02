package pkg

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToPgText(value string) pgtype.Text {
	if value == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: value, Valid: true}
}

func ToPgInt4(value int) pgtype.Int4 {
	if value == 0 {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(value), Valid: true}
}

func ParseTime(s string) (time.Time, error) {
	// Expect format HH:MM:SS (e.g. "10:00:00") and attach today's date (UTC).
	t, err := time.Parse("15:04:05", s)
	if err != nil {
		return time.Time{}, err
	}
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC), nil
}
