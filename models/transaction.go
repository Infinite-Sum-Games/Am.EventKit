package models

type VerifyTransactionRequest struct {
	TxnID string `json:"txn_id" binding:"required"`
}

const (
	StatusFailed   = "FAILED"
	StatusSuccess  = "SUCCESS"
	StatusPending  = "PENDING"
	StatusNotFound = "NOT_FOUND"
)
