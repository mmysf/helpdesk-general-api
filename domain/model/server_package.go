package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServerPackage struct {
	ID           primitive.ObjectID  `bson:"_id" json:"id"`
	Name         string              `bson:"name" json:"name"`
	Description  string              `bson:"description" json:"description"`
	Benefit      []string            `bson:"benefit" json:"benefit"`
	Price        float64             `bson:"price" json:"price"`
	Customizable bool                `bson:"customizable" json:"customizable"`
	Validity     int64               `bson:"validity" json:"validity"`
	Status       ServerPackageStatus `bson:"status" json:"status"`
	CreatedAt    time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time           `bson:"updatedAt" json:"updatedAt"`
	DeletedAt    *time.Time          `bson:"deletedAt" json:"-"`
}

type ServerPackageStatus string

const (
	ServerPackageActive   ServerPackageStatus = "active"
	ServerPackageInactive ServerPackageStatus = "inactive"
)

type ServerPackageFK struct {
	ID           string   `bson:"id" json:"id"`
	Name         string   `bson:"name" json:"name"`
	Price        float64  `bson:"price" json:"price"`
	Customizable bool     `bson:"customizable" json:"customizable"`
	Validity     int64    `bson:"validity" json:"validity"`
	Benefit      []string `bson:"benefit" json:"benefit"`
}
