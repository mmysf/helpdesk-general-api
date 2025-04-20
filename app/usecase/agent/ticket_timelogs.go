package usecase_agent

import (
	"app/domain"
	"app/domain/model"
	"app/helpers"
	"context"
	"net/url"
	"time"

	yurekahelpers "github.com/Yureka-Teknologi-Cipta/yureka/helpers"
	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/sirupsen/logrus"
)

func (u *agentUsecase) GetTicketTimeLogsList(ctx context.Context, claim domain.JWTClaimAgent, query url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, limit, offset := yurekahelpers.GetLimitOffset(query)

	fetchOptions := map[string]interface{}{
		"limit":     limit,
		"offset":    offset,
		"companyID": claim.CompanyID,
	}

	// filtering
	if query.Get("sort") != "" {
		fetchOptions["sort"] = query.Get("sort")
	}

	if query.Get("dir") != "" {
		fetchOptions["dir"] = query.Get("dir")
	}

	if query.Get("ticketId") != "" {
		fetchOptions["ticketId"] = query.Get("ticketId")
	}

	if query.Get("customerID") != "" {
		fetchOptions["customerID"] = query.Get("customerID")
	}

	if query.Get("projectID") != "" {
		fetchOptions["projectID"] = query.Get("projectID")
	}

	// count first
	totalDocuments := u.mongodbRepo.CountTicketTimelogs(ctx, fetchOptions)

	if totalDocuments == 0 {
		return response.Success(domain.ResponseList{
			List: response.List{
				List:  []interface{}{},
				Page:  page,
				Limit: limit,
				Total: totalDocuments,
			},
			TotalPage: 0,
		})
	}

	// check ticket list
	cur, err := u.mongodbRepo.FetchTicketTimelogsList(ctx, fetchOptions)

	if err != nil {
		return response.Success(domain.ResponseList{
			List: response.List{
				List:  []interface{}{},
				Page:  page,
				Limit: limit,
				Total: totalDocuments,
			},
			TotalPage: 0,
		})
	}

	defer cur.Close(ctx)

	list := make([]interface{}, 0)
	for cur.Next(ctx) {
		row := model.TicketTimeLogs{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("Time log Decode ", err)
			return response.Success(
				domain.ResponseList{
					List: response.List{
						List:  []interface{}{},
						Page:  page,
						Limit: limit,
						Total: totalDocuments,
					},
					TotalPage: 0,
				},
			)
		}

		//adding duration active per pause
		for i, pause := range row.PauseHistory {
			pausedAt := pause.PausedAt
			var durationActive time.Duration

			//first pause case
			if i == 0 {
				durationActive = pausedAt.Sub(*row.StartAt)
			} else {
				if row.PauseHistory[i-1].ResumedAt != nil {
					lastResume := *row.PauseHistory[i-1].ResumedAt
					durationActive = pausedAt.Sub(lastResume)
				} else {
					logrus.Warn("ResumedAt is nil for previous pause")
					continue
				}
			}

			//Add duration active to pause history data
			row.PauseHistory[i].DurationActive = int(durationActive.Seconds())
		}

		//count active duration after last pause
		if len(row.PauseHistory) > 0 {
			lastPause := row.PauseHistory[len(row.PauseHistory)-1]
			if lastPause.ResumedAt != nil && row.EndAt != nil {
				endAt := row.EndAt
				durationAfterLastPause := endAt.Sub(*lastPause.ResumedAt)
				row.PauseHistory = append(row.PauseHistory, model.PauseHistory{
					DurationActive: int(durationAfterLastPause.Seconds()),
				})
			}
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
