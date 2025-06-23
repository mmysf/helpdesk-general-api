package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerBalanceHistory struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Customer  CustomerFK         `bson:"customer" json:"customer"`
	In        int64              `bson:"in" json:"in"`
	Out       int64              `bson:"out" json:"out"`
	Reference Reference          `bson:"reference" json:"reference"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt *time.Time         `bson:"deletedAt" json:"-"`
}

type Reference struct {
	UniqueID string        `bson:"unique_id" json:"unique_id"`
	Type     ReferenceType `bson:"type" json:"type"`
}

type ReferenceType string

const (
	OrderReference  ReferenceType = "order"
	TicketReference ReferenceType = "ticket"
)
