package usecase_member

import (
	"context"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
)

func (u *appUsecase) GetConfig(ctx context.Context) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get from config
	config := u._CacheConfig(ctx)

	return response.Success(config.ToPublic())
}
