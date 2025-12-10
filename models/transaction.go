package models

type VerifyTransactionRequest struct {
	TxnID string `json:"txn_id" binding:"required"`
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
