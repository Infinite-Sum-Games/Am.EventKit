package models

import (
	"regexp"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

var Txn_regex = regexp.MustCompile(`^TXN-ANK26`)

type TeamBookingRequest struct {
	TeamName    string       `json:"team_name" binding:"required"`
	TeamMembers []TeamMember `json:"team_members" binding:"required"`
	ProblemStmt *string      `json:"ps,omitempty"`
}

type TeamMember struct {
	StudentEmail string `json:"student_email" binding:"required"`
	StudentRole  string `json:"student_role" binding:"required"`
}

// Validation for TeamMember
func (t TeamMember) Validate() error {
	return v.ValidateStruct(&t,
		v.Field(&t.StudentEmail, v.Required, v.Length(3, 200), is.Email),
		v.Field(&t.StudentRole, v.Required, v.Length(1, 100)),
	)
}

func (e TeamBookingRequest) Validate() error {
	return v.ValidateStruct(&e,
		v.Field(&e.TeamName, v.Required, v.Length(2, 100)),
		v.Field(&e.TeamMembers,
			v.Length(0, 100),
			v.Each(v.By(func(value any) error {
				if tm, ok := value.(TeamMember); ok {
					return tm.Validate()
				}
				return nil
			})),
		),
		// asuming the only metadata we get is related to problem statement type
		v.Field(
			&e.ProblemStmt,
			v.When(
				e.ProblemStmt != nil,
				v.In("agentic_ai", "generative_ai", "aiot"),
			),
		),
	)
}

type VerifyTransactionRequest struct {
	TxnID string `json:"txn_id" binding:"required"`
}

func (s VerifyTransactionRequest) Validate() error {
	return v.ValidateStruct(&s,
		v.Field(&s.TxnID,
			v.Required,
			v.Match(Txn_regex).
				Error("txn_id is not valid")))
}

const (
	PaymentFailed   = "FAILED"
	PaymentSuccess  = "SUCCESS"
	PaymentPending  = "PENDING"
	PaymentNotFound = "NOT_FOUND"
)

type PayUVerifyResponse struct {
	Status             int                      `json:"status"`
	Msg                string                   `json:"msg"`
	TransactionDetails map[string]PayUTxnDetail `json:"transaction_details"`
}

type PayUTxnDetail struct {
	Status string `json:"status"`
}

func MapPayUStatus(res PayUVerifyResponse, txnID string) string {
	// res.Status meaning:
	// - 0: Failed to fetch data
	// - 1: Successfully fetched
	if res.Status == 0 {
		return PaymentNotFound
	}

	detail, ok := res.TransactionDetails[txnID]
	if !ok {
		return PaymentNotFound
	}

	switch detail.Status {
	case "success":
		return PaymentSuccess
	case "failure":
		return PaymentFailed
	case "pending":
		return PaymentPending
	default:
		return PaymentNotFound
	}
}
