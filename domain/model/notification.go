package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationType string

const (
	NotificationTypePractitioner = "practitioner"
	NotificationTypeClient       = "client"
)

type Meta map[string]interface{}

type Notification struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Content     string             `bson:"content" json:"content"`
	Meta        Meta               `bson:"meta" json:"meta"`
	Link        string             `bson:"link" json:"link"`
	IsRead      bool               `bson:"isRead" json:"isRead"`
	IsImportant bool               `bson:"isImportant" json:"isImportant"`
	UserId      string             `bson:"userId" json:"userId"`
	Type        NotificationType   `bson:"type" json:"type"`
	Category    string             `bson:"category" json:"category"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}