package pkg

import "github.com/jackc/pgx/v5/pgtype"

func ToPgText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{String: "", Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}
