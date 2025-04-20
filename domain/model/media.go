package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Media struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Provider     string             `bson:"provider" json:"provider"`
	ProviderKey  string             `bson:"providerKey" json:"providerKey"`
	Type         string             `bson:"type" json:"type"`
	Category     MediaCategory      `bson:"category" json:"category"`
	Size         int64              `bson:"size" json:"size"`
	URL          string             `bson:"url" json:"url"`
	ExpiredUrlAt *time.Time         `bson:"expiredUrlAt" json:"expiredUrlAt"`
	IsUsed       bool               `bson:"isUsed" json:"isUsed"`
	IsPrivate    bool               `bson:"isPrivate" json:"isPrivate"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt    *time.Time         `bson:"deletedAt" json:"-"`
}

type MediaFK struct {
	ID   string `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
	Size int64  `bson:"size" json:"size"`
	URL  string `bson:"url" json:"url"`
	// ExpiredUrlAt *time.Time         `bson:"expiredUrlAt" json:"expiredUrlAt"`
	Type        string        `bson:"type" json:"type"`
	Category    MediaCategory `bson:"category" json:"category"`
	IsPrivate   bool          `bson:"isPrivate" json:"isPrivate"`
	ProviderKey string        `bson:"providerKey" json:"providerKey"`
}

type MediaCategory string

const (
	CompanyLogo MediaCategory = "COMPANY_LOGO"
	BrandLogo   MediaCategory = "BRAND_LOGO"
	Avatar      MediaCategory = "AVATAR"
	Other       MediaCategory = "OTHER"
)
