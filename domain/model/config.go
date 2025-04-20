package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConfigPublic struct {
	ID                 primitive.ObjectID `bson:"_id" json:"id"`
	AppName            string             `bson:"appName" json:"appName"`
	CloseTicketLink    string             `bson:"closeTicketLink" json:"-"`
	AgentLink          string             `bson:"agentLink" json:"-"`
	Maintenance        ConfigMaintenance  `bson:"maintenance" json:"maintenance"`
	Email              ConfigEmail        `bson:"email" json:"-"`
	Registration       ConfigRegistration `bson:"registration" json:"-"`
	ResetPasswordLink  string             `bson:"resetPasswordLink" json:"-"`
	LoginLink          string             `bson:"loginLink" json:"-"`
	MainDomain         string             `bson:"mainDomain" json:"mainDomain"`
	CS                 string             `bson:"cs" json:"cs"`
	MinimumCredit      int64              `bson:"minimumCredit" json:"minimumCredit"`
	DollarInIdr        float64            `bson:"dollarInIdr" json:"dollarInIdr"`
	DefaultColor       ColorMode          `bson:"defaultColor" json:"defaultColor"`
	ManualPayment      ManualPayment      `bson:"manualPayment" json:"manualPayment"`
	BlacklistSubdomain []string           `bson:"blacklistSubdomain" json:"-"`
	CreatedAt          time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type Config struct {
	ID                 primitive.ObjectID `bson:"_id" json:"id"`
	AppName            string             `bson:"appName" json:"appName"`
	CloseTicketLink    string             `bson:"closeTicketLink" json:"closeTicketLink"`
	AgentLink          string             `bson:"agentLink" json:"agentLink"`
	Maintenance        ConfigMaintenance  `bson:"maintenance" json:"maintenance"`
	Email              ConfigEmail        `bson:"email" json:"email"`
	Registration       ConfigRegistration `bson:"registration" json:"registration"`
	ResetPasswordLink  string             `bson:"resetPasswordLink" json:"resetPasswordLink"`
	LoginLink          string             `bson:"loginLink" json:"loginLink"`
	MainDomain         string             `bson:"mainDomain" json:"mainDomain"`
	CS                 string             `bson:"cs" json:"cs"`
	MinimumCredit      int64              `bson:"minimumCredit" json:"minimumCredit"`
	DollarInIdr        float64            `bson:"dollarInIdr" json:"dollarInIdr"`
	DefaultColor       ColorMode          `bson:"defaultColor" json:"defaultColor"`
	ManualPayment      ManualPayment      `bson:"manualPayment" json:"manualPayment"`
	BlacklistSubdomain []string           `bson:"blacklistSubdomain" json:"blacklistSubdomain"`
	CreatedAt          time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time          `bson:"updatedAt" json:"updatedAt"`
}

func (t *Config) ToPublic() ConfigPublic {
	return ConfigPublic(*t)
}

type ConfigMaintenance struct {
	IsMaintenance bool   `bson:"isMaintenance" json:"isMaintenance"`
	Message       string `bson:"message" json:"message"`
}

type ConfigRegistration struct {
	VerifyEmail      bool   `bson:"verifyEmail" json:"verifyEmail"`
	VerificationLink string `bson:"verificationLink" json:"verificationLink"`
}

type ConfigEmail struct {
	Template    EmailTemplate `bson:"template" json:"template"`
	SenderName  string        `bson:"senderName" json:"senderName"`
	SenderEmail string        `bson:"senderEmail" json:"senderEmail"`
}

type EmailTemplate struct {
	Register           TemplateEmailConfig `bson:"register" json:"register"`
	DefaultUser        TemplateEmailConfig `bson:"defaultUser" json:"defaultUser"`
	ResetPassword      TemplateEmailConfig `bson:"resetPassword" json:"resetPassword"`
	ConfirmCloseTicket TemplateEmailConfig `bson:"confirmCloseTicket" json:"confirmCloseTicket"`
	TicketComment      TemplateEmailConfig `bson:"ticketComment" json:"ticketComment"`
	Stamped            TemplateEmailConfig `bson:"stamped" json:"stamped"`
	LinkGenerated      TemplateEmailConfig `bson:"linkGenerated" json:"linkGenerated"`
	CreateTicket       TemplateEmailConfig `bson:"createTicket" json:"createTicket"`
	CloseTicket        TemplateEmailConfig `bson:"closeTicket" json:"closeTicket"`
	ReopenTicket       TemplateEmailConfig `bson:"reopenTicket" json:"reopenTicket"`
	PackageActivated   TemplateEmailConfig `bson:"packageActivated" json:"packageActivated"`
	PackageExpired     TemplateEmailConfig `bson:"packageExpired" json:"packageExpired"`
}

type TemplateEmailConfig struct {
	Title string `bson:"title" json:"title"`
	Body  string `bson:"body" json:"body"`
}

type ManualPayment struct {
	IsActive         bool   `bson:"isActive" json:"isActive"`
	DurationInSecond int64  `bson:"durationInSecond" json:"durationInSecond"`
	AccountName      string `bson:"accountName" json:"accountName"`
	AccountNumber    string `bson:"accountNumber" json:"accountNumber"`
	BankName         string `bson:"bankName" json:"bankName"`
	SwiftCode        string `bson:"swiftCode" json:"swiftCode"`
}
