package usecase_member

import (
	"app/domain"
	"app/domain/model"
	"app/helpers"
	"context"
	"net/http"
	"net/url"
	"time"

	yurekahelpers "github.com/Yureka-Teknologi-Cipta/yureka/helpers"
	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/sirupsen/logrus"
)

func (u *appUsecase) GetCustomerSubscriptionList(ctx context.Context, claim domain.JWTClaimUser, query url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, limit, offset := yurekahelpers.GetLimitOffset(query)

	fetchOptions := map[string]interface{}{
		"limit":      limit,
		"offset":     offset,
		"customerID": claim.UserID,
	}

	// filtering
	if query.Get("sort") != "" {
		fetchOptions["sort"] = query.Get("sort")
	}
	if query.Get("dir") != "" {
		fetchOptions["dir"] = query.Get("dir")
	}
	if query.Get("status") != "" {
		fetchOptions["status"] = query.Get("status")
	}
	if query.Get("orderType") != "" {
		fetchOptions["orderType"] = query.Get("orderType")
	}

	// count first
	totalDocuments := u.mongodbRepo.CountCustomerSubscription(ctx, fetchOptions)

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

	// check customer subscription list
	cur, err := u.mongodbRepo.FetchCustomerSubscriptionList(ctx, fetchOptions, true)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	defer cur.Close(ctx)

	list := make([]interface{}, 0)
	for cur.Next(ctx) {
		row := model.CustomerSubscription{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("Topup Decode ", err)
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

func (u *appUsecase) GetCustomerSubscriptionDetail(ctx context.Context, claim domain.JWTClaimUser, customerSubscriptionID string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check customer subscription
	customerSubscription, err := u.mongodbRepo.FetchOneCustomerSubscription(ctx, map[string]interface{}{
		"id": customerSubscriptionID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if customerSubscription == nil {
		return response.Error(http.StatusBadRequest, "customer subscription not found")
	}

	// update order status if expired
	if customerSubscription.Status == model.Active && customerSubscription.ExpiredAt.Before(time.Now()) {
		customerSubscription.Status = model.Expired
		customerSubscription.UpdatedAt = time.Now()

		// update in background
		go func() {
			u.mongodbRepo.UpdateOneCustomerSubscription(context.Background(), customerSubscription)
		}()

	}

	return response.Success(customerSubscription)
}
