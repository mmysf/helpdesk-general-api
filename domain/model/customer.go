package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	Company CompanyNested      `bson:"company" json:"company"`
	// CompanyProduct     CompanyProductNested `bson:"companyProduct" json:"companyProduct"`
	Name               string        `bson:"name" json:"name"`
	Email              string        `bson:"email" json:"email"`
	Password           string        `bson:"password" json:"-"`
	IsNeedBalance      bool          `bson:"isNeedBalance" json:"isNeedBalance"`
	Subscription       *Subscription `bson:"subscription" json:"subscription"`
	ProfilePicture     MediaFK       `bson:"profilePicture" json:"profilePicture"`
	JobTitle           string        `bson:"jobTitle" json:"jobTitle"`
	Bio                string        `bson:"bio" json:"bio"`
	Role               UserRole      `bson:"role" json:"role"`
	TickeTotal         int64         `bson:"ticketTotal" json:"ticketTotal"`
	Token              string        `bson:"token" json:"-"`
	PasswordResetToken string        `bson:"passwordResetToken" json:"-"`
	IsVerified         bool          `bson:"isVerified" json:"isVerified"`
	VerifiedAt         *time.Time    `bson:"verifiedAt" json:"-"`
	LastActivityAt     *time.Time    `bson:"lastActivityAt" json:"lastActivityAt"`
	CreatedAt          time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time     `bson:"updatedAt" json:"updatedAt"`
	DeletedAt          *time.Time    `bson:"deletedAt" json:"-"`
}

type CustomerFK struct {
	ID    string `bson:"id" json:"id,omitempty"`
	Name  string `bson:"name" json:"name,omitempty"`
	Email string `bson:"email" json:"email"`
}

type Subscription struct {
	Status      SubscriptionStatus `bson:"status" json:"status"`
	HourPackage *HourPackageFK     `bson:"hourPackage" json:"hourPackage"`
	Balance     *Balance           `bson:"balance" json:"balance"`
	StartAt     time.Time          `bson:"startAt" json:"startAt"`
	EndAt       time.Time          `bson:"endAt" json:"endAt"`
}

type SubscriptionStatus string

const (
	Active  SubscriptionStatus = "active"
	Expired SubscriptionStatus = "expired"
)

type Balance struct {
	Time   TimeBalance   `bson:"time" json:"time"`
	Ticket TicketBalance `bson:"ticket" json:"ticket"`
}

type TimeBalance struct {
	Total     int64           `bson:"total" json:"total"`
	Remaining RemainingNested `bson:"-" json:"remaining"`
	Used      int64           `bson:"used" json:"used"`
}

type RemainingNested struct {
	Total  int64 `bson:"-" json:"total"`
	Hour   int64 `bson:"-" json:"hour"`
	Minute int64 `bson:"-" json:"minute"`
	Second int64 `bson:"-" json:"second"`
}

type TicketBalance struct {
	Remaining int64 `bson:"remaining" json:"remaining"`
	Used      int64 `bson:"used" json:"used"`
}

type UserRole string

const (
	AdminRole    UserRole = "admin"
	AgentRole    UserRole = "agent"
	CustomerRole UserRole = "customer"
)
