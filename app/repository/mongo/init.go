package mongorepo

import (
	"app/domain/model"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDBRepo struct {
	Conn                             *mongo.Database
	ConfigCollection                 string
	CustomerCollection               string
	AgentCollection                  string
	CompanyCollection                string
	CompanyProductCollection         string
	TicketCollection                 string
	TicketCommentCollection          string
	SuperuserCollection              string
	AttachmentCollection             string
	TicketTimelogsCollection         string
	SuperadminCollection             string
	HourPackageCollection            string
	OrderCollection                  string
	CustomerSubscriptionCollection   string
	CustomerBalanceHistoryCollection string
	ProjectCollection                string
	MediaCollection                  string
	TicketCategoryCollection         string
	ServerPackageCollection          string
}

func NewMongodbRepo(Conn *mongo.Database) MongoDBRepo {
	return &mongoDBRepo{
		Conn:                             Conn,
		ConfigCollection:                 "config",
		CustomerCollection:               "customers",
		AgentCollection:                  "agents",
		CompanyCollection:                "companies",
		CompanyProductCollection:         "company_products",
		TicketCollection:                 "tickets",
		TicketCommentCollection:          "ticket_comments",
		SuperuserCollection:              "superusers",
		AttachmentCollection:             "attachments",
		TicketTimelogsCollection:         "ticket_timelogs",
		SuperadminCollection:             "superadmins",
		HourPackageCollection:            "hour_packages",
		OrderCollection:                  "orders",
		CustomerSubscriptionCollection:   "customer_subscriptions",
		CustomerBalanceHistoryCollection: "customer_balance_histories",
		ProjectCollection:                "projects",
		MediaCollection:                  "medias",
		TicketCategoryCollection:         "ticket_categories",
		ServerPackageCollection:          "server_packages",
	}
}

type MongoDBRepo interface {
	// Customer
	FetchCustomerList(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error)
	CountCustomer(ctx context.Context, options map[string]interface{}) (total int64)
	FetchOneCustomer(ctx context.Context, options map[string]interface{}) (*model.Customer, error)
	CreateCustomer(ctx context.Context, usermodel *model.Customer) (err error)
	CreateManyCustomer(ctx context.Context, rows []*model.Customer) (err error)
	UpdateOneCustomer(ctx context.Context, query, payload map[string]interface{}) (err error)
	UpdatePartialCustomer(ctx context.Context, options, field map[string]interface{}) (err error)
	UpdateManyPartialCustomer(ctx context.Context, ids []primitive.ObjectID, field map[string]interface{}) (err error)

	// Agent
	FetchOneAgent(ctx context.Context, options map[string]interface{}) (*model.Agent, error)
	FetchAgentList(ctx context.Context, options map[string]interface{}) (*mongo.Cursor, error)
	CountAgent(ctx context.Context, options map[string]interface{}) (total int64)
	CreateAgent(ctx context.Context, row *model.Agent) (err error)
	UpdateAgent(ctx context.Context, row *model.Agent) (err error)
	UpdateOneAgent(ctx context.Context, query, payload map[string]interface{}) (err error)
	DeleteAgent(ctx context.Context, row *model.Agent) (err error)
	UpdatePartialAgent(ctx context.Context, options, field map[string]interface{}) (err error)

	// Config
	FetchOneConfig(ctx context.Context, options map[string]interface{}) (row *model.Config, err error)

	// Company
	FetchOneCompany(ctx context.Context, options map[string]interface{}) (*model.Company, error)
	FetchCompanyList(ctx context.Context, options map[string]interface{}) (*mongo.Cursor, error)
	CountCompany(ctx context.Context, options map[string]interface{}) int64
	CreateCompany(ctx context.Context, company *model.Company) (err error)
	UpdateCompany(ctx context.Context, company *model.Company) (err error)
	UpdatePartialCompany(ctx context.Context, options map[string]interface{}, field map[string]interface{}) (err error)
	IncrementOneCompany(ctx context.Context, id string, payload map[string]int64) (err error)
	DeleteCompany(ctx context.Context, company *model.Company) (err error)

	// CompanyProduct
	CreateCompanyProduct(ctx context.Context, company *model.CompanyProduct) (err error)
	UpdateCompanyProduct(ctx context.Context, company *model.CompanyProduct) (err error)
	UpdatePartialCompanyProduct(ctx context.Context, options map[string]interface{}, field map[string]interface{}) (err error)
	DeleteCompanyProduct(ctx context.Context, company *model.CompanyProduct) (err error)
	FetchOneCompanyProduct(ctx context.Context, options map[string]interface{}) (*model.CompanyProduct, error)
	FetchCompanyProductList(ctx context.Context, options map[string]interface{}) (*mongo.Cursor, error)
	CountCompanyProduct(ctx context.Context, options map[string]interface{}) int64
	IncrementOneCompanyProduct(ctx context.Context, id string, payload map[string]int64) (err error)

	// Ticket
	FetchTicketList(ctx context.Context, options map[string]interface{}) (*mongo.Cursor, error)
	CreateTicket(ctx context.Context, ticket *model.Ticket) (err error)
	FetchOneTicket(ctx context.Context, options map[string]interface{}) (*model.Ticket, error)
	UpdateTicket(ctx context.Context, ticket *model.Ticket) (err error)
	CountTicket(ctx context.Context, options map[string]interface{}) int64
	UpdateTicketPartial(ctx context.Context, ids primitive.ObjectID, field map[string]interface{}) error

	// Ticket Comment
	CreateTicketComment(ctx context.Context, ticket *model.TicketComment) (err error)
	FetchTicketCommentList(ctx context.Context, options map[string]interface{}) (*mongo.Cursor, error)
	FetchOneTicketComment(ctx context.Context, options map[string]interface{}) (*model.TicketComment, error)
	CountTicketComment(ctx context.Context, options map[string]interface{}) int64
	UpdateTicketCommentPartial(ctx context.Context, id primitive.ObjectID, field map[string]interface{}) error

	// Superuser
	FetchOneSuperuser(ctx context.Context, options map[string]interface{}) (*model.Superuser, error)

	// Attachment
	FetchAttachmentList(ctx context.Context, options map[string]interface{}) (*mongo.Cursor, error)
	CreateAttachment(ctx context.Context, ticket *model.Attachment) (err error)
	FetchOneAttachment(ctx context.Context, options map[string]interface{}) (*model.Attachment, error)
	UpdateManyAttachmentPartial(ctx context.Context, ids []primitive.ObjectID, field map[string]interface{}) error
	UpdateAttachmentPartial(ctx context.Context, id primitive.ObjectID, field map[string]interface{}) error

	// Media
	FetchMediaList(ctx context.Context, options map[string]interface{}) (*mongo.Cursor, error)
	CreateMedia(ctx context.Context, ticket *model.Media) (err error)
	FetchOneMedia(ctx context.Context, options map[string]interface{}) (*model.Media, error)
	UpdateManyMediaPartial(ctx context.Context, ids []primitive.ObjectID, field map[string]interface{}) error
	UpdateMediaPartial(ctx context.Context, id primitive.ObjectID, field map[string]interface{}) error

	// Ticket Timelogs
	FetchTicketTimelogsList(ctx context.Context, options map[string]interface{}) (*mongo.Cursor, error)
	CountTicketTimelogs(ctx context.Context, options map[string]interface{}) int64
	CreateTicketTimelogs(ctx context.Context, ticket *model.TicketTimeLogs) (err error)
	FetchOneTicketlogs(ctx context.Context, options map[string]interface{}) (*model.TicketTimeLogs, error)
	UpdateTicketlogs(ctx context.Context, ticket *model.TicketTimeLogs) (err error)

	// Superadmin
	FetchSuperadminList(ctx context.Context, options map[string]interface{}) (*mongo.Cursor, error)
	FetchOneSuperadmin(ctx context.Context, options map[string]interface{}) (*model.Superadmin, error)

	//hour package
	FetchHourPackageList(ctx context.Context, options map[string]interface{}) (*mongo.Cursor, error)
	FetchOneHourPackage(ctx context.Context, options map[string]interface{}) (*model.HourPackage, error)
	CountHourPackage(ctx context.Context, options map[string]interface{}) (count int64)
	CreateHourPackage(ctx context.Context, packages *model.HourPackage) (err error)
	UpdateHourPackage(ctx context.Context, packages *model.HourPackage) (err error)
	DeleteHourPackage(ctx context.Context, packages *model.HourPackage) (err error)
	UpdatePartialHourPackage(ctx context.Context, options map[string]interface{}, field map[string]interface{}) (err error)

	// Order
	FetchOrderList(ctx context.Context, options map[string]interface{}) (*mongo.Cursor, error)
	CountOrder(ctx context.Context, options map[string]interface{}) int64
	FetchOneOrder(ctx context.Context, options map[string]interface{}) (row *model.Order, err error)
	CreateOrder(ctx context.Context, order *model.Order) (err error)
	UpdateOneOrder(ctx context.Context, order *model.Order) (err error)

	// Customer Balance History
	CreateCustomerBalanceHistory(ctx context.Context, row *model.CustomerBalanceHistory) (err error)

	//Customer Subscription
	FetchCustomerSubscriptionList(ctx context.Context, options map[string]interface{}, withOptions bool) (*mongo.Cursor, error)
	CountCustomerSubscription(ctx context.Context, options map[string]interface{}) int64
	FetchOneCustomerSubscription(ctx context.Context, options map[string]interface{}) (*model.CustomerSubscription, error)
	CreateCustomerSubscription(ctx context.Context, row *model.CustomerSubscription) (err error)
	UpdateOneCustomerSubscription(ctx context.Context, customerSubscription *model.CustomerSubscription) (err error)
	UpdateManyPartialCustomerSubscription(ctx context.Context, ids []primitive.ObjectID, field map[string]interface{}) (err error)
	UpdatePartialCustomerSubscription(ctx context.Context, options, field map[string]interface{}) (err error)

	// Project
	FetchProjectList(ctx context.Context, options map[string]interface{}) (*mongo.Cursor, error)
	CountProject(ctx context.Context, options map[string]interface{}) int64
	FetchOneProject(ctx context.Context, options map[string]interface{}) (*model.Project, error)
	CreateProject(ctx context.Context, project *model.Project) (err error)
	UpdateOneProject(ctx context.Context, project *model.Project) (err error)

	// Ticket Category
	FetchTicketCategoryList(ctx context.Context, options map[string]interface{}) (*mongo.Cursor, error)
	FetchOneTicketCategory(ctx context.Context, options map[string]interface{}) (row *model.TicketCategory, err error)
	CountTicketCategory(ctx context.Context, options map[string]interface{}) (total int64)
	CreateTicketCategory(ctx context.Context, ticketsCategory *model.TicketCategory) (err error)
	UpdateOneTicketCategory(ctx context.Context, ticketsCategory *model.TicketCategory) (err error)

	// Server Product
	FetchServerPackageList(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error)
	FetchOneServerPackage(ctx context.Context, options map[string]interface{}) (row *model.ServerPackage, err error)
	CountServerPackage(ctx context.Context, options map[string]interface{}) (total int64)
	CreateServerPackage(ctx context.Context, ServerPackages *model.ServerPackage) (err error)
	UpdateServerPackage(ctx context.Context, ServerPackages *model.ServerPackage) (err error)
	DeleteServerPackage(ctx context.Context, ServerPackages *model.ServerPackage) (err error)
	UpdatePartialServerPackage(ctx context.Context, options map[string]interface{}, field map[string]interface{}) (err error)
}
