package models

type VerifyTransactionRequest struct {
	TxnID string `json:"txn_id" binding:"required"`
}

const (
	StatusFailed    = "FAILED"
	StatusSuccess   = "SUCCESS"
	StatusPending   = "PENDING"
	PaymentNotFound = "NOT_FOUND"
)

type PayUVerifyResponse struct {
	Status             int                      `json:"status"`
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
		return StatusSuccess
	case "failure":
		return StatusFailed
	case "pending":
		return StatusPending
	default:
		return PaymentNotFound
	}
}
