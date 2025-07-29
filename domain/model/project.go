package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	Company CompanyNested      `bson:"company" json:"company"`
	// CompanyProduct CompanyProductNested `bson:"companyProduct" json:"companyProduct"`
	Name        string     `bson:"name" json:"name"`
	Description string     `bson:"description" json:"description"`
	CreatedAt   time.Time  `bson:"createdAt" json:"createdAt"`
	CreatedBy   CustomerFK `bson:"createdBy" json:"createdBy"`
	UpdatedAt   time.Time  `bson:"updatedAt" json:"updatedAt"`
	DeletedAt   *time.Time `bson:"deletedAt" json:"-"`
}

type ProjectFK struct {
	ID   string `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
}
