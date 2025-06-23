package usecase_member

import (
	"app/domain"
	"app/domain/model"
	"app/helpers"
	"context"
	"net/http"
	"net/url"

	yurekahelpers "github.com/Yureka-Teknologi-Cipta/yureka/helpers"
	"github.com/Yureka-Teknologi-Cipta/yureka/response"
)

func (u *appUsecase) GetServerPackageList(ctx context.Context, claim domain.JWTClaimUser, query url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, limit, offset := yurekahelpers.GetLimitOffset(query)

	fetchOptions := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
		"status": string(model.ServerPackageActive),
	}

	// filtering
	if query.Get("sort") != "" {
		fetchOptions["sort"] = query.Get("sort")
	}
	if query.Get("dir") != "" {
		fetchOptions["dir"] = query.Get("dir")
	}
	if query.Get("dir") != "" {
		fetchOptions["dir"] = query.Get("dir")
	}
	if query.Get("q") != "" {
		fetchOptions["q"] = query.Get("q")
	}

	// count first
	totalServerPackage := u.mongodbRepo.CountServerPackage(ctx, fetchOptions)
	if totalServerPackage == 0 {
		return response.Success(domain.ResponseList{
			List: response.List{
				List:  []interface{}{},
				Page:  page,
				Limit: limit,
				Total: totalServerPackage,
			},
			TotalPage: helpers.GetTotalPage(totalServerPackage, limit),
		})
	}

	// check package list
	serverPackages, err := u.mongodbRepo.FetchServerPackageList(ctx, fetchOptions)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	list := make([]interface{}, 0)
	for serverPackages.Next(ctx) {
		row := model.ServerPackage{}
		err := serverPackages.Decode(&row)
		if err != nil {
			return response.Error(http.StatusInternalServerError, err.Error())
		}
		list = append(list, row)
	}

	return response.Success(domain.ResponseList{
		List: response.List{
			List:  list,
			Page:  page,
			Limit: limit,
			Total: totalServerPackage,
		},
		TotalPage: helpers.GetTotalPage(totalServerPackage, limit),
	})
}

func (u *appUsecase) GetServerPackageDetail(ctx context.Context, packageId string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validating
	if packageId == "" {
		return response.Error(http.StatusBadRequest, "packages id is required")
	}

	serverPackages, err := u.mongodbRepo.FetchOneServerPackage(ctx, map[string]interface{}{
		"id":     packageId,
		"status": string(model.ServerPackageActive),
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if serverPackages == nil {
		return response.Error(http.StatusBadRequest, "server package not found")
	}

	return response.Success(serverPackages)
}
