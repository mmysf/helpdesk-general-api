package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Agent struct {
	ID                 primitive.ObjectID `bson:"_id" json:"id"`
	Company            CompanyNested      `bson:"company" json:"company"`
	Name               string             `bson:"name" json:"name"`
	Email              string             `bson:"email" json:"email"`
	Password           string             `bson:"password" json:"-"`
	JobTitle           string             `bson:"jobTitle" json:"jobTitle"`
	ProfilePicture     MediaFK            `bson:"profilePicture" json:"profilePicture"`
	Bio                string             `bson:"bio" json:"bio"`
	Role               UserRole           `bson:"role" json:"role"`
	Category           TicketCategoryFK   `bson:"category" json:"category"`
	LastActivityAt     *time.Time         `bson:"lastActivityAt" json:"lastActivityAt"`
	PasswordResetToken string             `bson:"passwordResetToken" json:"-"`
	CreatedAt          time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt          *time.Time         `bson:"deletedAt" json:"-"`
}

type AgentNested struct {
	ID    string `bson:"id" json:"id,omitempty"`
	Name  string `bson:"name" json:"name,omitempty"`
	Email string `bson:"email" json:"email,omitempty"`
}
