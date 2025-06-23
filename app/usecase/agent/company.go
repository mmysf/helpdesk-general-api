package usecase_agent

import (
	"app/domain"
	"context"
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
)

func (u *agentUsecase) GetCompanyDetail(ctx context.Context, claim domain.JWTClaimAgent) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get id
	companyID := claim.CompanyID

	// check company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": companyID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if company == nil {
		return response.Error(http.StatusBadRequest, "CompanyProduct not found")
	}

	return response.Success(company)
}
