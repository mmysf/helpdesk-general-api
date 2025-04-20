package usecase_agent

import (
	"app/domain"
	"app/domain/model"
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
)

func (u *agentUsecase) GetTotalTicket(ctx context.Context, claim domain.JWTClaimAgent) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	total := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		"companyID": claim.CompanyID})

	totalOpen := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		"status":    []string{"open"},
		"companyID": claim.CompanyID,
	})

	totalInProgress := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		"status":    []string{"in_progress"},
		"companyID": claim.CompanyID,
	})

	totalResolved := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		"status":    []string{"resolve"},
		"companyID": claim.CompanyID,
	})

	totalClosed := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		"status":    []string{"closed"},
		"companyID": claim.CompanyID,
	})

	result := map[string]interface{}{
		"totalTicket":           total,
		"totalTicketOpen":       totalOpen,
		"totalTicketInProgress": totalInProgress,
		"totalTicketResolved":   totalResolved,
		"totalTicketClosed":     totalClosed,
	}

	return response.Success(result)
}

func (u *agentUsecase) GetTotalTicketNow(ctx context.Context, claim domain.JWTClaimAgent) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	totalOpen := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		"status":    []string{"open"},
		"companyID": claim.CompanyID,
		"startDate": startOfDay,
		"endDate":   endOfDay,
	})

	totalInProgress := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		"status":    []string{"in_progress"},
		"companyID": claim.CompanyID,
		"startDate": startOfDay,
		"endDate":   endOfDay,
	})

	totalResolved := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		"status":    []string{"resolve"},
		"companyID": claim.CompanyID,
		"startDate": startOfDay,
		"endDate":   endOfDay,
	})

	totalClosed := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		"status":    []string{"closed"},
		"companyID": claim.CompanyID,
		"startDate": startOfDay,
		"endDate":   endOfDay,
	})

	result := map[string]interface{}{
		"totalTicketOpen":       totalOpen,
		"totalTicketInProgress": totalInProgress,
		"totalTicketResolved":   totalResolved,
		"totalTicketClosed":     totalClosed,
	}

	return response.Success(result)
}

func (u *agentUsecase) GetDataDashboard(ctx context.Context, claim domain.JWTClaimAgent, query url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	filterOpen := map[string]interface{}{
		"status":    []string{"open", "processing"},
		"companyID": claim.CompanyID,
	}
	filterInProgress := map[string]interface{}{
		"status":    []string{"in_progress"},
		"companyID": claim.CompanyID,
	}
	filterResolved := map[string]interface{}{
		"status":    []string{"resolve"},
		"companyID": claim.CompanyID,
	}
	filterClosed := map[string]interface{}{
		"status":    []string{"closed"},
		"companyID": claim.CompanyID,
	}

	if query.Get("companyProductName") != "" {
		filterOpen["companyProductName"] = query.Get("companyProductName")
		filterInProgress["companyProductName"] = query.Get("companyProductName")
		filterResolved["companyProductName"] = query.Get("companyProductName")
		filterClosed["companyProductName"] = query.Get("companyProductName")
	}

	if query.Get("companyProductID") != "" {
		filterOpen["companyProductID"] = query.Get("companyProductID")
		filterInProgress["companyProductID"] = query.Get("companyProductID")
		filterResolved["companyProductID"] = query.Get("companyProductID")
		filterClosed["companyProductID"] = query.Get("companyProductID")
	}

	if query.Get("projectID") != "" {
		filterOpen["projectID"] = query.Get("projectID")
		filterInProgress["projectID"] = query.Get("projectID")
		filterResolved["projectID"] = query.Get("projectID")
		filterClosed["projectID"] = query.Get("projectID")
	}

	if query.Get("customerID") != "" {
		filterOpen["customerID"] = query.Get("customerID")
		filterInProgress["customerID"] = query.Get("customerID")
		filterResolved["customerID"] = query.Get("customerID")
		filterClosed["customerID"] = query.Get("customerID")
	}

	if query.Get("agentID") != "" {
		filterOpen["agentID"] = query.Get("agentID")
		filterInProgress["agentID"] = query.Get("agentID")
		filterResolved["agentID"] = query.Get("agentID")
		filterClosed["agentID"] = query.Get("agentID")
	}

	if query.Get("status") != "" {
		filterOpen["status"] = strings.Split(query.Get("status"), ",")
		filterInProgress["status"] = strings.Split(query.Get("status"), ",")
		filterResolved["status"] = strings.Split(query.Get("status"), ",")
		filterClosed["status"] = strings.Split(query.Get("status"), ",")
	}

	if query.Get("code") != "" {
		filterOpen["code"] = query.Get("code")
		filterInProgress["code"] = query.Get("code")
		filterResolved["code"] = query.Get("code")
		filterClosed["code"] = query.Get("code")
	}

	if query.Get("subject") != "" {
		filterOpen["subject"] = query.Get("subject")
		filterInProgress["subject"] = query.Get("subject")
		filterResolved["subject"] = query.Get("subject")
		filterClosed["subject"] = query.Get("subject")
	}

	responseDay := make([]map[string]interface{}, 0)
	for _, day := range model.Weekdays {
		u._countTicketPerDay(ctx, day, filterOpen, filterInProgress, filterResolved, filterClosed, &responseDay)
	}

	return response.Success(responseDay)
}

func (u *agentUsecase) _countTicketPerDay(ctx context.Context, day string, optionsOpen, optionsInProgress, optionsResolved, optionsClosed map[string]interface{}, responseDay *[]map[string]interface{}) {
	optionsOpen["day"] = day
	optionsInProgress["day"] = day
	optionsResolved["day"] = day
	optionsClosed["day"] = day

	ticketOpen := u.mongodbRepo.CountTicket(ctx, optionsOpen)
	ticketInProgress := u.mongodbRepo.CountTicket(ctx, optionsInProgress)
	ticketResolved := u.mongodbRepo.CountTicket(ctx, optionsResolved)
	ticketClosed := u.mongodbRepo.CountTicket(ctx, optionsClosed)

	*responseDay = append(*responseDay, map[string]interface{}{
		"dayName":    day,
		"open":       ticketOpen,
		"inProgress": ticketInProgress,
		"resolve":    ticketResolved,
		"close":      ticketClosed,
	})
}
