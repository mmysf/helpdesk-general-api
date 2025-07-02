package usecase_superadmin

import (
	"app/domain"
	"app/domain/model"
	"app/helpers"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	yurekahelpers "github.com/Yureka-Teknologi-Cipta/yureka/helpers"
	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *superadminUsecase) GetTotalTicket(ctx context.Context, claim domain.JWTClaimSuperadmin) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	total := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{})

	totalOpen := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		"status": []string{"open"},
	})

	totalInProgress := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		"status": []string{"in_progress"},
	})

	totalResolved := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		"status": []string{"resolve"},
	})

	totalClosed := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		"status": []string{"closed"},
	})

	result := map[string]interface{}{
		"total_ticket":             total,
		"total_ticket_open":        totalOpen,
		"total_ticket_resolved":    totalResolved,
		"total_ticket_closed":      totalClosed,
		"total_ticket_in_progress": totalInProgress,
	}

	return response.Success(result)
}

func (u *superadminUsecase) GetTicketList(ctx context.Context, claim domain.JWTClaimSuperadmin, query url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, limit, offset := yurekahelpers.GetLimitOffset(query)

	fetchOptions := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	// filtering
	if query.Get("sort") != "" {
		fetchOptions["sort"] = query.Get("sort")
	}
	if query.Get("dir") != "" {
		fetchOptions["dir"] = query.Get("dir")
	}
	if query.Get("subject") != "" {
		fetchOptions["subject"] = query.Get("subject")
	}
	if query.Get("code") != "" {
		fetchOptions["code"] = query.Get("code")
	}
	if query.Get("agentID") != "" {
		fetchOptions["agentID"] = query.Get("agentID")
	}
	if query.Get("companyProductID") != "" {
		fetchOptions["companyProductID"] = query.Get("companyProductID")
	}

	if query.Get("status") != "" {
		fetchOptions["status"] = strings.Split(query.Get("status"), ",")
	} else {
		fetchOptions["status"] = []string{string(model.Open), string(model.InProgress), string(model.Resolve), string(model.Cancel), string(model.Closed)}
	}

	if query.Get("priority") != "" {
		fetchOptions["priority"] = query.Get("priority")
	}
	if query.Get("companyID") != "" {
		fetchOptions["companyID"] = query.Get("companyID")
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
			TotalPage: helpers.GetTotalPage(totalDocuments, limit),
		})
	}

	// check ticket list
	cur, err := u.mongodbRepo.FetchTicketList(ctx, fetchOptions)

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	defer cur.Close(ctx)

	list := make([]interface{}, 0)
	for cur.Next(ctx) {
		row := model.Ticket{}
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

func (u *superadminUsecase) GetTicketDetail(ctx context.Context, ticketId string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if ticketId == "" {
		return response.Error(http.StatusBadRequest, "ticket id is required")
	}

	// check ticket
	ticket, err := u.mongodbRepo.FetchOneTicket(ctx, map[string]interface{}{
		"id": ticketId,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if ticket == nil {
		return response.Error(http.StatusBadRequest, "ticket not found")
	}

	return response.Success(ticket)
}

func (u *superadminUsecase) AssignAgent(ctx context.Context, claim domain.JWTClaimSuperadmin, ticketId string, payload domain.AssignAgentRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Validate agent names
	if len(payload.AgentIds) == 0 {
		return response.ErrorValidation(
			map[string]string{"agentIds": "agentIds is required"},
			"error validation",
		)
	}

	// check ticket
	ticket, err := u.mongodbRepo.FetchOneTicket(ctx, map[string]interface{}{
		"id": ticketId,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if ticket == nil {
		return response.Error(http.StatusBadRequest, "ticket not found")
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
	if company.Type != "B2C" {
		return response.Error(http.StatusBadRequest, "company is not B2C")
	}

	// Process each agent
	for _, agentId := range payload.AgentIds {
		// Check agent
		agent, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{
			"id": agentId,
		})
		if err != nil {
			return response.Error(http.StatusInternalServerError, err.Error())
		}
		if agent == nil {
			return response.Error(http.StatusBadRequest, fmt.Sprintf("agent '%s' not found", agentId))
		}

		// Check agent company
		if agent.Company.ID != ticket.Company.ID {
			return response.Error(http.StatusBadRequest, fmt.Sprintf("agent '%s' is not in company '%s'", agentId, ticket.Company.ID))
		}

		// Check if agent is already assigned
		for _, assignedAgent := range ticket.Agent {
			if assignedAgent.ID == agent.ID.Hex() {
				return response.Error(http.StatusBadRequest, "agent is already assigned to this ticket")
			}
		}

		// Assign agent
		newAgent := model.AgentNested{
			ID:    agent.ID.Hex(),
			Name:  agent.Name,
			Email: agent.Email,
		}
		ticket.Agent = append(ticket.Agent, newAgent)
	}

	if err := u.mongodbRepo.UpdateTicket(ctx, ticket); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(ticket)
}

func (u *superadminUsecase) GetDataClientTicket(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get id
	companyID := options["id"]

	filterOpen := map[string]interface{}{
		"status":    []string{"open"},
		"companyID": companyID,
	}

	filterInProgress := map[string]interface{}{
		"status":    []string{"in_progress"},
		"companyID": companyID,
	}

	filterClosed := map[string]interface{}{
		"status":    []string{"closed"},
		"companyID": companyID,
	}

	responseDay := make([]map[string]interface{}, 0)
	for _, day := range model.Weekdays {
		u._countClientTicketPerDay(ctx, day, filterOpen, filterClosed, filterInProgress, &responseDay)
	}

	return response.Success(responseDay)
}

func (u *superadminUsecase) _countClientTicketPerDay(ctx context.Context, day string, optionsOpen, optionsClosed, optionsInProgress map[string]interface{}, responseDay *[]map[string]interface{}) {
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

func (u *superadminUsecase) GetAverageDurationClient(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get id
	companyID := options["id"]

	fetchOptions := map[string]interface{}{
		"status":    []string{"closed"},
		"companyID": companyID,
	}

	responseDay := make([]map[string]interface{}, 0)

	// Calculate the average duration
	for _, day := range model.Weekdays {
		u._calculateAverageDurationPerDate(ctx, day, fetchOptions, &responseDay)
	}

	return response.Success(responseDay)
}

func (u *superadminUsecase) _calculateAverageDurationPerDate(ctx context.Context, day string, options map[string]interface{}, responseDay *[]map[string]interface{}) {
	options["day"] = day

	cur, err := u.mongodbRepo.FetchTicketList(ctx, options)
	if err != nil {
		logrus.Error("FetchTicketList:", err)
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
		*responseDay = append(*responseDay, map[string]interface{}{
			"day":             day,
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

	*responseDay = append(*responseDay, map[string]interface{}{
		"day":             day,
		"averageDuration": averageDuration,
	})
}

func (u *superadminUsecase) PauseLoggingTicket(ctx context.Context, claim domain.JWTClaimSuperadmin, payload domain.LoggingTicketRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating
	if payload.TicketId == "" {
		errValidation["id"] = "id field is required"
	}

	// check ticket
	ticket, err := u.mongodbRepo.FetchOneTicket(ctx, map[string]interface{}{
		"id": payload.TicketId,
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
	now := time.Now()
	if len(ticket.LogTime.PauseHistory) == 0 {
		ticket.LogTime.DurationInSeconds = int(now.Sub(*ticket.LogTime.StartAt).Seconds())
	} else {
		lastPause := ticket.LogTime.PauseHistory[len(ticket.LogTime.PauseHistory)-1]
		if lastPause.ResumedAt != nil {
			duration := int(now.Sub(*lastPause.ResumedAt).Seconds())
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
	timelogs.DurationInSeconds = int(now.Sub(*timelogs.StartAt).Seconds())

	if err := u.mongodbRepo.UpdateTicketlogs(ctx, timelogs); err != nil {
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

func (u *superadminUsecase) ResumeLoggingTicket(ctx context.Context, claim domain.JWTClaimSuperadmin, payload domain.LoggingTicketRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating
	if payload.TicketId == "" {
		errValidation["id"] = "id field is required"
	}

	// check ticket
	ticket, err := u.mongodbRepo.FetchOneTicket(ctx, map[string]interface{}{
		"id": payload.TicketId,
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
