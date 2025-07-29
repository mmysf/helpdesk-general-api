package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ticket struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	Company CompanyNested      `bson:"company" json:"company"`
	// Product      CompanyProductNested `bson:"product" json:"product"`
	Project      *ProjectFK        `bson:"project" json:"project"`
	Category     *TicketCategoryFK `bson:"category" json:"category"`
	Customer     CustomerFK        `bson:"customer" json:"customer"`
	Agent        []AgentNested     `bson:"agents" json:"agents"`
	AssignedToMe *bool             `bson:"-" json:"assignedToMe,omitempty"`
	Subject      string            `bson:"subject" json:"subject"`
	Content      string            `bson:"content" json:"content"`
	Code         string            `bson:"code" json:"code"`
	Name         string            `bson:"name" json:"name"`
	Attachments  []AttachmentFK    `bson:"attachments" json:"attachments"`
	LogTime      LogTime           `bson:"logTime" json:"logTime"`
	Priority     TicketPriority    `bson:"priority" json:"priority"`
	Status       TicketStatus      `bson:"status" json:"status"`
	ReminderSent bool              `bson:"reminderSent" json:"reminderSent"`
	Token        string            `bson:"token" json:"-"`
	DetailTime   DetailTime        `bson:"detailTime" json:"detailTime"`
	Parent       *TicketNested     `bson:"parent" json:"parent"`
	CompletedBy  *AgentNested      `bson:"completedBy" json:"completedBy"`
	ClosedAt     *time.Time        `bson:"closedAt" json:"closedAt"`
	CreatedAt    time.Time         `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time         `bson:"updatedAt" json:"updatedAt"`
	DeletedAt    *time.Time        `bson:"deletedAt" json:"-"`
}

type TicketStatus string

const (
	Open       TicketStatus = "open"
	Processing TicketStatus = "processing"
	Closed     TicketStatus = "closed"
	InProgress TicketStatus = "in_progress"
	Resolve    TicketStatus = "resolve"
	Cancel     TicketStatus = "cancel"
)

type TicketNested struct {
	ID       string         `bson:"id" json:"id"`
	Subject  string         `bson:"subject" json:"subject"`
	Content  string         `bson:"content" json:"content"`
	Priority TicketPriority `bson:"priority" json:"priority"`
}

type TicketComment struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	Company CompanyNested      `bson:"company" json:"company"`
	// Product     CompanyProductNested `bson:"product" json:"product"`
	Ticket      TicketNested   `bson:"ticket" json:"ticket"`
	Agent       AgentNested    `bson:"agent" json:"agent"`
	Customer    CustomerFK     `bson:"customer" json:"customer"`
	Sender      SenderType     `bson:"sender" json:"sender"` // agent | customer
	Content     string         `bson:"content" json:"content"`
	Attachments []AttachmentFK `bson:"attachments" json:"attachments"`
	CreatedAt   time.Time      `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time      `bson:"updatedAt" json:"updatedAt"`
	DeletedAt   *time.Time     `bson:"deletedAt" json:"-"`
}

type SenderType string

const (
	AgentSender    SenderType = "agent"
	CustomerSender SenderType = "customer"
)

type LogTime struct {
	StartAt                      *time.Time     `bson:"startAt" json:"startAt"`
	EndAt                        *time.Time     `bson:"endAt" json:"endAt"`
	DurationInSeconds            int            `bson:"durationInSeconds" json:"durationInSeconds"`
	PauseDurationInSeconds       int            `bson:"pauseDurationInSeconds" json:"pauseDurationInSeconds"`
	Status                       LogTimeStatus  `bson:"status" json:"status"`
	TotalDurationInSeconds       int            `bson:"totalDurationInSeconds" json:"totalDurationInSeconds"`
	TotalPausedDurationInSeconds int            `bson:"totalPauseDurationInSeconds" json:"totalPauseDurationInSeconds"`
	PauseHistory                 []PauseHistory `bson:"pauseHistory" json:"pauseHistory"`
}

type PauseHistory struct {
	PausedAt       time.Time  `bson:"pausedAt" json:"pausedAt"`
	ResumedAt      *time.Time `bson:"resumedAt" json:"resumedAt"`
	DurationActive int        `bson:"-" json:"durationActive"`
}

type LogTimeStatus string

const (
	NotStarted LogTimeStatus = "not_started"
	Paused     LogTimeStatus = "paused"
	Running    LogTimeStatus = "running"
	Done       LogTimeStatus = "done"
)

type TicketPriority string

const (
	PriorityLow      TicketPriority = "(P4) Low"
	PriorityMedium   TicketPriority = "(P3) Medium"
	PriorityHigh     TicketPriority = "(P2) High"
	PriorityCritical TicketPriority = "(P1) Critical"
)

type DetailTime struct {
	Year    int    `bson:"year" json:"year"`
	Month   int    `bson:"month" json:"month"`
	Day     int    `bson:"day" json:"day"`
	DayName string `bson:"dayName" json:"dayName"`
}

var Weekdays = []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}

func (t *Ticket) Format(userID string) *Ticket {
	// check is ticket assigned to me
	if t.Company.Type == "B2C" {
		condition := false
		for _, assignedAgent := range t.Agent {
			if assignedAgent.ID == userID {
				condition = true
				break
			}
		}
		t.AssignedToMe = &condition
	}

	return t
}
