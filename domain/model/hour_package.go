package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HourPackage struct {
	ID          primitive.ObjectID    `bson:"_id" json:"id"`
	Name        string                `bson:"name" json:"name"`
	Description string                `bson:"description" json:"description"`
	Benefit     []string              `bson:"benefit" json:"benefit"`
	Price       float64               `bson:"price" json:"price"`
	Duration    HourPackageDuration   `bson:"duration" json:"duration"`
	Additional  HourPackageAdditional `bson:"additional" json:"additional"`
	Status      HourPackageStatus     `bson:"status" json:"status"`
	CreatedAt   time.Time             `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time             `bson:"updatedAt" json:"updatedAt"`
	DeletedAt   *time.Time            `bson:"deletedAt" json:"-"`
}

type HourPackageFK struct {
	ID      string   `bson:"id" json:"id"`
	Name    string   `bson:"name" json:"name"`
	Hours   int64    `bson:"hours" json:"hours"`
	Price   float64  `bson:"price" json:"price"`
	Benefit []string `bson:"benefit" json:"benefit"`
}

type HourPackageDuration struct {
	Hours          int64 `bson:"hours" json:"hours"`
	TotalinSeconds int64 `bson:"totalInSeconds" json:"totalInSeconds"`
}

type HourPackageAdditional struct {
	TicketLimit *int64 `bson:"ticketLimit" json:"ticketLimit"`
}

type HourPackageStatus string

const (
	HourPackageActive   HourPackageStatus = "active"
	HourPackageInactive HourPackageStatus = "inactive"
)
