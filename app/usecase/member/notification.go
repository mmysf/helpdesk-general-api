package usecase_member

import (
	"app/domain"
	"app/domain/model"
	"context"
	"net/url"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"

	yurekaHelpers "github.com/Yureka-Teknologi-Cipta/yureka/helpers"
)

func (u *appUsecase) GetNotificationList(ctx context.Context, claim domain.JWTClaimUser, query url.Values) response.Base {

	page, limit, offset := yurekaHelpers.GetLimitOffset(query)

	fetchOptions := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
		"page":   page,
		"sort":   "createdAt",
		"dir":    "desc",
	}

	// filtering
	if query.Get("sort") != "" {
		fetchOptions["sort"] = query.Get("sort")
	}

	if query.Get("dir") != "" {
		fetchOptions["dir"] = query.Get("dir")
	}

	if isLastTwo := query.Get("isLastTwo"); isLastTwo != "" {
		if isLastTwo == "1" {
			fetchOptions["isLastTwo"] = true
		} else {
			fetchOptions["isLastTwo"] = false
		}
	}

	if category := query.Get("category"); category != "" {
		fetchOptions["category"] = category
	}

	// count
	totalDocuments := u.mongodbRepo.CountNotification(ctx, fetchOptions)
	if totalDocuments == 0 {
		return response.Success(
			map[string]interface{}{
				"list":  []interface{}{},
				"meta":  nil,
				"limit": limit,
				"page":  page,
				"total": totalDocuments,
			},
		)
	}

	data, _ := u.mongodbRepo.FetchNotificationList(ctx, fetchOptions)

	defer data.Close(ctx)

	var r []interface{} = make([]interface{}, 0)
	for data.Next(ctx) {
		var t model.Notification
		if err := data.Decode(&t); err != nil {
			return response.Error(500, err.Error())
		}
		r = append(r, t)
	}

	meta := u.GetNotificationCount(ctx, claim)

	return response.Success(
		map[string]interface{}{
			"list":  r,
			"meta":  meta.Data,
			"limit": limit,
			"page":  page,
			"total": totalDocuments,
		},
	)
}

func (u *appUsecase) GetNotificationDetail(ctx context.Context, claim domain.JWTClaimUser, id string) response.Base {
	data, err := u.mongodbRepo.FetchOneNotification(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return response.Error(500, err.Error())
	}

	data.IsRead = true

	u.mongodbRepo.UpdateNotification(ctx, data)
	return response.Success(data)
}

func (u *appUsecase) ReadAllNotification(ctx context.Context, claim domain.JWTClaimUser) response.Base {
	u.mongodbRepo.ReadAllNotification(ctx, claim.UserID)
	return response.Success(nil)
}

func (u *appUsecase) GetNotificationCount(ctx context.Context, claim domain.JWTClaimUser) response.Base {
	unread := u.mongodbRepo.CountNotification(ctx, map[string]interface{}{
		"isRead": false,
	})
	readed := u.mongodbRepo.CountNotification(ctx, map[string]interface{}{
		"isRead": true,
	})

	resp := map[string]interface{}{
		"unread": unread,
		"readed": readed,
	}

	return response.Success(resp)
}
