package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Company   CompanyNested      `bson:"company" json:"company"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	IsRead    bool               `bson:"isRead" json:"isRead"`
	UserRole  UserRole           `bson:"userRole" json:"userRole"`
	User      UserNested         `bson:"user" json:"user"`
	Type      NotificationType   `bson:"type" json:"type"`
	Ticket    TicketNested       `bson:"ticket" json:"ticket"`
	Category  TicketCategoryFK   `bson:"category" json:"category"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type NotificationType string

const (
	TicketCreated NotificationType = "ticketCreated"
	TicketUpdated NotificationType = "ticketUpdated"
	TicketClosed  NotificationType = "ticketClosed"
)
