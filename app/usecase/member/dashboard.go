package usecase_member

import (
	"app/domain"
	"app/domain/model"
	"context"
	"maps"
	"net/url"
	"time"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/sirupsen/logrus"
)

func (u *appUsecase) GetTotalTicket(ctx context.Context, claim domain.JWTClaimUser) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// filter
	fetchOption := map[string]interface{}{
		"companyID": claim.CompanyID,
	}

	// filter user if b2c company
	if claim.Company.Type == "B2B" {
		fetchOption["customerID"] = claim.UserID
	}

	// option total
	total := u.mongodbRepo.CountTicket(ctx, fetchOption)

	// option open
	fetchOptionOpen := maps.Clone(fetchOption)
	fetchOptionOpen["status"] = []string{"open"}
	totalOpen := u.mongodbRepo.CountTicket(ctx, fetchOptionOpen)

	// option in progress
	fetchOptionInProgress := maps.Clone(fetchOption)
	fetchOptionInProgress["status"] = []string{"in_progress"}
	totalInProgress := u.mongodbRepo.CountTicket(ctx, fetchOptionInProgress)

	// option resolve
	fetchOptionResolve := maps.Clone(fetchOption)
	fetchOptionResolve["status"] = []string{"resolve"}
	totalResolve := u.mongodbRepo.CountTicket(ctx, fetchOptionResolve)

	// option closed
	fetchOptionClosed := maps.Clone(fetchOption)
	fetchOptionClosed["status"] = []string{"closed"}
	totalClosed := u.mongodbRepo.CountTicket(ctx, fetchOptionClosed)

	// option last month
	fetchOptionLastMonth := maps.Clone(fetchOption)
	fetchOptionLastMonth["lastMonth"] = true
	ticketLastMonth := u.mongodbRepo.CountTicket(ctx, fetchOptionLastMonth)

	result := map[string]interface{}{
		"totalTicket":           total,
		"totalTicketOpen":       totalOpen,
		"totalTicketClosed":     totalClosed,
		"totalTicketInProgress": totalInProgress,
		"totalTicketResolve":    totalResolve,
		"ticketLastMonth":       ticketLastMonth,
	}

	return response.Success(result)
}

func (u *appUsecase) GetTotalTicketNow(ctx context.Context, claim domain.JWTClaimUser) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	totalOpen := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		"status": []string{"open"},
		// "companyProductID": claim.CompanyProductID,
		"startDate": startOfDay,
		"endDate":   endOfDay,
	})

	totalInProgress := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		"status": []string{"in_progress"},
		// "companyProductID": claim.CompanyProductID,
		"startDate": startOfDay,
		"endDate":   endOfDay,
	})

	totalClosed := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		"status": []string{"closed"},
		// "companyProductID": claim.CompanyProductID,
		"startDate": startOfDay,
		"endDate":   endOfDay,
	})

	result := map[string]interface{}{
		"totalTicketOpen":       totalOpen,
		"totalTicketClosed":     totalClosed,
		"totalTicketInProgress": totalInProgress,
	}

	return response.Success(result)
}

func (u *appUsecase) GetDataDashboard(ctx context.Context, claim domain.JWTClaimUser, query url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	filterOpen := map[string]interface{}{
		"status": []string{"open"},
		// "companyProductID": claim.CompanyProductID,
	}

	filterInProgress := map[string]interface{}{
		"status": []string{"in_progress"},
		// "companyProductID": claim.CompanyProductID,
	}

	filterClosed := map[string]interface{}{
		"status": []string{"closed"},
		// "companyProductID": claim.CompanyProductID,
	}

	responseDay := make([]map[string]interface{}, 0)
	for _, day := range model.Weekdays {
		u._countTicketPerDay(ctx, day, filterOpen, filterClosed, filterInProgress, &responseDay)
	}

	return response.Success(responseDay)
}

