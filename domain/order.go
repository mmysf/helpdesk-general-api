package domain

type OrderRequest struct {
	PackageID string `json:"packageId"`
	Amount    int64  `json:"amount"`
}

type ConfrimOrderRequest struct {
	OrderID       string `json:"orderId"`
	AttachmentId  string `json:"attachmentId"`
	AccountName   string `json:"accountName"`
	AccountNumber string `json:"accountNumber"`
	BankName      string `json:"BankName"`
	Note          string `json:"note"`
}

type UpdateManualPaymentRequest struct {
	Status        string `json:"status"`
	Note          string `json:"note"`
	AttachmentId  string `json:"attachmentId"`
	AccountName   string `json:"accountName"`
	AccountNumber string `json:"accountNumber"`
	BankName      string `json:"BankName"`
}
