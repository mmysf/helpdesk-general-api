package usecase_superadmin

import (
	mongorepo "app/app/repository/mongo"
	redisrepo "app/app/repository/redis"
	s3repo "app/app/repository/s3"
	"app/domain"
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
)

type superadminUsecase struct {
	mongodbRepo    mongorepo.MongoDBRepo
	contextTimeout time.Duration
	redisRepo      redisrepo.RedisRepo
	s3Repo         s3repo.S3Repo
}

type RepoInjection struct {
	MongoDBRepo mongorepo.MongoDBRepo
	Redis       redisrepo.RedisRepo
	S3Repo      s3repo.S3Repo
}

func NewAppSuperadminUsecase(r RepoInjection, timeout time.Duration) SuperadminUsecase {
	return &superadminUsecase{
		mongodbRepo:    r.MongoDBRepo,
		contextTimeout: timeout,
		redisRepo:      r.Redis,
		s3Repo:         r.S3Repo,
	}
}

type SuperadminUsecase interface {
	// auth
	Login(ctx context.Context, payload domain.SuperadminLoginRequest) response.Base
	GetMe(ctx context.Context, claim domain.JWTClaimSuperadmin) response.Base

	// ticket
	GetTotalTicket(ctx context.Context, claim domain.JWTClaimSuperadmin) response.Base
	GetTicketList(ctx context.Context, claim domain.JWTClaimSuperadmin, query url.Values) response.Base
	GetTicketDetail(ctx context.Context, ticketId string) response.Base
	AssignAgent(ctx context.Context, claim domain.JWTClaimSuperadmin, ticketId string, payload domain.AssignAgentRequest) response.Base
	GetDataClientTicket(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base
	GetAverageDurationClient(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base
	PauseLoggingTicket(ctx context.Context, claim domain.JWTClaimSuperadmin, payload domain.LoggingTicketRequest) response.Base
	ResumeLoggingTicket(ctx context.Context, claim domain.JWTClaimSuperadmin, payload domain.LoggingTicketRequest) response.Base

	// ticket comment
	GetTicketCommentList(ctx context.Context, claim domain.JWTClaimSuperadmin, ticketId string, query url.Values) response.Base
	GetTicketCommentDetail(ctx context.Context, claim domain.JWTClaimSuperadmin, commentId string) response.Base
	CreateTicketComment(ctx context.Context, claim domain.JWTClaimSuperadmin, payload domain.SuperadminTicketCommentRequest) response.Base

	// order
	GetOrderList(ctx context.Context, claim domain.JWTClaimSuperadmin, query url.Values) response.Base
	GetOrderDetail(ctx context.Context, orderID string) response.Base
	UploadAttachmentOrder(ctx context.Context, claim domain.JWTClaimSuperadmin, payload domain.UploadAttachment, request *http.Request) response.Base
	UpdateManualPayment(ctx context.Context, claim domain.JWTClaimSuperadmin, orderID string, payload domain.UpdateManualPaymentRequest) response.Base

	// dashboard
	GetDataDashboard(ctx context.Context, claim domain.JWTClaimSuperadmin) response.Base
	GetHourPackagesDashboard(ctx context.Context, claim domain.JWTClaimSuperadmin, query url.Values) response.Base

	// customer
	GetCustomers(ctx context.Context, claim domain.JWTClaimSuperadmin, query url.Values) response.Base
	CreateCustomer(ctx context.Context, claim domain.JWTClaimSuperadmin, payload domain.AccountRequest) response.Base
	GetCustomerDetail(ctx context.Context, claim domain.JWTClaimSuperadmin, customerId string) response.Base
	UpdateCustomer(ctx context.Context, claim domain.JWTClaimSuperadmin, customerId string, payload domain.AccountRequest) response.Base
	DeleteCustomer(ctx context.Context, claim domain.JWTClaimSuperadmin, customerId string) response.Base
	ResetPasswordCustomer(ctx context.Context, claim domain.JWTClaimSuperadmin, customerId string) response.Base
	ImportCustomer(ctx context.Context, claim domain.JWTClaimSuperadmin, request *http.Request) response.Base

	// hour package
	GetHourPackages(ctx context.Context, claim domain.JWTClaimSuperadmin, query url.Values) response.Base
	CreateHourPackage(ctx context.Context, claim domain.JWTClaimSuperadmin, payload domain.HourPackageRequest) response.Base
	GetHourPackageDetail(ctx context.Context, productId string) response.Base
	UpdateHourPackage(ctx context.Context, claim domain.JWTClaimSuperadmin, packageId string, payload domain.HourPackageUpdate) response.Base
	UpdateStatusHourPackage(ctx context.Context, claim domain.JWTClaimSuperadmin, packageId string, payload domain.HourPackageStatusUpdate) response.Base
	DeleteHourPackage(ctx context.Context, claim domain.JWTClaimSuperadmin, packageId string) response.Base

	// agent
	GetAgents(ctx context.Context, claim domain.JWTClaimSuperadmin, query url.Values) response.Base
	CreateAgent(ctx context.Context, claim domain.JWTClaimSuperadmin, payload domain.AccountRequest) response.Base
	GetAgentDetail(ctx context.Context, claim domain.JWTClaimSuperadmin, agentId string) response.Base
	UpdateAgent(ctx context.Context, claim domain.JWTClaimSuperadmin, agentId string, payload domain.AccountRequest) response.Base
	DeleteAgent(ctx context.Context, claim domain.JWTClaimSuperadmin, agentId string) response.Base
	ResetPasswordAgent(ctx context.Context, claim domain.JWTClaimSuperadmin, agentId string) response.Base

	// company
	UploadCompanyLogo(ctx context.Context, claim domain.JWTClaimSuperadmin, payload domain.UploadAttachment, request *http.Request) response.Base
	GetCompanyList(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base
	GetCompanyDetail(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base
	CreateCompany(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base
	UpdateCompany(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base
	DeleteCompany(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base

	// company product
	UploadCompanyProductLogo(ctx context.Context, claim domain.JWTClaimSuperadmin, payload domain.UploadAttachment, request *http.Request) response.Base
	GetCompanyProductList(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base
	GetCompanyProductDetail(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base
	CreateCompanyProduct(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base
	UpdateCompanyProduct(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base
	DeleteCompanyProduct(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base

	// config
	GetConfig(ctx context.Context) response.Base

	// server package
	GetServerPackageList(ctx context.Context, claim domain.JWTClaimSuperadmin, query url.Values) response.Base
	CreateServerPackage(ctx context.Context, claim domain.JWTClaimSuperadmin, payload domain.ServerPackageRequest) response.Base
	GetServerPackageDetail(ctx context.Context, packageId string) response.Base
	UpdateServerPackage(ctx context.Context, claim domain.JWTClaimSuperadmin, serverpackageId string, payload domain.ServerPackageUpdate) response.Base
	UpdateStatusServerPackage(ctx context.Context, claim domain.JWTClaimSuperadmin, serverpackageId string, payload domain.ServerPackageStatusUpdate) response.Base
	DeleteServerPackage(ctx context.Context, claim domain.JWTClaimSuperadmin, serverpackageId string) response.Base
}