func (u *appUsecase) _countTicketPerDay(ctx context.Context, day string, optionsOpen, optionsClosed, optionsInProgress map[string]interface{}, responseDay *[]map[string]interface{}) {
	optionsOpen["day"] = day
	optionsClosed["day"] = day
	optionsInProgress["day"] = day

	ticketOpen := u.mongodbRepo.CountTicket(ctx, optionsOpen)
	ticketClosed := u.mongodbRepo.CountTicket(ctx, optionsClosed)
	ticketInProgress := u.mongodbRepo.CountTicket(ctx, optionsInProgress)

	*responseDay = append(*responseDay, map[string]interface{}{
		"dayName":    day,
		"open":       ticketOpen,
		"close":      ticketClosed,
		"inProgress": ticketInProgress,
	})
}
func (u *appUsecase) GetAverageDurationDashboard(ctx context.Context, claim domain.JWTClaimUser, query url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	fetchOptions := map[string]interface{}{
		"status": []string{"closed"},
		// "companyProductID": claim.CompanyProductID,
	}

	// Get current date
	now := time.Now()
	startOfDay := now.Truncate(24 * time.Hour)

	// Set default start and end dates
	defaultStartDate := now.AddDate(0, 0, -7).Truncate(24 * time.Hour)
	defaultEndDate := startOfDay

	var startDate, endDate time.Time
	var err error

	startDateStr := query.Get("startDate")
	endDateStr := query.Get("endDate")

	// If startDate is provided
	if startDateStr != "" {
		startDate, err = time.ParseInLocation(time.RFC3339, startDateStr, time.Local)
		if err != nil {
			return response.Error(400, "Invalid startDate format")
		}
		// If endDate is not provided, set it 7 days after the startDate
		if endDateStr == "" {
			endDate = startDate.AddDate(0, 0, 7)
		} else {
			endDate, err = time.ParseInLocation(time.RFC3339, endDateStr, time.Local)
			if err != nil {
				return response.Error(400, "Invalid endDate format")
			}
		}
		// If endDate is provided without a startDate
	} else if endDateStr != "" {
		endDate, err = time.ParseInLocation(time.RFC3339, endDateStr, time.Local)
		if err != nil {
			return response.Error(400, "Invalid endDate format")
		}
		startDate = endDate.AddDate(0, 0, -7)
	} else {
		// If neither startDate nor endDate is provided, use defaults
		startDate = defaultStartDate
		endDate = defaultEndDate
	}

	// Ensure startDate is not after endDate
	if startDate.After(endDate) {
		return response.Error(400, "startDate cannot be after endDate")
	}

	responseDate := make([]map[string]interface{}, 0)

	// Calculate the average duration
	for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
		u._calculateAverageDurationPerDate(ctx, date, fetchOptions, &responseDate)
	}

	return response.Success(responseDate)
}

func (u *appUsecase) _calculateAverageDurationPerDate(ctx context.Context, date time.Time, optionsTime map[string]interface{}, responseDate *[]map[string]interface{}) {
	optionsTime["date"] = date
	optionsTime["month"] = date.Month()
	optionsTime["year"] = date.Year()

	cur, err := u.mongodbRepo.FetchTicketList(ctx, optionsTime)
	if err != nil {
		*responseDate = append(*responseDate, map[string]interface{}{
			"date":            date,
			"averageDuration": "error",
			"error":           err.Error(),
		})
		return
	}
	defer cur.Close(ctx)

	totalDuration := 0
	count := 0

	for cur.Next(ctx) {
		var ticket map[string]interface{}
		if err := cur.Decode(&ticket); err != nil {
			logrus.Error("Error decoding ticket:", err)
			continue
		}

		// Accessing field "logTime.totalDurationInSeconds"
		if logTime, ok := ticket["logTime"].(map[string]interface{}); ok {
			if duration, ok := logTime["totalDurationInSeconds"].(int32); ok {
				totalDuration += int(duration)
				count++
			} else {
				logrus.Warnf("Field 'totalDurationInSeconds' is of unexpected type: %T, value: %v", logTime["totalDurationInSeconds"], logTime["totalDurationInSeconds"])
			}
		} else {
			logrus.Warn("Field 'logTime' is not a map or does not exist in ticket: ", ticket)
		}
	}

	if err := cur.Err(); err != nil {
		logrus.Error("Cursor error: ", err)
		*responseDate = append(*responseDate, map[string]interface{}{
			"date":            date,
			"averageDuration": "error",
			"error":           err.Error(),
		})
		return
	}

	// Calculate average duration
	averageDuration := 0
	if count > 0 {
		averageDuration = totalDuration / count
	}

	*responseDate = append(*responseDate, map[string]interface{}{
		"date":            date,
		"averageDuration": averageDuration,
	})
}
