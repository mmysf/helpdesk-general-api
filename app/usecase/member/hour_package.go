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
	"github.com/sirupsen/logrus"
)

func (u *appUsecase) GetHourPackageList(ctx context.Context, claim domain.JWTClaimUser, query url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, limit, offset := yurekahelpers.GetLimitOffset(query)

	fetchOptions := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
		"status": string(model.HourPackageActive),
	}

	// filtering
	if query.Get("sort") != "" {
		fetchOptions["sort"] = query.Get("sort")
	}
	if query.Get("dir") != "" {
		fetchOptions["dir"] = query.Get("dir")
	}
	if query.Get("categoryID") != "" {
		fetchOptions["categoryID"] = query.Get("categoryID")
	}

	// count first
	totalDocuments := u.mongodbRepo.CountHourPackage(ctx, fetchOptions)

	if totalDocuments == 0 {
		return response.Success(domain.ResponseList{
			List: response.List{
				List:  []interface{}{},
				Page:  page,
				Limit: limit,
				Total: totalDocuments,
			},
			TotalPage: helpers.GetTotalPage(totalDocuments, limit),
		})
	}

	// check hour package list
	cur, err := u.mongodbRepo.FetchHourPackageList(ctx, fetchOptions)

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	defer cur.Close(ctx)

	list := make([]interface{}, 0)
	for cur.Next(ctx) {
		row := model.HourPackage{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("HourPackage Decode ", err)
			return response.Error(http.StatusInternalServerError, err.Error())
		}

		list = append(list, row)
	}

	return response.Success(domain.ResponseList{
		List: response.List{
			List:  list,
			Page:  page,
			Limit: limit,
			Total: totalDocuments,
		},
		TotalPage: helpers.GetTotalPage(totalDocuments, limit),
	})
}

func (u *appUsecase) GetHourPackageDetail(ctx context.Context, claim domain.JWTClaimUser, packageID string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check hour package
	oneHourPackage, err := u.mongodbRepo.FetchOneHourPackage(ctx, map[string]interface{}{
		"id": packageID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if oneHourPackage == nil {
		return response.Error(http.StatusBadRequest, "hour package not found")
	}

	return response.Success(oneHourPackage)
}
