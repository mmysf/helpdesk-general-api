package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	HourPackage     *HourPackageFK     `bson:"hourPackage" json:"hourPackage"`
	ServerPackage   *ServerPackageFK   `bson:"serverPackage" json:"serverPackage"`
	Customer        CustomerFK         `bson:"customer" json:"customer"`
	Invoice         InvoiceNested      `bson:"invoice" json:"invoice"`
	OrderNumber     string             `bson:"orderNumber" json:"orderNumber"`
	Status          OrderStatus        `bson:"status" json:"status"`
	Type            OrderType          `bson:"type" json:"type"`
	Amount          int64              `bson:"amount" json:"amount"`
	Tax             float64            `bson:"tax" json:"tax"`
	AdminFee        float64            `bson:"adminFee" json:"adminFee"`
	Discount        float64            `bson:"discount" json:"discount"`
	SubTotal        float64            `bson:"subTotal" json:"subTotal"`
	GrandTotal      float64            `bson:"grandTotal" json:"grandTotal"`
	GrandTotalinIdr float64            `bson:"-" json:"grandTotalInIdr"`
	Note            string             `bson:"note" json:"note"`
	Payment         Payment            `bson:"payment" json:"payment"`
	PaidAt          *time.Time         `bson:"paidAt" json:"paidAt"`
	ExpiredAt       time.Time          `bson:"expiredAt" json:"expiredAt"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt       *time.Time         `bson:"deletedAt" json:"-"`
}

func (o *Order) Format(c *Config) *Order {
	o.GrandTotalinIdr = o.GrandTotal * c.DollarInIdr
	return o
}

type OrderFK struct {
	ID          string    `bson:"id" json:"id"`
	OrderNumber string    `bson:"orderNumber" json:"orderNumber"`
	Type        OrderType `bson:"type" json:"type"`
}

type InvoiceNested struct {
	InvoiceURL         string `bson:"invoiceURL" json:"invoiceURL"`
	InvoiceExternalId  string `bson:"invoiceExternalId" json:"invoiceExternalId"`
	InvoiceXenditId    string `bson:"invoiceXenditId" json:"invoiceXenditId"`
	MerchantName       string `bson:"merchantName" json:"merchantName"`
	PaymentMethod      string `bson:"PaymentMethod" json:"PaymentMethod"`
	BankCode           string `bson:"bankCode" json:"bankCode"`
	PaymentChannel     string `bson:"paymentChannel" json:"paymentChannel"`
	PaymentDestination string `bson:"paymentDestination" json:"paymentDestination"`
	SwiftCode          string `bson:"swiftCode" json:"swiftCode"`
}

type OrderStatus string

const (
	STATUS_PENDING          OrderStatus = "pending"
	STATUS_WAITING_APPROVAL OrderStatus = "waiting_approval"
	STATUS_REJECT           OrderStatus = "reject"
	STATUS_PAID             OrderStatus = "paid"
	STATUS_EXPIRED          OrderStatus = "expired"
)

type OrderType string

const (
	HOUR_TYPE   OrderType = "HOUR"
	SERVER_TYPE OrderType = "SERVER"
)

type Payment struct {
	Status     string      `bson:"status" json:"status"`
	ManualPaid *ManualPaid `bson:"manualPaid" json:"manualPaid"`
	PaidAt     *time.Time  `bson:"paidAt" json:"paidAt"`
	Webhook    Webhook     `bson:"webhook" json:"-"`
	Snap       interface{} `bson:"snap" json:"-"`
}

type ManualPaid struct {
	AccountName   string    `bson:"accountName" json:"accountName"`
	AccountNumber string    `bson:"accountNumber" json:"accountNumber"`
	BankName      string    `bson:"bankName" json:"bankName"`
	Note          string    `bson:"note" json:"note"`
	Attachment    MediaFK   `bson:"attachment" json:"attachment"`
	Approval      *Approval `bson:"approval" json:"approval"`
}

type Approval struct {
	User UserNested `bson:"user" json:"user"`
	Note string     `bson:"note" json:"note"`
	At   time.Time  `bson:"at" json:"at"`
}

type Webhook struct {
	Detail  interface{}   `bson:"detail" json:"detail"`
	History []interface{} `bson:"history" json:"history"`
}
