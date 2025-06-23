package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Company struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	AccessKey    string             `bson:"accessKey" json:"accessKey"`
	Name         string             `bson:"name" json:"name"`
	Bio          string             `bson:"bio" json:"bio"`
	Type         string             `bson:"type" json:"type"`
	ProductTotal int64              `bson:"productTotal" json:"productTotal"`
	TicketTotal  int64              `bson:"ticketTotal" json:"ticketTotal"`
	Logo         MediaFK            `bson:"logo" json:"logo"`
	Settings     CompanySeting      `bson:"settings" json:"settings"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt    *time.Time         `bson:"deletedAt" json:"-"`
}

type CompanySeting struct {
	Code      string        `bson:"code" json:"code"`
	Email     string        `bson:"email" json:"email"`
	ColorMode ColorMode     `bson:"colorMode" json:"colorMode"`
	Domain    CompanyDomain `bson:"domain" json:"domain"`
	SMTP      SMTP          `bson:"smtp" json:"smtp"`
}

type SMTP struct {
	FromAddress string `bson:"fromAddress" json:"fromAddress"`
	FromName    string `bson:"fromName" json:"fromName"`
}

type ColorMode struct {
	Light Color `bson:"light" json:"light"`
	Dark  Color `bson:"dark" json:"dark"`
}

type Color struct {
	Primary   string `bson:"primary" json:"primary"`
	Secondary string `bson:"secondary" json:"secondary"`
}

type CompanyDomain struct {
	IsCustom  bool   `bson:"isCustom" json:"isCustom"`
	Subdomain string `bson:"subdomain" json:"subdomain"`
	FullUrl   string `bson:"fullUrl" json:"fullUrl"`
}

type CompanyNested struct {
	ID    string `bson:"id" json:"id"`
	Name  string `bson:"name" json:"name"`
	Image string `bson:"image" json:"image"`
	Type  string `bson:"type" json:"type"`
}

type CompanyProduct struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	Company        CompanyNested      `bson:"company" json:"company"`
	Name           string             `bson:"name" json:"name"`
	Code           string             `bson:"code" json:"code"`
	Logo           MediaFK            `bson:"logo" json:"logo"`
	TicketTotal    int64              `bson:"ticketTotal" json:"ticketTotal"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
	LastActivityAt *time.Time         `bson:"lastActivityAt" json:"lastActivityAt"`
	DeletedAt      *time.Time         `bson:"deletedAt" json:"-"`
}

type CompanyProductNested struct {
	ID    string `bson:"id" json:"id"`
	Name  string `bson:"name" json:"name"`
	Image string `bson:"image" json:"image"`
	Code  string `bson:"code" json:"code"`
}
