package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TicketCategory struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Company   CompanyNested      `bson:"company" json:"company"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt *time.Time         `bson:"deletedAt" json:"-"`
}

type TicketCategoryFK struct {
	ID   string `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
}
