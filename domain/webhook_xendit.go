package domain

import (
	"time"
)

type QrisWebhookRequest struct {
	Event      string                 `json:"event"`
	Created    time.Time              `json:"created"`
	BusinessID string                 `json:"business_id"`
	Data       DataQrisWebhookRequest `json:"data"`
}

type DataQrisWebhookRequest struct {
	ID            string        `json:"id"`
	BusinessID    string        `json:"business_id"`
	Currency      string        `json:"currency"`
	Amount        int64         `json:"amount"`
	Status        string        `json:"status"`
	Created       time.Time     `json:"created"`
	QrID          string        `json:"qr_id"`
	QrString      string        `json:"qr_string"`
	ReferenceID   string        `json:"reference_id"`
	Type          string        `json:"type"`
	ChannelCode   string        `json:"channel_code"`
	ExpiresAt     time.Time     `json:"expires_at"`
	Metadata      Metadata      `json:"metadata"`
	PaymentDetail PaymentDetail `json:"payment_detail"`
}

type Metadata struct {
	BranchCode string `json:"branch_code"`
}

type PaymentDetail struct {
	ReceiptID string `json:"receipt_id"`
	Source    string `json:"source"`
}

type SnapWebhookRequest struct {
	ID                 string    `json:"id"`
	ExternalID         string    `json:"external_id"`
	UserID             string    `json:"user_id"`
	IsHigh             bool      `json:"is_high"`
	PaymentMethod      string    `json:"payment_method"`
	Status             string    `json:"status"`
	MerchantName       string    `json:"merchant_name"`
	Amount             int64     `json:"amount"`
	PaidAmount         int64     `json:"paid_amount"`
	BankCode           string    `json:"bank_code"`
	PaidAt             time.Time `json:"paid_at"`
	PayerEmail         string    `json:"payer_email"`
	Description        string    `json:"description"`
	Created            time.Time `json:"created"`
	Updated            time.Time `json:"updated"`
	Currency           string    `json:"currency"`
	PaymentChannel     string    `json:"payment_channel"`
	PaymentDestination string    `json:"payment_destination"`
}

type XenditGenereteQRResponseSuccess struct {
	ReferenceID string      `json:"reference_id"`
	Type        string      `json:"type"`
	Currency    string      `json:"currency"`
	ChannelCode string      `json:"channel_code"`
	Amount      int64       `json:"amount"`
	ExpiresAt   time.Time   `json:"expires_at"`
	Basket      []Basket    `json:"basket"`
	Metadata    interface{} `json:"metadata"`
	BusinessID  string      `json:"business_id"`
	ID          string      `json:"id"`
	Created     time.Time   `json:"created"`
	Updated     time.Time   `json:"updated"`
	QrString    string      `json:"qr_string"`
	Status      string      `json:"status"`
}

type Basket struct {
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Price    int64  `json:"price"`
	Quantity int64  `json:"quantity"`
	Type     string `json:"type"`
}

type XenditResponseError struct {
	ErrorCode string  `json:"error_code"`
	Message   string  `json:"message"`
	Errors    []Error `json:"errors"`
}

type Error struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

type TopupPayloadPaymentLink struct {
	ID         string `json:"id"`
	QrString   string `json:"qrString"`
	QrId       string `json:"qrId"`
	ExpiresAt  string `json:"expiresAt"`
	GrandTotal int64  `json:"grandTotal"`
}

type XenditGenereteSnapLinkResponseSuccess struct {
	ID                        string                  `json:"id"`
	ExternalID                string                  `json:"external_id"`
	UserID                    string                  `json:"user_id"`
	Status                    string                  `json:"status"`
	MerchantName              string                  `json:"merchant_name"`
	MerchantProfilePictureURL string                  `json:"merchant_profile_picture_url"`
	Amount                    int64                   `json:"amount"`
	PayerEmail                string                  `json:"payer_email"`
	Description               string                  `json:"description"`
	ExpiryDate                time.Time               `json:"expiry_date"`
	InvoiceURL                string                  `json:"invoice_url"`
	AvailableBanks            []AvailableBank         `json:"available_banks"`
	AvailableRetailOutlets    []AvailableRetailOutlet `json:"available_retail_outlets"`
	AvailableEwallets         []AvailableEwallet      `json:"available_ewallets"`
	AvailableQrCodes          []AvailableQrCode       `json:"available_qr_codes"`
	AvailableDirectDebits     []AvailableDirectDebit  `json:"available_direct_debits"`
	AvailablePaylaters        []AvailablePaylater     `json:"available_paylaters"`
	ShouldExcludeCreditCard   bool                    `json:"should_exclude_credit_card"`
	ShouldSendEmail           bool                    `json:"should_send_email"`
	Created                   time.Time               `json:"created"`
	Updated                   time.Time               `json:"updated"`
	Currency                  string                  `json:"currency"`
	ReminderDate              time.Time               `json:"reminder_date"`
	Metadata                  interface{}             `json:"metadata"`
}

type AvailableBank struct {
	BankCode          string `json:"bank_code"`
	CollectionType    string `json:"collection_type"`
	TransferAmount    int64  `json:"transfer_amount"`
	BankBranch        string `json:"bank_branch"`
	AccountHolderName string `json:"account_holder_name"`
	IdentityAmount    int64  `json:"identity_amount"`
	BankAccountNumber string `json:"bank_account_number"`
}

type AvailableDirectDebit struct {
	DirectDebitType string `json:"direct_debit_type"`
}

type AvailableEwallet struct {
	EwalletType string `json:"ewallet_type"`
}

type AvailablePaylater struct {
	PaylaterType string `json:"paylater_type"`
}

type AvailableQrCode struct {
	QrCodeType string `json:"qr_code_type"`
}

type AvailableRetailOutlet struct {
	RetailOutletName string `json:"retail_outlet_name"`
	PaymentCode      string `json:"payment_code"`
	TransferAmount   int64  `json:"transfer_amount"`
	MerchantName     string `json:"merchant_name"`
}
