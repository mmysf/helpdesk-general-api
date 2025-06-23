package usecase_superadmin

import (
	"app/domain"
	"app/domain/model"
	"app/helpers"
	"context"
	"net/url"

	yurekahelpers "github.com/Yureka-Teknologi-Cipta/yureka/helpers"
	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/sirupsen/logrus"
)

func (u *superadminUsecase) GetDataDashboard(ctx context.Context, claim domain.JWTClaimSuperadmin) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	totalClient := u.mongodbRepo.CountCompany(ctx, map[string]interface{}{})
	totalOrder := u.mongodbRepo.CountOrder(ctx, map[string]interface{}{})
	totalCustomer := u.mongodbRepo.CountCustomer(ctx, map[string]interface{}{})
	totalNewTicket := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		"status": []string{"open"},
	})
	totalAgent := u.mongodbRepo.CountAgent(ctx, map[string]interface{}{})

	result := map[string]interface{}{
		"totalClient":    totalClient,
		"totalCustomer":  totalCustomer,
		"totalNewTicket": totalNewTicket,
		"totalAgent":     totalAgent,
		"totalOrder":     totalOrder,
	}

	return response.Success(result)
}

func (u *superadminUsecase) GetHourPackagesDashboard(ctx context.Context, claim domain.JWTClaimSuperadmin, query url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, limit, offset := yurekahelpers.GetLimitOffset(query)

	fetchOptions := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
		"status": "active",
	}

	//filtering
	if query.Get("categoryName") != "" {
		fetchOptions["categoryName"] = query.Get("categoryName")
	}

	// count first
	totalpackages := u.mongodbRepo.CountHourPackage(ctx, fetchOptions)

	if totalpackages == 0 {
		return response.Success(domain.ResponseList{
			List: response.List{
				List:  []interface{}{},
				Page:  page,
				Limit: limit,
				Total: totalpackages,
			},
			TotalPage: helpers.GetTotalPage(totalpackages, limit),
		})
	}

	// check hour package list
	packages, err := u.mongodbRepo.FetchHourPackageList(ctx, fetchOptions)

	if err != nil {
		return response.Success(domain.ResponseList{
			List: response.List{
				List:  []interface{}{},
				Page:  page,
				Limit: limit,
				Total: totalpackages,
			},
			TotalPage: 0,
		})
	}

	defer packages.Close(ctx)

	list := make([]interface{}, 0)
	for packages.Next(ctx) {
		row := model.HourPackage{}
		err := packages.Decode(&row)
		if err != nil {
			logrus.Error("Topup Decode ", err)
			return response.Success(response.List{
				List:  []interface{}{},
				Page:  page,
				Limit: limit,
				Total: totalpackages,
			})
		}
		list = append(list, row)
	}

	return response.Success(domain.ResponseList{
		List: response.List{
			List:  list,
			Page:  page,
			Limit: limit,
			Total: totalpackages,
		},
		TotalPage: helpers.GetTotalPage(totalpackages, limit),
	})
}
