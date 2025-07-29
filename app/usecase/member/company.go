package usecase_member

import (
	"app/domain/model"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
)

func (u *appUsecase) GetCompanyDetailByDomain(ctx context.Context, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get id
	domain, ok := options["domain"].(string)
	if !ok {
		return response.Error(http.StatusBadRequest, "Invalid domain type")
	}

	company := &model.Company{}

	if u.redisRepo.Enabled() {
		if strval, err := u.redisRepo.Get(ctx, domain); err == nil {
			if err = json.Unmarshal(strval, &company); err == nil {
				// from cache
				return response.Success(company)
			}
		}
	}

	// check company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		// "fullUrl": domain,
		"subdomain": domain,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if company == nil {
		return response.Error(http.StatusBadRequest, "Company not found")
	}

	// set cache
	if u.redisRepo.Enabled() {
		ttl := 5 * time.Minute
		byteData, _ := json.Marshal(company)
		u.redisRepo.Set(ctx, domain, byteData, &ttl)
	}

	return response.Success(company)
}
