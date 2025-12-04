package pkg

import (
	"time"

	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
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

func ToPgTextPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func ToPgNumericFromFloat(f float64) (pgtype.Numeric, error) {
	// Convert to string for safe parsing
	s := strconv.FormatFloat(f, 'f', -1, 64)

	parts := strings.Split(s, ".")
	var intStr string
	var exp int32

	if len(parts) == 1 {
		// whole number
		intStr = parts[0]
		exp = 0
	} else {
		// decimal number
		intStr = parts[0] + parts[1]
		exp = -int32(len(parts[1]))
	}

	bigInt := new(big.Int)
	_, ok := bigInt.SetString(intStr, 10)
	if !ok {
		return pgtype.Numeric{}, fmt.Errorf("invalid numeric")
	}

	return pgtype.Numeric{
		Int:   bigInt,
		Exp:   exp,
		Valid: true,
	}, nil
}

func GenerateTxnID(studentID uuid.UUID, eventID uuid.UUID) string {
	// Shorten UUIDs to 4 chars for readability
	shortStudent := studentID.String()[:4]
	shortEvent := eventID.String()[:4]

	// Add a timestamp
	timestamp := time.Now().UTC().Format("20060102T150405")

	// Small random suffix to avoid collisions
	randomBytes := make([]byte, 2)
	_, _ = rand.Read(randomBytes)
	randomHex := hex.EncodeToString(randomBytes)

	return fmt.Sprintf("TXN-%s-%s-%s-%s", timestamp, shortStudent, shortEvent, randomHex)
}
