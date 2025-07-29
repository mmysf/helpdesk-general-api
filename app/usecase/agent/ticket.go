package usecase_agent

import (
	"app/domain"
	"app/domain/model"
	"app/helpers"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	yurekahelpers "github.com/Yureka-Teknologi-Cipta/yureka/helpers"
	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *agentUsecase) GetTicketList(ctx context.Context, claim domain.JWTClaimAgent, query url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, limit, offset := yurekahelpers.GetLimitOffset(query)

	//get agent
	agent, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{
		"id": claim.UserID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	fetchOptions := map[string]interface{}{
		"limit":     limit,
		"offset":    offset,
		"companyID": claim.CompanyID,
	}

	// filtering
	if agent.Role == "agent" {
		fetchOptions["categoryID"] = agent.Category.ID
	}

	if query.Get("sort") != "" {
		fetchOptions["sort"] = query.Get("sort")
	}

	if query.Get("dir") != "" {
		fetchOptions["dir"] = query.Get("dir")
	}

	if query.Get("companyProductName") != "" {
		fetchOptions["companyProductName"] = query.Get("companyProductName")
	}

	if query.Get("companyProductID") != "" {
		fetchOptions["companyProductID"] = query.Get("companyProductID")
	}

	if query.Get("projectID") != "" {
		fetchOptions["projectID"] = query.Get("projectID")
	}

	if query.Get("customerID") != "" {
		fetchOptions["customerID"] = query.Get("customerID")
	}

	if query.Get("agentID") != "" {
		fetchOptions["agentID"] = query.Get("agentID")
	}

	if query.Get("status") != "" {
		fetchOptions["status"] = strings.Split(query.Get("status"), ",")
	}

	if query.Get("code") != "" {
		fetchOptions["code"] = query.Get("code")
	}

	if query.Get("subject") != "" {
		fetchOptions["subject"] = query.Get("subject")
	}

	if query.Get("priority") != "" {
		fetchOptions["priority"] = query.Get("priority")
	}

	if query.Get("categoryID") != "" {
		fetchOptions["categoryID"] = query.Get("categoryID")
	}

	if query.Get("completedBy") != "" {
		fetchOptions["completedBy"] = query.Get("completedBy")
	}

	// count first
	totalDocuments := u.mongodbRepo.CountTicket(ctx, fetchOptions)

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
	cur, err := u.mongodbRepo.FetchTicketList(ctx, fetchOptions)

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
		row := model.Ticket{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("Ticket Decode ", err)
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

		row.Format(claim.UserID)

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

func (u *agentUsecase) GetMyTicketList(ctx context.Context, claim domain.JWTClaimAgent, query url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, limit, offset := yurekahelpers.GetLimitOffset(query)

	fetchOptions := map[string]interface{}{
		"limit":     limit,
		"offset":    offset,
		"companyID": claim.CompanyID,
		"agentID":   claim.UserID,
	}

	// filtering
	if query.Get("sort") != "" {
		fetchOptions["sort"] = query.Get("sort")
	}

	if query.Get("dir") != "" {
		fetchOptions["dir"] = query.Get("dir")
	}

	if query.Get("companyProductName") != "" {
		fetchOptions["companyProductName"] = query.Get("companyProductName")
	}

	if query.Get("companyProductID") != "" {
		fetchOptions["companyProductID"] = query.Get("companyProductID")
	}

	if query.Get("projectID") != "" {
		fetchOptions["projectID"] = query.Get("projectID")
	}

	if query.Get("customerID") != "" {
		fetchOptions["customerID"] = query.Get("customerID")
	}

	if query.Get("status") != "" {
		fetchOptions["status"] = strings.Split(query.Get("status"), ",")
	}

	if query.Get("code") != "" {
		fetchOptions["code"] = query.Get("code")
	}

	if query.Get("subject") != "" {
		fetchOptions["subject"] = query.Get("subject")
	}

	if query.Get("priority") != "" {
		fetchOptions["priority"] = query.Get("priority")
	}

	if query.Get("categoryID") != "" {
		fetchOptions["categoryID"] = query.Get("categoryID")
	}

	// count first
	totalDocuments := u.mongodbRepo.CountTicket(ctx, fetchOptions)

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
	cur, err := u.mongodbRepo.FetchTicketList(ctx, fetchOptions)

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
		row := model.Ticket{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("Ticket Decode ", err)
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

		row.Format(claim.UserID)

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

func (u *agentUsecase) GetTicketDetail(ctx context.Context, claim domain.JWTClaimAgent, ticketId string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check ticket
	ticket, err := u.mongodbRepo.FetchOneTicket(ctx, map[string]interface{}{
		"id":         ticketId,
		"commpanyID": claim.CompanyID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if ticket == nil {
		return response.Error(http.StatusBadRequest, "ticket not found")
	}

	return response.Success(ticket.Format(claim.UserID))
}

func (u *agentUsecase) CloseTicket(ctx context.Context, claim domain.JWTClaimAgent, payload domain.CloseTicketRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	now := time.Now()

	errValidation := make(map[string]string)
	// validating
	if payload.TicketId == "" {
		errValidation["id"] = "id field is required"
	}

	// check ticket
	ticket, err := u.mongodbRepo.FetchOneTicket(ctx, map[string]interface{}{
		"id":        payload.TicketId,
		"companyID": claim.CompanyID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if ticket == nil {
		return response.Error(http.StatusBadRequest, "ticket not found")
	}

	// check log status
	if ticket.LogTime.Status == model.Running {
		return response.Error(http.StatusBadRequest, "log still running")
	}

	// update ticket
	ticket.Status = model.Closed
	ticket.ClosedAt = &now
	ticket.UpdatedAt = now

	if err := u.mongodbRepo.UpdateTicket(ctx, ticket); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(ticket)
}

func (u *agentUsecase) ReopenTicket(ctx context.Context, claim domain.JWTClaimAgent, payload domain.ReopenTicketRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating
	if payload.TicketId == "" {
		errValidation["id"] = "id field is required"
	}

	// check ticket
	ticket, err := u.mongodbRepo.FetchOneTicket(ctx, map[string]interface{}{
		"id":        payload.TicketId,
		"companyID": claim.CompanyID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if ticket == nil {
		return response.Error(http.StatusBadRequest, "ticket not found")
	}

	// update ticket
	ticket.Status = model.Open
	ticket.UpdatedAt = time.Now()

	if err := u.mongodbRepo.UpdateTicket(ctx, ticket); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(ticket)
}

func (u *agentUsecase) StartLoggingTicket(ctx context.Context, claim domain.JWTClaimAgent, payload domain.LoggingTicketRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating
	if payload.TicketId == "" {
		errValidation["id"] = "id field is required"
	}

	// check ticket
	ticket, err := u.mongodbRepo.FetchOneTicket(ctx, map[string]interface{}{
		"id":        payload.TicketId,
		"companyID": claim.CompanyID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if ticket == nil {
		return response.Error(http.StatusBadRequest, "ticket not found")
	}

	// check log status
	if ticket.LogTime.Status == model.Running {
		return response.Error(http.StatusBadRequest, "log already running")
	}

	// update ticket
	startAt := time.Now()
	ticket.LogTime.StartAt = &startAt
	ticket.LogTime.DurationInSeconds = 0
	ticket.ReminderSent = false
	ticket.LogTime.PauseDurationInSeconds = 0
	ticket.LogTime.PauseHistory = []model.PauseHistory{}
	ticket.LogTime.Status = model.Running
	ticket.LogTime.EndAt = nil
	ticket.Status = model.InProgress
	ticket.UpdatedAt = time.Now()

	if err := u.mongodbRepo.UpdateTicket(ctx, ticket); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	u.mongodbRepo.CreateTicketTimelogs(ctx, &model.TicketTimeLogs{
		ID:       primitive.NewObjectID(),
		Company:  ticket.Company,
		Customer: ticket.Customer,
		// Product:  ticket.Product,
		Ticket: model.TicketNested{
			ID:       ticket.ID.Hex(),
			Subject:  ticket.Subject,
			Content:  ticket.Content,
			Priority: ticket.Priority,
		},
		DurationInSeconds:      ticket.LogTime.DurationInSeconds,
		PauseDurationInSeconds: ticket.LogTime.PauseDurationInSeconds,
		StartAt:                ticket.LogTime.StartAt,
		EndAt:                  ticket.LogTime.EndAt,
		PauseHistory:           ticket.LogTime.PauseHistory,
		IsManual:               false,
		ActivityTtype:          "start log",
		CreatedAt:              time.Now(),
	})

	// find company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": ticket.Company.ID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if company == nil {
		return response.Error(http.StatusBadRequest, "company not found")
	}

	go _sendActivityNotification(ticket, company)

	return response.Success(ticket)
}

func (u *agentUsecase) StopLoggingTicket(ctx context.Context, claim domain.JWTClaimAgent, payload domain.LoggingTicketRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating
	if payload.TicketId == "" {
		errValidation["id"] = "id field is required"
	}

	// check ticket
	ticket, err := u.mongodbRepo.FetchOneTicket(ctx, map[string]interface{}{
		"id":        payload.TicketId,
		"companyID": claim.CompanyID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if ticket == nil {
		return response.Error(http.StatusBadRequest, "ticket not found")
	}

	// check log status
	if ticket.LogTime.Status != model.Running {
		return response.Error(http.StatusBadRequest, "log not running")
	}

	// update ticket
	endAt := time.Now()
	ticket.LogTime.EndAt = &endAt
	ticket.LogTime.DurationInSeconds = int(endAt.Sub(*ticket.LogTime.StartAt).Seconds()) - ticket.LogTime.PauseDurationInSeconds
	ticket.LogTime.TotalDurationInSeconds += ticket.LogTime.DurationInSeconds
	ticket.LogTime.TotalPausedDurationInSeconds += ticket.LogTime.PauseDurationInSeconds
	ticket.LogTime.Status = model.Done
	ticket.UpdatedAt = time.Now()

	if err := u.mongodbRepo.UpdateTicket(ctx, ticket); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// create ticket timelogs
	u.mongodbRepo.CreateTicketTimelogs(ctx, &model.TicketTimeLogs{
		ID:       primitive.NewObjectID(),
		Company:  ticket.Company,
		Customer: ticket.Customer,
		// Product:  ticket.Product,
		Ticket: model.TicketNested{
			ID:       ticket.ID.Hex(),
			Subject:  ticket.Subject,
			Content:  ticket.Content,
			Priority: ticket.Priority,
		},
		DurationInSeconds:      ticket.LogTime.DurationInSeconds,
		PauseDurationInSeconds: ticket.LogTime.PauseDurationInSeconds,
		StartAt:                ticket.LogTime.StartAt,
		EndAt:                  ticket.LogTime.EndAt,
		PauseHistory:           ticket.LogTime.PauseHistory,
		IsManual:               false,
		ActivityTtype:          "stop log",
		CreatedAt:              time.Now(),
	})

	// find company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": ticket.Company.ID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if company == nil {
		return response.Error(http.StatusBadRequest, "company not found")
	}

	go _sendActivityNotification(ticket, company)

	return response.Success(ticket)
}

func (u *agentUsecase) PauseLoggingTicket(ctx context.Context, claim domain.JWTClaimAgent, payload domain.LoggingTicketRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating
	if payload.TicketId == "" {
		errValidation["id"] = "id field is required"
	}

	// check ticket
	ticket, err := u.mongodbRepo.FetchOneTicket(ctx, map[string]interface{}{
		"id":        payload.TicketId,
		"companyID": claim.CompanyID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if ticket == nil {
		return response.Error(http.StatusBadRequest, "ticket not found")
	}

	// check log status
	if ticket.LogTime.Status == model.Paused {
		return response.Error(http.StatusBadRequest, "log already paused")
	}

	if ticket.LogTime.Status != model.Running {
		return response.Error(http.StatusBadRequest, "log not running")
	}

	// update ticket
	duration := 0
	now := time.Now()
	if len(ticket.LogTime.PauseHistory) == 0 {
		duration = int(now.Sub(*ticket.LogTime.StartAt).Seconds())
		ticket.LogTime.DurationInSeconds = duration
	} else {
		lastPause := ticket.LogTime.PauseHistory[len(ticket.LogTime.PauseHistory)-1]
		if lastPause.ResumedAt != nil {
			duration = int(now.Sub(*lastPause.ResumedAt).Seconds())
			ticket.LogTime.DurationInSeconds += duration
		}
	}
	ticket.LogTime.PauseHistory = append(ticket.LogTime.PauseHistory, model.PauseHistory{
		PausedAt: now,
	})
	ticket.LogTime.Status = model.Paused
	ticket.UpdatedAt = now

	if err := u.mongodbRepo.UpdateTicket(ctx, ticket); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// check ticket timelogs
	timelogs, e := u.mongodbRepo.FetchOneTicketlogs(ctx, map[string]interface{}{
		"ticket.id": payload.TicketId,
		"sort":      "createdAt",
		"dir":       "desc",
	})

	if e != nil {
		return response.Error(http.StatusInternalServerError, e.Error())
	}

	if timelogs == nil {
		return response.Error(http.StatusBadRequest, "time log not found")
	}

	timelogs.EndAt = &now
	timelogs.DurationInSeconds = duration
	timelogs.UpdatedAt = &now

	if err := u.mongodbRepo.UpdateTicketlogs(ctx, timelogs); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// update time balance
	if err := u._updateTimeBalance(ctx, ticket, duration); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// find company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": ticket.Company.ID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if company == nil {
		return response.Error(http.StatusBadRequest, "company not found")
	}

	// notification
	go _sendActivityNotification(ticket, company)

	return response.Success(ticket)
}

func (u *agentUsecase) ResumeLoggingTicket(ctx context.Context, claim domain.JWTClaimAgent, payload domain.LoggingTicketRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating
	if payload.TicketId == "" {
		errValidation["id"] = "id field is required"
	}

	// check ticket
	ticket, err := u.mongodbRepo.FetchOneTicket(ctx, map[string]interface{}{
		"id":        payload.TicketId,
		"companyID": claim.CompanyID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if ticket == nil {
		return response.Error(http.StatusBadRequest, "ticket not found")
	}

	// check log status
	if ticket.LogTime.Status != model.Paused {
		return response.Error(http.StatusBadRequest, "log is not paused")
	}

	now := time.Now()

	// check customer
	customer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"id": ticket.Customer.ID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if customer == nil {
		return response.Error(http.StatusBadRequest, "customer not found")
	}

	if customer.IsNeedBalance {
		assigned := false
		for _, agent := range ticket.Agent {
			if agent.ID == claim.UserID {
				assigned = true
				break
			}
		}
		if !assigned {
			return response.Error(http.StatusBadRequest, "you are not assigned to this ticket")
		}

		if len(ticket.Agent) == 0 {
			return response.Error(http.StatusBadRequest, "you are not assigned to this ticket")
		}
	}

	//check pause history
	if len(ticket.LogTime.PauseHistory) == 0 {
		return response.Error(http.StatusBadRequest, "no pause history found")
	}

	lastPause := &ticket.LogTime.PauseHistory[len(ticket.LogTime.PauseHistory)-1]
	if lastPause.ResumedAt != nil {
		return response.Error(http.StatusBadRequest, "log already resumed")
	}

	// update ticket
	lastPause.ResumedAt = &now
	ticket.LogTime.PauseDurationInSeconds += int(now.Sub(lastPause.PausedAt).Seconds())
	ticket.LogTime.Status = model.Running
	ticket.UpdatedAt = now

	if ticket.LogTime.PauseDurationInSeconds < 1 {
		return response.Error(http.StatusBadRequest, "pause duration is invalid")
	}

	if err := u.mongodbRepo.UpdateTicket(ctx, ticket); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// create ticket timelogs
	u.mongodbRepo.CreateTicketTimelogs(ctx, &model.TicketTimeLogs{
		ID:       primitive.NewObjectID(),
		Company:  ticket.Company,
		Customer: ticket.Customer,
		// Product:  ticket.Product,
		Ticket: model.TicketNested{
			ID:       ticket.ID.Hex(),
			Subject:  ticket.Subject,
			Content:  ticket.Content,
			Priority: ticket.Priority,
		},
		DurationInSeconds:      0,
		PauseDurationInSeconds: 0,
		StartAt:                &now,
		EndAt:                  nil,
		IsManual:               false,
		ActivityTtype:          "resume log",
		CreatedAt:              time.Now(),
	})

	// find company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": ticket.Company.ID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if company == nil {
		return response.Error(http.StatusBadRequest, "company not found")
	}

	// notification
	go _sendActivityNotification(ticket, company)

	return response.Success(ticket)
}

func (u *agentUsecase) CreateTicketComment(ctx context.Context, claim domain.JWTClaimAgent, payload domain.TicketCommentRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	UserID := claim.UserID

	now := time.Now()

	errValidation := make(map[string]string)
	// validating

	if payload.TicketId == "" {
		errValidation["ticketId"] = "ticketId field is required"
	}

	if payload.Content == "" {
		errValidation["content"] = "content field is required"
	}

	if payload.Status == "" {
		errValidation["status"] = "status field is required"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check ticket
	ticket, err := u.mongodbRepo.FetchOneTicket(ctx, map[string]interface{}{
		"id":        payload.TicketId,
		"companyID": claim.CompanyID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if ticket == nil {
		return response.Error(http.StatusBadRequest, "ticked not found")
	}

	// check company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": ticket.Company.ID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if company == nil {
		return response.Error(http.StatusBadRequest, "company not found")
	}

	// check customer
	customer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"id": ticket.Customer.ID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if customer == nil {
		return response.Error(http.StatusBadRequest, "customer not found")
	}

	if customer.IsNeedBalance {
		assigned := false
		for _, agent := range ticket.Agent {
			if agent.ID == UserID {
				assigned = true
				break
			}
		}
		if !assigned {
			return response.Error(http.StatusBadRequest, "you are not assigned to this ticket")
		}

		if len(ticket.Agent) == 0 {
			return response.Error(http.StatusBadRequest, "you are not assigned to this ticket")
		}
	}

	// validate status
	switch ticket.Status {
	case model.Open:
		if !helpers.InArrayString(string(payload.Status), []string{string(model.Open), string(model.InProgress)}) {
			return response.Error(400, "status allowed only open or in progress")
		}
	case model.InProgress:
		if !helpers.InArrayString(string(payload.Status), []string{string(model.InProgress), string(model.Resolve)}) {
			return response.Error(400, "status allowed only in progress or resolve")
		}
	case model.Resolve:
		if !helpers.InArrayString(string(payload.Status), []string{string(model.InProgress), string(model.Resolve)}) {
			return response.Error(400, "status allowed only in progress or resolve")
		}
	}

	// check customer balance
	if customer.IsNeedBalance {
		if payload.Status == model.InProgress {
			if !(now.After(customer.Subscription.StartAt) && now.Before(customer.Subscription.EndAt)) || customer.Subscription.Status != model.Active {
				return response.Error(http.StatusBadRequest, "the customer doesn't have active subscription")
			}
			timeRemaining := customer.Subscription.Balance.Time.Total - customer.Subscription.Balance.Time.Used
			if customer.Subscription.Balance == nil || timeRemaining < 0 {
				return response.Error(http.StatusBadRequest, "the customer doesn't have enough time balance")
			}
		}
	}

	// get detail attachments
	ticketCommentAttachments := make([]model.AttachmentFK, 0)
	if len(payload.AttachIds) > 0 {
		cur, err := u.mongodbRepo.FetchAttachmentList(ctx, map[string]interface{}{
			"ids":        payload.AttachIds,
			"company_id": claim.CompanyID,
		})
		if err != nil {
			return response.Error(http.StatusBadRequest, "attachments not found")
		}

		defer cur.Close(ctx)

		listAttachmentOri := make([]model.Attachment, 0)
		for cur.Next(ctx) {
			row := model.Attachment{}
			err := cur.Decode(&row)
			if err != nil {
				logrus.Error("Attachment Decode ", err)
				return response.Error(http.StatusInternalServerError, err.Error())
			}

			listAttachmentOri = append(listAttachmentOri, row)
		}

		for _, attachment := range listAttachmentOri {
			attach := model.AttachmentFK{
				ID:          attachment.ID.Hex(),
				Name:        attachment.Name,
				Size:        attachment.Size,
				URL:         attachment.URL,
				Type:        attachment.Type,
				ProviderKey: attachment.ProviderKey,
				IsPrivate:   attachment.IsPrivate,
			}

			ticketCommentAttachments = append(ticketCommentAttachments, attach)
		}
	}

	// update attachment isUsed
	attachIds := make([]primitive.ObjectID, 0)
	for _, attachID := range payload.AttachIds {
		primitiveID, _ := primitive.ObjectIDFromHex(attachID)
		attachIds = append(attachIds, primitive.ObjectID(primitiveID))
	}

	if err := u.mongodbRepo.UpdateManyAttachmentPartial(ctx, attachIds, map[string]interface{}{
		"isUsed": true,
	}); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// create ticket comment
	ticketComment := &model.TicketComment{
		ID:      primitive.NewObjectID(),
		Company: claim.Company,
		// Product: model.CompanyProductNested{
		// 	ID:    ticket.Product.ID,
		// 	Name:  ticket.Product.Name,
		// 	Image: ticket.Product.Image,
		// 	Code:  ticket.Product.Code,
		// },
		Agent: model.AgentNested{
			ID:   claim.User.ID,
			Name: claim.User.Name,
		},
		Ticket: model.TicketNested{
			ID:      ticket.ID.Hex(),
			Subject: ticket.Subject,
		},
		Content:     payload.Content,
		Sender:      model.AgentSender,
		Attachments: ticketCommentAttachments,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := u.mongodbRepo.CreateTicketComment(ctx, ticketComment); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// update ticket & ticket LogTime
	if err := u._updateTicketAndTimelog(ctx, ticket, payload.Status, claim.User, company); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// get from config
	config := u._CacheConfig(ctx)

	// notification
	go _sendCommentNotfication(config, ticket, ticketComment, &claim.User, company)

	return response.Success(ticketComment)
}

func _sendCommentNotfication(config model.Config, ticket *model.Ticket, ticketComment *model.TicketComment, agent *model.UserNested, company *model.Company) {
	//mail content
	mailer := helpers.NewSMTPMailer(company)
	mailer.To([]string{ticket.Customer.Email})
	mailer.Subject(config.Email.Template.TicketComment.Title)
	mailer.Body(helpers.StringReplacer(config.Email.Template.TicketComment.Body, map[string]string{
		"title":          config.Email.Template.TicketComment.Title,
		"customer_name":  ticket.Customer.Name,
		"ticket_subject": ticket.Subject,
		"agent_name":     agent.Name,
		"comment":        ticketComment.Content,
	}))

	//send mail
	if err := mailer.Send(); err != nil {
		logrus.Errorf("Send Email to %s error %v", ticket.Customer.Email, err)
	}
}

func _sendActivityNotification(ticket *model.Ticket, company *model.Company) {
	//get email receiver
	receiverEmail := ticket.Customer.Email

	// assign helper mailer
	mailer := helpers.NewSMTPMailer(company)

	// setup mail content
	subject := fmt.Sprintf("New activity on your ticket : %s", ticket.Subject)
	body := fmt.Sprintf(`
		<p>Hello,</p>
		<p>Your ticket has been updated its status to: <strong>%s</strong></p>
	`, ticket.LogTime.Status)

	//assign mail content
	mailer.To([]string{receiverEmail})
	mailer.Subject(subject)
	mailer.Body(body)

	//send mail
	if err := mailer.Send(); err != nil {
		logrus.WithFields(logrus.Fields{
			"ticketID": ticket.ID.Hex(),
			"subject":  ticket.Subject,
			"receiver": receiverEmail,
		}).Errorf("Failed to send email: %s", err.Error())
	}
}

func (u *agentUsecase) GetTicketCommentList(ctx context.Context, claim domain.JWTClaimAgent, ticketId string, query url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, limit, offset := yurekahelpers.GetLimitOffset(query)

	fetchOptions := map[string]interface{}{
		"limit":     limit,
		"offset":    offset,
		"ticketID":  ticketId,
		"companyID": claim.CompanyID,
	}

	// filtering
	if query.Get("sort") != "" {
		fetchOptions["sort"] = query.Get("sort")
	}

	if query.Get("dir") != "" {
		fetchOptions["dir"] = query.Get("dir")
	}

	if query.Get("companyProductID") != "" {
		fetchOptions["companyProductID"] = query.Get("companyProductID")
	}

	if query.Get("projectID") != "" {
		fetchOptions["projectID"] = query.Get("projectID")
	}

	if query.Get("q") != "" {
		fetchOptions["q"] = query.Get("q")
	}

	// count first
	totalDocuments := u.mongodbRepo.CountTicketComment(ctx, fetchOptions)

	if totalDocuments == 0 {
		return response.Success(response.List{
			List:  []interface{}{},
			Page:  page,
			Limit: limit,
			Total: totalDocuments,
		})
	}

	// check ticket list
	cur, err := u.mongodbRepo.FetchTicketCommentList(ctx, fetchOptions)

	if err != nil {
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

	defer cur.Close(ctx)

	list := make([]interface{}, 0)
	for cur.Next(ctx) {
		row := model.TicketComment{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("Tickect comment decode", err)
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

		list = append(list, row)
	}

	return response.Success(
		domain.ResponseList{
			List: response.List{
				List:  list,
				Page:  page,
				Limit: limit,
				Total: totalDocuments,
			},
			TotalPage: helpers.GetTotalPage(totalDocuments, limit),
		},
	)
}

func (u *agentUsecase) GetTicketCommentDetail(ctx context.Context, claim domain.JWTClaimAgent, commentId string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check ticket
	ticket, err := u.mongodbRepo.FetchOneTicketComment(ctx, map[string]interface{}{
		"id":        commentId,
		"companyID": claim.CompanyID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if ticket == nil {
		return response.Error(http.StatusBadRequest, "ticket not found")
	}

	return response.Success(ticket)
}

func (u *agentUsecase) EditTimeTrack(ctx context.Context, claim domain.JWTClaimAgent, ticketId string, payload domain.TimeTrackRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	now := time.Now()
	second := int(payload.Hour)*3600 + payload.Minute*60 + payload.Second

	// check ticket
	ticket, err := u.mongodbRepo.FetchOneTicket(ctx, map[string]interface{}{
		"id": ticketId,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if ticket == nil {
		return response.Error(http.StatusBadRequest, "ticked not found")
	}

	ticket.LogTime.StartAt = nil
	ticket.LogTime.EndAt = nil
	ticket.LogTime.DurationInSeconds = 0
	ticket.LogTime.TotalDurationInSeconds = second
	ticket.UpdatedAt = now

	if err := u.mongodbRepo.UpdateTicket(ctx, ticket); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// create ticket timelogs
	endAt := now.Add(time.Second * time.Duration(second))
	u.mongodbRepo.CreateTicketTimelogs(ctx, &model.TicketTimeLogs{
		ID:       primitive.NewObjectID(),
		Company:  ticket.Company,
		Customer: ticket.Customer,
		// Product:  ticket.Product,
		Ticket: model.TicketNested{
			ID:       ticket.ID.Hex(),
			Subject:  ticket.Subject,
			Content:  ticket.Content,
			Priority: ticket.Priority,
		},
		DurationInSeconds: second,
		StartAt:           &now,
		EndAt:             &endAt,
		IsManual:          true,
		ActivityTtype:     "edit time track",
		CreatedAt:         now,
	})

	return response.Success(ticket)
}

func (u *agentUsecase) _updateTicketAndTimelog(ctx context.Context, ticket *model.Ticket, status model.TicketStatus, agent model.UserNested, company *model.Company) (err error) {
	now := time.Now()

	u.mongodbRepo.UpdateOneAgent(
		ctx,
		map[string]interface{}{"id": agent.ID},
		map[string]interface{}{
			"updatedAt":      now,
			"lastActivityAt": now,
		})

	if ticket.Status != model.Open {
		if status == model.Open {
			ticket.Status = model.Open
			ticket.UpdatedAt = time.Now()

			if err := u.mongodbRepo.UpdateTicket(ctx, ticket); err != nil {
				return err
			}

			// create ticket timelogs
			u.mongodbRepo.CreateTicketTimelogs(ctx, &model.TicketTimeLogs{
				ID:       primitive.NewObjectID(),
				Company:  ticket.Company,
				Customer: ticket.Customer,
				// Product:  ticket.Product,
				Ticket: model.TicketNested{
					ID:       ticket.ID.Hex(),
					Subject:  ticket.Subject,
					Content:  ticket.Content,
					Priority: ticket.Priority,
				},
				DurationInSeconds:      ticket.LogTime.DurationInSeconds,
				PauseDurationInSeconds: ticket.LogTime.PauseDurationInSeconds,
				StartAt:                ticket.LogTime.StartAt,
				EndAt:                  ticket.LogTime.EndAt,
				PauseHistory:           ticket.LogTime.PauseHistory,
				IsManual:               false,
				ActivityTtype:          "ticket open",
				CreatedAt:              now,
			})
		}
	}

	if ticket.Status != model.InProgress {
		if status == model.InProgress {
			ticket.LogTime.StartAt = &now
			ticket.LogTime.EndAt = nil
			ticket.LogTime.Status = model.Running
			ticket.LogTime.DurationInSeconds = 0
			ticket.LogTime.PauseDurationInSeconds = 0
			ticket.LogTime.PauseHistory = []model.PauseHistory{}
			ticket.ReminderSent = false
			ticket.Status = model.InProgress
			ticket.UpdatedAt = time.Now()

			if err := u.mongodbRepo.UpdateTicket(ctx, ticket); err != nil {
				return err
			}

			// create ticket timelogs
			u.mongodbRepo.CreateTicketTimelogs(ctx, &model.TicketTimeLogs{
				ID:       primitive.NewObjectID(),
				Company:  ticket.Company,
				Customer: ticket.Customer,
				// Product:  ticket.Product,
				Ticket: model.TicketNested{
					ID:       ticket.ID.Hex(),
					Subject:  ticket.Subject,
					Content:  ticket.Content,
					Priority: ticket.Priority,
				},
				DurationInSeconds:      ticket.LogTime.DurationInSeconds,
				PauseDurationInSeconds: ticket.LogTime.PauseDurationInSeconds,
				StartAt:                ticket.LogTime.StartAt,
				EndAt:                  ticket.LogTime.EndAt,
				PauseHistory:           ticket.LogTime.PauseHistory,
				IsManual:               false,
				ActivityTtype:          "ticket in progress",
				CreatedAt:              now,
			})
		}
	}

	if ticket.Status != model.Resolve {
		if status == model.Resolve {
			//default token
			defaultToken := helpers.RandomString(64)

			ticket.LogTime.EndAt = &now

			//get agent
			ticketAgent, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{"id": agent.ID})
			if err != nil {
				return err
			}

			//init logs end time
			logsEndAt := now

			//check pause
			duration := 0
			if len(ticket.LogTime.PauseHistory) == 0 {
				duration = int(now.Sub(*ticket.LogTime.StartAt).Seconds())
				ticket.LogTime.DurationInSeconds = duration
			} else { //count active duration after last pause
				lastPause := &ticket.LogTime.PauseHistory[len(ticket.LogTime.PauseHistory)-1]
				if lastPause.ResumedAt != nil {
					duration = int(now.Sub(*lastPause.ResumedAt).Seconds())
					ticket.LogTime.DurationInSeconds += duration
				} else {
					lastPause.ResumedAt = &now
					logsEndAt = lastPause.PausedAt
				}
			}
			ticket.LogTime.TotalDurationInSeconds += ticket.LogTime.DurationInSeconds
			ticket.LogTime.TotalPausedDurationInSeconds += ticket.LogTime.PauseDurationInSeconds
			ticket.LogTime.Status = model.Done
			ticket.ReminderSent = true
			ticket.Token = defaultToken
			ticket.Status = model.Resolve
			ticket.CompletedBy = &model.AgentNested{
				ID:    ticketAgent.ID.Hex(),
				Name:  ticketAgent.Name,
				Email: ticketAgent.Email,
			}
			ticket.UpdatedAt = now

			ticketAgent.TotalTicketCompleted++
			ticketAgent.UpdatedAt = now

			if err := u.mongodbRepo.UpdateTicket(ctx, ticket); err != nil {
				return err
			}

			if err := u.mongodbRepo.UpdateAgent(ctx, ticketAgent); err != nil {
				return err
			}

			if err := u._updateTimeBalance(ctx, ticket, duration); err != nil {
				return err
			}

			// get from config
			config := u._CacheConfig(ctx)

			go _sendConfirmCloseTicketNotification(config, ticket, company)

			// check ticket timelogs
			timelogs, err := u.mongodbRepo.FetchOneTicketlogs(ctx, map[string]interface{}{
				"ticket.id": ticket.ID.Hex(),
				"sort":      "createdAt",
				"dir":       "desc",
			})
			if err != nil {
				return err
			}

			if timelogs == nil {
				return fmt.Errorf("time log not found")
			}

			timelogs.EndAt = &logsEndAt
			timelogs.DurationInSeconds += duration
			timelogs.UpdatedAt = &now

			if err := u.mongodbRepo.UpdateTicketlogs(ctx, timelogs); err != nil {
				return err
			}
		}
	}

	// create notification
	if err := u._createNotification(ctx, ticket, &ticket.Company); err != nil {
		return err
	}

	return nil
}

func (u *agentUsecase) _createNotification(ctx context.Context, ticket *model.Ticket, company *model.CompanyNested) (err error) {
	// notif
	var title string
	var content string

	if ticket.Status == model.Open {
		title = "Status ticket change"
		content = "Ticket opened"
	}

	if ticket.Status == model.InProgress {
		title = "Status ticket change"
		content = "Ticket in progress"
	}

	if ticket.Status == model.Resolve {
		title = "Status ticket change"
		content = "Ticket resolved"
	}

	// create notification
	notification := &model.Notification{
		ID:       primitive.NewObjectID(),
		Company:  model.CompanyNested{ID: company.ID, Name: company.Name},
		Title:    title,
		Content:  content,
		IsRead:   false,
		UserRole: model.AgentRole,
		User:     model.UserNested(ticket.Customer),
		Type:     model.TicketUpdated,
		Ticket: model.TicketNested{
			ID:      ticket.ID.Hex(),
			Subject: ticket.Subject,
		},
		Category:  *ticket.Category,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.mongodbRepo.CreateNotification(ctx, notification); err != nil {
		logrus.Error(err)
	}

	return nil

}

func _sendConfirmCloseTicketNotification(config model.Config, ticket *model.Ticket, company *model.Company) {
	//check customer email
	if ticket.Customer.Email == "" {
		logrus.Error("Customer email not found on ticket")
		return
	}

	closeTicketLink := helpers.StringReplacer(config.CloseTicketLink, map[string]string{
		"base_url_frontend": company.Settings.Domain.FullUrl,
		"token":             ticket.Token,
	})

	//mail content
	mailer := helpers.NewSMTPMailer(company)
	mailer.To([]string{ticket.Customer.Email})
	mailer.Subject(config.Email.Template.ConfirmCloseTicket.Title)
	mailer.Body(helpers.StringReplacer(config.Email.Template.ConfirmCloseTicket.Body, map[string]string{
		"title":          config.Email.Template.ConfirmCloseTicket.Title,
		"customer_name":  ticket.Customer.Name,
		"ticket_subject": ticket.Subject,
		"close_link":     closeTicketLink,
	}))

	//send mail
	if err := mailer.Send(); err != nil {
		logrus.Errorf("Send Email to %s error %v", ticket.Customer.Email, err)
	}
}

func (u *agentUsecase) ExportTicketsToCSV(ctx context.Context, claim domain.JWTClaimAgent, query url.Values, w http.ResponseWriter) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	fetchOptions := map[string]interface{}{
		"companyID": claim.CompanyID,
	}

	// Filtering
	if query.Get("status") != "" {
		fetchOptions["status"] = strings.Split(query.Get("status"), ",")
	} else {
		fetchOptions["status"] = []string{
			string(model.Open), string(model.InProgress), string(model.Resolve), string(model.Closed), string(model.Processing), string(model.Cancel),
		}
	}

	if startDateStr := query.Get("startDate"); startDateStr != "" {
		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err == nil {
			fetchOptions["startDate"] = startDate
		} else {
			return response.Error(http.StatusBadRequest, "Invalid startDate format")
		}
	}

	if endDateStr := query.Get("endDate"); endDateStr != "" {
		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err == nil {
			fetchOptions["endDate"] = endDate
		} else {
			return response.Error(http.StatusBadRequest, "Invalid endDate format")
		}
	}

	if query.Get("sort") != "" {
		fetchOptions["sort"] = query.Get("sort")
	}

	if query.Get("dir") != "" {
		fetchOptions["dir"] = query.Get("dir")
	}

	if query.Get("companyProductName") != "" {
		fetchOptions["companyProductName"] = query.Get("companyProductName")
	}

	if query.Get("companyProductID") != "" {
		fetchOptions["companyProductID"] = query.Get("companyProductID")
	}

	if query.Get("projectID") != "" {
		fetchOptions["projectID"] = query.Get("projectID")
	}

	if query.Get("customerID") != "" {
		fetchOptions["customerID"] = query.Get("customerID")
	}

	if query.Get("agentID") != "" {
		fetchOptions["agentID"] = query.Get("agentID")
	}

	if query.Get("status") != "" {
		fetchOptions["status"] = strings.Split(query.Get("status"), ",")
	}

	if query.Get("code") != "" {
		fetchOptions["code"] = query.Get("code")
	}

	if query.Get("subject") != "" {
		fetchOptions["subject"] = query.Get("subject")
	}

	// Fetch tickets from MongoDB
	cursor, err := u.mongodbRepo.FetchTicketList(ctx, fetchOptions)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "Error fetching tickets from MongoDB")
	}
	defer cursor.Close(ctx)

	// Prepare CSV headers and response writer
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=tickets.csv")

	// CSV writer
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	// Write CSV headers
	err = csvWriter.Write([]string{"ID", "Company", "Customer", "Subject", "Code", "Status", "Priority", "CreatedAt", "ClosedAt\\n"})
	if err != nil {
		return response.Error(http.StatusInternalServerError, "Error writing CSV header")
	}

	// Flag to check if any records are found
	foundRecords := false

	// write data to the CSV
	for cursor.Next(ctx) {
		foundRecords = true

		var ticket model.Ticket
		// Decode the current ticket
		if err := cursor.Decode(&ticket); err != nil {
			log.Printf("Error decoding ticket: %v", err)
			continue
		}

		// Prepare the CSV row
		row := []string{
			ticket.ID.Hex(),
			ticket.Company.Name,
			// ticket.Product.Name,
			ticket.Customer.Name,
			ticket.Subject,
			ticket.Code,
			string(ticket.Status),
			string(ticket.Priority),
			ticket.CreatedAt.Format(time.RFC3339),
		}

		// create closedAt
		if ticket.ClosedAt != nil {
			row = append(row, ticket.ClosedAt.Format(time.RFC3339))
		} else {
			row = append(row, "")
		}

		row[len(row)-1] += "\\n"

		// Write the row to the CSV
		if err := csvWriter.Write(row); err != nil {
			log.Printf("Error writing ticket to CSV: %v", err)
			continue
		}
	}

	// Check for cursor errors
	if err := cursor.Err(); err != nil {
		return response.Error(http.StatusInternalServerError, "Error iterating cursor")
	}

	// return an empty list
	if !foundRecords {
		return response.Success(domain.ResponseList{
			List: response.List{
				List: []interface{}{},
			},
		})
	}

	csvWriter.Flush()

	// Check for flushing errors
	if err := csvWriter.Error(); err != nil {
		return response.Error(http.StatusInternalServerError, "Error flushing CSV data")
	}

	return response.Success(domain.ResponseList{
		List: response.List{
			List: []interface{}{},
		},
	})
}

func (u *agentUsecase) _updateTimeBalance(ctx context.Context, ticket *model.Ticket, duration int) (err error) {
	customer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"id": ticket.Customer.ID,
	})
	if err != nil {
		return err
	}
	if customer == nil {
		return fmt.Errorf("customer not found")
	}

	// check need balance
	if !customer.IsNeedBalance {
		return
	}

	usedTime := int64(customer.Subscription.Balance.Time.Used) + int64(duration)

	if err := u.mongodbRepo.UpdateOneCustomer(
		ctx,
		map[string]interface{}{"id": customer.ID},
		map[string]interface{}{
			"subscription.balance.time.used": usedTime,
			"updatedAt":                      time.Now(),
		}); err != nil {
		return err

	}

	if err := u.mongodbRepo.CreateCustomerBalanceHistory(ctx, &model.CustomerBalanceHistory{
		ID:       primitive.NewObjectID(),
		Customer: ticket.Customer,
		Out:      int64(ticket.LogTime.DurationInSeconds),
		Reference: model.Reference{
			UniqueID: ticket.ID.Hex(),
			Type:     model.TicketReference},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}); err != nil {
		return err
	}

	return nil
}

func (u *agentUsecase) AssignTicketToMe(ctx context.Context, claim domain.JWTClaimAgent, ticketId string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check ticket
	ticket, err := u.mongodbRepo.FetchOneTicket(ctx, map[string]interface{}{
		"id":        ticketId,
		"companyID": claim.CompanyID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if ticket == nil {
		return response.Error(http.StatusBadRequest, "ticket not found")
	}

	// check ticket company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": ticket.Company.ID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if company == nil {
		return response.Error(http.StatusBadRequest, "company not found")
	}

	// check ticket company type
	if company.Type != "B2C" {
		return response.Error(http.StatusBadRequest, "company type must be B2C")
	}

	// check if agent is already assigned
	for _, assignedAgent := range ticket.Agent {
		if assignedAgent.ID == claim.User.ID {
			return response.Error(http.StatusBadRequest, "agent is already assigned to this ticket")
		}
	}

	// assign agent
	ticket.Agent = append(ticket.Agent, model.AgentNested{
		ID:    claim.User.ID,
		Name:  claim.User.Name,
		Email: claim.User.Email,
	})

	ticket.UpdatedAt = time.Now()

	if err := u.mongodbRepo.UpdateTicket(ctx, ticket); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(map[string]interface{}{
		"message": "Ticket successfully assigned",
		"ticket":  ticket,
	})
}

func (u *agentUsecase) GetTotalTicketCustomer(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	total := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		// "companyProductID": options["id"],
		"companyID": claim.CompanyID,
	})

	totalOpen := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		"status":    []string{"open"},
		"companyID": claim.CompanyID,
		// "companyProductID": options["id"],
	})

	totalInProgress := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		"status":    []string{"in_progress"},
		"companyID": claim.CompanyID,
		// "companyProductID": options["id"],
	})

	totalClosed := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		"status": []string{"closed"},
		// "companyProductID": options["id"],
	})

	result := map[string]interface{}{
		"totalTicket":           total,
		"totalTicketOpen":       totalOpen,
		"totalTicketClosed":     totalClosed,
		"totalTicketInProgress": totalInProgress,
	}

	return response.Success(result)
}

func (u *agentUsecase) GetDataCustomerTicket(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	filterOpen := map[string]interface{}{
		"status":           []string{"open"},
		"companyProductID": options["id"],
		"companyID":        claim.CompanyID,
	}

	filterInProgress := map[string]interface{}{
		"status":           []string{"in_progress"},
		"companyProductID": options["id"],
		"companyID":        claim.CompanyID,
	}

	filterClosed := map[string]interface{}{
		"status":           []string{"closed"},
		"companyProductID": options["id"],
		"companyID":        claim.CompanyID,
	}

	responseDay := make([]map[string]interface{}, 0)
	for _, day := range model.Weekdays {
		u._countCustomerTicketPerDay(ctx, day, filterOpen, filterClosed, filterInProgress, &responseDay)
	}

	return response.Success(responseDay)
}

func (u *agentUsecase) _countCustomerTicketPerDay(ctx context.Context, day string, optionsOpen, optionsClosed, optionsInProgress map[string]interface{}, responseDay *[]map[string]interface{}) {
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
