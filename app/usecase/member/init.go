package usecase_member

import (
	mongorepo "app/app/repository/mongo"
	redisrepo "app/app/repository/redis"
	s3repo "app/app/repository/s3"
	xenditrepo "app/app/repository/xendit"
	"app/domain"
	"net/http"
	"net/url"

	"context"
	"time"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
)

type appUsecase struct {
	mongodbRepo    mongorepo.MongoDBRepo
	contextTimeout time.Duration
	redisRepo      redisrepo.RedisRepo
	s3Repo         s3repo.S3Repo
	xenditRepo     xenditrepo.XenditRepo
}

type RepoInjection struct {
	MongoDBRepo mongorepo.MongoDBRepo
	Redis       redisrepo.RedisRepo
	S3Repo      s3repo.S3Repo
	XenditRepo  xenditrepo.XenditRepo
}

func NewAppUsecase(r RepoInjection, timeout time.Duration) AppUsecase {
	return &appUsecase{
		mongodbRepo:    r.MongoDBRepo,
		contextTimeout: timeout,
		redisRepo:      r.Redis,
		s3Repo:         r.S3Repo,
		xenditRepo:     r.XenditRepo,
	}
}

type AppUsecase interface {
	// Auth
	Login(ctx context.Context, payload domain.LoginRequest) response.Base
	Register(ctx context.Context, payload domain.RegisterRequest) response.Base
	VerifyRegistration(ctx context.Context, payload domain.VerifyRegisterRequest) response.Base
	SendEmailPasswordReset(ctx context.Context, payload domain.EmailPasswordResetRequest) response.Base
	PasswordReset(ctx context.Context, payload domain.PasswordResetRequest) response.Base
	GetMe(ctx context.Context, claim domain.JWTClaimUser) response.Base
	RegisterB2B(ctx context.Context, payload domain.RegisterRequest) response.Base

	// Ticket
	GetTicketList(ctx context.Context, claim domain.JWTClaimUser, query url.Values) response.Base
	CreateTicket(ctx context.Context, claim domain.JWTClaimUser, payload domain.TicketRequest) response.Base
	GetTicketDetail(ctx context.Context, claim domain.JWTClaimUser, TicketID string) response.Base
	CloseTicket(ctx context.Context, claim domain.JWTClaimUser, payload domain.CloseTicketRequest) response.Base
	CloseTicketByEmail(ctx context.Context, payload domain.CloseTicketbyEmailRequest) response.Base
	ReopenTicket(ctx context.Context, claim domain.JWTClaimUser, payload domain.ReopenTicketRequest) response.Base
	CancelTicket(ctx context.Context, claim domain.JWTClaimUser, payload domain.CancelTicketRequest) response.Base

	// Tickect comment
	CreateTicketComment(ctx context.Context, claim domain.JWTClaimUser, paylaod domain.TicketCommentRequest) response.Base
	GetTicketCommentList(ctx context.Context, claim domain.JWTClaimUser, ticketId string, query url.Values) response.Base
	GetTicketCommentDetail(ctx context.Context, claim domain.JWTClaimUser, commentId string) response.Base

	// Product
	GetProductList(ctx context.Context, claim domain.JWTClaimUser, query url.Values) response.Base

	// Attachment
	UploadAttachment(ctx context.Context, claim domain.JWTClaimUser, payload domain.UploadAttachment, request *http.Request) response.Base
	GetAttachmentDetail(ctx context.Context, claim domain.JWTClaimUser, attachmentId string) response.Base

	// Dashboard
	GetTotalTicket(ctx context.Context, claim domain.JWTClaimUser) response.Base
	GetTotalTicketNow(ctx context.Context, claim domain.JWTClaimUser) response.Base
	GetDataDashboard(ctx context.Context, claim domain.JWTClaimUser, query url.Values) response.Base
	GetAverageDurationDashboard(ctx context.Context, claim domain.JWTClaimUser, query url.Values) response.Base

	// Order
	GetOrderList(ctx context.Context, claim domain.JWTClaimUser, query url.Values) response.Base
	GetOrderDetail(ctx context.Context, claim domain.JWTClaimUser, orderID string) response.Base
	CreateHourOrder(ctx context.Context, claim domain.JWTClaimUser, payload domain.OrderRequest) response.Base
	CreateServerOrder(ctx context.Context, claim domain.JWTClaimUser, payload domain.OrderRequest) response.Base
	ConfirmOrder(ctx context.Context, claim domain.JWTClaimUser, payload domain.ConfrimOrderRequest) response.Base
	UploadAttachmentOrder(ctx context.Context, claim domain.JWTClaimUser, payload domain.UploadAttachment, request *http.Request) response.Base

	// Hour Package
	GetHourPackageList(ctx context.Context, claim domain.JWTClaimUser, query url.Values) response.Base
	GetHourPackageDetail(ctx context.Context, claim domain.JWTClaimUser, packageID string) response.Base

	// Customer Subscription
	GetCustomerSubscriptionList(ctx context.Context, claim domain.JWTClaimUser, query url.Values) response.Base
	GetCustomerSubscriptionDetail(ctx context.Context, claim domain.JWTClaimUser, customerSubscriptionID string) response.Base

	// User
	GetUserList(ctx context.Context, claim domain.JWTClaimUser, options map[string]interface{}) response.Base
	GetUserDetail(ctx context.Context, claim domain.JWTClaimUser, options map[string]interface{}) response.Base
	CreateUser(ctx context.Context, claim domain.JWTClaimUser, options map[string]interface{}) response.Base
	UpdateUser(ctx context.Context, claim domain.JWTClaimUser, options map[string]interface{}) response.Base
	DeleteUser(ctx context.Context, claim domain.JWTClaimUser, options map[string]interface{}) response.Base

	// Project
	GetProjectList(ctx context.Context, claim domain.JWTClaimUser, query url.Values) response.Base
	GetProjectDetail(ctx context.Context, claim domain.JWTClaimUser, projectID string) response.Base
	CreateProject(ctx context.Context, claim domain.JWTClaimUser, payload domain.ProjectRequest) response.Base
	UpdateProject(ctx context.Context, claim domain.JWTClaimUser, projectID string, payload domain.ProjectRequest) response.Base
	DeleteProject(ctx context.Context, claim domain.JWTClaimUser, projectID string) response.Base

	//Company
	GetCompanyDetailByDomain(ctx context.Context, options map[string]interface{}) response.Base

	// Setting
	ChangePassword(ctx context.Context, claim domain.JWTClaimUser, payload domain.ChangePasswordRequest) response.Base

	// Ticket Category
	GetTicketCategoriesList(ctx context.Context, claim domain.JWTClaimUser, query url.Values) response.Base

	// server package
	GetServerPackageList(ctx context.Context, claim domain.JWTClaimUser, query url.Values) response.Base
	GetServerPackageDetail(ctx context.Context, packageId string) response.Base

	// config
	GetConfig(ctx context.Context) response.Base

	// Notification
	GetNotificationList(ctx context.Context, claim domain.JWTClaimUser, query url.Values) response.Base
	GetNotificationDetail(ctx context.Context, claim domain.JWTClaimUser, id string) response.Base
	ReadAllNotification(ctx context.Context, claim domain.JWTClaimUser) response.Base
	GetNotificationCount(ctx context.Context, claim domain.JWTClaimUser) response.Base
}
