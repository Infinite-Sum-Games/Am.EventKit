package models

type VerifyTransactionRequest struct {
	TxnID string `json:"txn_id" binding:"required"`
}
