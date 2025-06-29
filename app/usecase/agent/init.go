package usecase_agent

import (
	mongorepo "app/app/repository/mongo"
	redisrepo "app/app/repository/redis"
	s3Repo "app/app/repository/s3"
	"app/domain"
	"context"
	"net/http"
	"net/url"

	"time"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
)

type agentUsecase struct {
	mongodbRepo    mongorepo.MongoDBRepo
	contextTimeout time.Duration
	redisRepo      redisrepo.RedisRepo
	s3Repo         s3Repo.S3Repo
}

type RepoInjection struct {
	MongoDBRepo mongorepo.MongoDBRepo
	Redis       redisrepo.RedisRepo
	S3Repo      s3Repo.S3Repo
}

func NewAppAgentUsecase(r RepoInjection, timeout time.Duration) AgentUsecase {
	return &agentUsecase{
		mongodbRepo:    r.MongoDBRepo,
		contextTimeout: timeout,
		redisRepo:      r.Redis,
		s3Repo:         r.S3Repo,
	}
}

type AgentUsecase interface {
	// Auth
	Login(ctx context.Context, payload domain.LoginRequest) response.Base
	SendEmailPasswordReset(ctx context.Context, payload domain.EmailPasswordResetRequest) response.Base
	PasswordReset(ctx context.Context, payload domain.PasswordResetRequest) response.Base
	GetMe(ctx context.Context, claim domain.JWTClaimAgent) response.Base

	// Ticket
	GetTicketList(ctx context.Context, claim domain.JWTClaimAgent, query url.Values) response.Base
	GetMyTicketList(ctx context.Context, claim domain.JWTClaimAgent, query url.Values) response.Base
	GetTicketDetail(ctx context.Context, claim domain.JWTClaimAgent, ticketId string) response.Base
	CloseTicket(ctx context.Context, claim domain.JWTClaimAgent, payload domain.CloseTicketRequest) response.Base
	ReopenTicket(ctx context.Context, claim domain.JWTClaimAgent, payload domain.ReopenTicketRequest) response.Base
	StartLoggingTicket(ctx context.Context, claim domain.JWTClaimAgent, payload domain.LoggingTicketRequest) response.Base
	StopLoggingTicket(ctx context.Context, claim domain.JWTClaimAgent, payload domain.LoggingTicketRequest) response.Base
	PauseLoggingTicket(ctx context.Context, claim domain.JWTClaimAgent, payload domain.LoggingTicketRequest) response.Base
	ResumeLoggingTicket(ctx context.Context, claim domain.JWTClaimAgent, payload domain.LoggingTicketRequest) response.Base
	EditTimeTrack(ctx context.Context, claim domain.JWTClaimAgent, ticketId string, payload domain.TimeTrackRequest) response.Base
	AssignTicketToMe(ctx context.Context, claim domain.JWTClaimAgent, ticketId string) response.Base

	// Ticket Comment
	CreateTicketComment(ctx context.Context, claim domain.JWTClaimAgent, payload domain.TicketCommentRequest) response.Base
	GetTicketCommentList(ctx context.Context, claim domain.JWTClaimAgent, ticketId string, query url.Values) response.Base
	GetTicketCommentDetail(ctx context.Context, claim domain.JWTClaimAgent, commentId string) response.Base

	// Product
	GetProductList(ctx context.Context, claim domain.JWTClaimAgent, query url.Values) response.Base

	// Attachment
	UploadAttachment(ctx context.Context, claim domain.JWTClaimAgent, payload domain.UploadAttachment, request *http.Request) response.Base
	GetAttachmentDetail(ctx context.Context, claim domain.JWTClaimAgent, attachmentId string) response.Base

	// Dashboard
	GetTotalTicket(ctx context.Context, claim domain.JWTClaimAgent) response.Base
	GetTotalTicketNow(ctx context.Context, claim domain.JWTClaimAgent) response.Base
	GetDataDashboard(ctx context.Context, claim domain.JWTClaimAgent, query url.Values) response.Base

	// Ticket Timelogs
	GetTicketTimeLogsList(ctx context.Context, claim domain.JWTClaimAgent, query url.Values) response.Base

	// Export To CSV
	ExportTicketsToCSV(ctx context.Context, claim domain.JWTClaimAgent, query url.Values, w http.ResponseWriter) response.Base

	// total ticket customer
	GetTotalTicketCustomer(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base
	GetDataCustomerTicket(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base

	// company product
	UploadCompanyProductLogo(ctx context.Context, claim domain.JWTClaimAgent, payload domain.UploadAttachment, request *http.Request) response.Base
	GetCompanyProductList(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base
	GetCompanyProductDetail(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base
	CreateCompanyProduct(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base
	UpdateCompanyProduct(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base
	DeleteCompanyProduct(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base

	// Customer
	GetCustomerList(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base
	GetCustomerDetail(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base
	CreateCustomer(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base
	UpdateCustomer(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base
	DeleteCustomer(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base

	// company
	GetCompanyDetail(ctx context.Context, claim domain.JWTClaimAgent) response.Base

	// setting
	ChangePassword(ctx context.Context, claim domain.JWTClaimAgent, payload domain.ChangePasswordRequest) response.Base
	ChangeDomain(ctx context.Context, claim domain.JWTClaimAgent, payload domain.ChangeDomainRequest) response.Base
	UpdateProfile(ctx context.Context, claim domain.JWTClaimAgent, payload domain.UpdateProfileRequest) response.Base
	ChangeColor(ctx context.Context, claim domain.JWTClaimAgent, payload domain.ChangeColorMode) response.Base
	UploadAgentProfilePicture(ctx context.Context, claim domain.JWTClaimAgent, payload domain.UploadAttachment, request *http.Request) response.Base

	// Agent
	GetAgentList(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base
	GetAgentDetail(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base
	CreateAgent(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base
	UpdateAgent(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base
	DeleteAgent(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base

	// Config
	GetConfig(ctx context.Context) response.Base

	// Ticket Category
	GetTicketCategoriesList(ctx context.Context, claim domain.JWTClaimAgent, query url.Values) response.Base
	GetTicketCategoryDetail(ctx context.Context, claim domain.JWTClaimAgent, id string) response.Base
	CreateTicketCategory(ctx context.Context, claim domain.JWTClaimAgent, payload domain.TicketCategoryRequest) response.Base
	UpdateTicketCategory(ctx context.Context, claim domain.JWTClaimAgent, id string, payload domain.TicketCategoryRequest) response.Base
	DeleteTicketCategory(ctx context.Context, claim domain.JWTClaimAgent, id string) response.Base

	// Notification
	GetNotificationList(ctx context.Context, claim domain.JWTClaimAgent, query url.Values) response.Base
	GetNotificationDetail(ctx context.Context, claim domain.JWTClaimAgent, id string) response.Base
	ReadAllNotification(ctx context.Context, claim domain.JWTClaimAgent) response.Base
	GetNotificationCount(ctx context.Context, claim domain.JWTClaimAgent) response.Base
}
