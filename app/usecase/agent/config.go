package usecase_agent

import (
	"context"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
)

func (u *agentUsecase) GetConfig(ctx context.Context) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get from config
	config := u._CacheConfig(ctx)

	return response.Success(config.ToPublic())
}
