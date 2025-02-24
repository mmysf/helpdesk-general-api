package usecase_member

import (
	mongorepo "app/app/repository/mongo"
	"app/domain"
	"context"
	"time"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
)

type appUsecase struct {
	mongodbRepo    mongorepo.MongoDBRepo
	contextTimeout time.Duration
}

type RepoInjection struct {
	MongoDBRepo mongorepo.MongoDBRepo
}

func NewAppUsecase(r RepoInjection, timeout time.Duration) AppUsecase {
	return &appUsecase{
		mongodbRepo:    r.MongoDBRepo,
		contextTimeout: timeout,
	}
}

type AppUsecase interface {
	Login(ctx context.Context, payload domain.LoginRequest) response.Base
	Register(ctx context.Context, payload domain.RegisterRequest) response.Base
	GetMe(ctx context.Context, claim domain.JWTClaimUser) response.Base
}
