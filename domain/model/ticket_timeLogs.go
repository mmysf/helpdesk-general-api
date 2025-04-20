package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TicketTimeLogs struct {
	ID                     primitive.ObjectID   `bson:"_id" json:"id"`
	Company                CompanyNested        `bson:"company" json:"company"`
	Customer               CustomerFK           `bson:"customer" json:"customer"`
	Product                CompanyProductNested `bson:"product" json:"product"`
	Ticket                 TicketNested         `bson:"ticket" json:"ticket"`
	DurationInSeconds      int                  `bson:"durationInSeconds" json:"durationInSeconds"`
	PauseDurationInSeconds int                  `bson:"pauseDurationInSeconds" json:"pauseDurationInSeconds"`
	StartAt                *time.Time           `bson:"startAt" json:"startAt"`
	EndAt                  *time.Time           `bson:"endAt" json:"endAt"`
	PauseHistory           []PauseHistory       `bson:"pauseHistory" json:"pauseHistory"`
	IsManual               bool                 `bson:"isManual" json:"isManual"`
	ActivityTtype          string               `bson:"activityType" json:"activityType"`
	CreatedAt              time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt              *time.Time           `bson:"updatedAt" json:"updatedAt"`
	DeletedAt              *time.Time           `bson:"deletedAt" json:"-"`
}
