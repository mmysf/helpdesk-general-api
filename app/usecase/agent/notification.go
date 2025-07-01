package usecase_agent

import (
	"app/domain"
	"app/domain/model"
	"app/helpers"
	"context"
	"net/url"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"

	yurekaHelpers "github.com/Yureka-Teknologi-Cipta/yureka/helpers"
)

func (u *agentUsecase) GetNotificationList(ctx context.Context, claim domain.JWTClaimAgent, query url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, limit, offset := yurekaHelpers.GetLimitOffset(query)

	fetchOptions := map[string]interface{}{
		"limit":     limit,
		"offset":    offset,
		"companyID": claim.CompanyID,
		"userRole":  model.AgentRole,
		"sort":      "createdAt",
		"dir":       "desc",
		"type":   "ticketUpdated",
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

	// count total notifications
	totalDocuments := u.mongodbRepo.CountNotification(ctx, fetchOptions)
	
	// count unread notifications
	unreadOptions := map[string]interface{}{
		"companyID": claim.CompanyID,
		"userRole":  model.AgentRole,
		"isRead":    false,
	}
	unreadCount := u.mongodbRepo.CountNotification(ctx, unreadOptions)

	if totalDocuments == 0 {
		return response.Success(domain.ResponseList{
			List: response.List{
				List:  []interface{}{},
				Page:  page,
				Limit: limit,
				Total: totalDocuments,
			},
			TotalPage:   0,
			UnreadCount: unreadCount,
		})
	}

	// fetch notifications
	cur, err := u.mongodbRepo.FetchNotificationList(ctx, fetchOptions)
	if err != nil {
		return response.Error(500, err.Error())
	}
	defer cur.Close(ctx)

	list := make([]interface{}, 0)
	for cur.Next(ctx) {
		var notification model.Notification
		if err := cur.Decode(&notification); err != nil {
			return response.Error(500, err.Error())
		}
		list = append(list, notification)
	}

	return response.Success(domain.ResponseList{
		List: response.List{
			List:  list,
			Page:  page,
			Limit: limit,
			Total: totalDocuments,
		},
		TotalPage:   helpers.GetTotalPage(totalDocuments, limit),
		UnreadCount: unreadCount,
	})
}

func (u *agentUsecase) GetNotificationDetail(ctx context.Context, claim domain.JWTClaimAgent, id string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Fetch notification dengan filter company
	data, err := u.mongodbRepo.FetchOneNotification(ctx, map[string]interface{}{
		"id":        id,
		"companyID": claim.CompanyID, // Security: pastikan hanya bisa akses notif company sendiri
	})
	if err != nil {
		return response.Error(500, err.Error())
	}

	if data == nil {
		return response.Error(404, "notification not found")
	}

	// Mark as read
	data.IsRead = true
	if err := u.mongodbRepo.UpdateNotification(ctx, data); err != nil {
		return response.Error(500, err.Error())
	}

	return response.Success(data)
}

func (u *agentUsecase) ReadAllNotification(ctx context.Context, claim domain.JWTClaimAgent) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return response.Success(map[string]interface{}{
		"message": "All notifications marked as read",
	})
}

func (u *agentUsecase) GetNotificationCount(ctx context.Context, claim domain.JWTClaimAgent) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	unread := u.mongodbRepo.CountNotification(ctx, map[string]interface{}{
		"companyID": claim.CompanyID,
		"userRole":  model.AgentRole,
		"isRead":    false,
	})
	
	total := u.mongodbRepo.CountNotification(ctx, map[string]interface{}{
		"companyID": claim.CompanyID,
		"userRole":  model.AgentRole,
	})

	resp := map[string]interface{}{
		"unread": unread,
		"total":  total,
	}

	return response.Success(resp)
}