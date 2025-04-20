package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerSubscription struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	Customer      CustomerFK         `bson:"customer" json:"customer"`
	HourPackage   *HourPackageFK     `bson:"hourPackage" json:"hourPackage"`
	ServerPackage *ServerPackageFK   `bson:"serverPackage" json:"serverPackage"`
	Order         OrderFK            `bson:"order" json:"order"`
	Status        SubscriptionStatus `bson:"status" json:"status"`
	ExpiredAt     time.Time          `bson:"expiredAt" json:"expiredAt"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt     *time.Time         `bson:"deletedAt" json:"-"`
}
