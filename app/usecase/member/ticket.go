package usecase_member

import (
	"app/domain"
	"app/domain/model"
	"app/helpers"
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	yurekahelpers "github.com/Yureka-Teknologi-Cipta/yureka/helpers"
	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *appUsecase) GetTicketList(ctx context.Context, claim domain.JWTClaimUser, query url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, limit, offset := yurekahelpers.GetLimitOffset(query)

	fetchOptions := map[string]interface{}{
		"limit":            limit,
		"offset":           offset,
		"companyID":        claim.CompanyID,
		"companyProductID": claim.CompanyProductID,
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

	if query.Get("customerID") != "" {
		fetchOptions["customerID"] = query.Get("customerID")
	}

	if query.Get("projectID") != "" {
		fetchOptions["projectID"] = query.Get("projectID")
	}

	if query.Get("status") != "" {
		fetchOptions["status"] = strings.Split(query.Get("status"), ",")
	}

	if query.Get("priority") != "" {
		fetchOptions["priority"] = query.Get("priority")
	}

	if query.Get("code") != "" {
		fetchOptions["code"] = query.Get("code")
	}

	if query.Get("categoryID") != "" {
		fetchOptions["categoryID"] = query.Get("categoryID")
	}

	// filter by customer id if b2c
	if claim.Company.Type == "B2C" {
		fetchOptions["customerID"] = claim.UserID
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
			return response.Success(response.List{
				List:  []interface{}{},
				Page:  page,
				Limit: limit,
				Total: totalDocuments,
			})
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

func (u *appUsecase) CreateTicket(ctx context.Context, claim domain.JWTClaimUser, payload domain.TicketRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating

	if payload.Subject == "" {
		errValidation["subject"] = "subject field is required"
	}

	if payload.Content == "" {
		errValidation["content"] = "content field is required"
	}

	if payload.Priority == "" {
		errValidation["priority"] = "priority field is required"
	}

	switch payload.Priority {
	case string(model.PriorityLow), string(model.PriorityMedium), string(model.PriorityHigh), string(model.PriorityCritical):
		break
	default:
		errValidation["priority"] = "priority field is invalid"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check customer
	customer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"id": claim.UserID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if customer == nil {
		return response.Error(http.StatusBadRequest, "customer not found")
	}

	helpers.CustomerBalanceFormat(customer)

	// check project
	var projectFK *model.ProjectFK
	if payload.ProjectId != "" {
		project, err := u.mongodbRepo.FetchOneProject(ctx, map[string]interface{}{
			"id":               payload.ProjectId,
			"companyID":        claim.CompanyID,
			"companyProductID": claim.CompanyProductID,
		})
		if err != nil {
			return response.Error(http.StatusInternalServerError, err.Error())
		}
		if project == nil {
			return response.Error(http.StatusBadRequest, "project not found")
		}
		// check project owner
		// if project.CompanyProduct.ID != claim.CompanyProductID {
		// 	return response.Error(http.StatusBadRequest, "project is not owned by company product")
		// }

		projectFK = &model.ProjectFK{
			ID:   project.ID.Hex(),
			Name: project.Name,
		}
	}

	// check category
	var ticketCategoryFK *model.TicketCategoryFK
	if payload.CategoryId != "" {
		ticketCategory, err := u.mongodbRepo.FetchOneTicketCategory(ctx, map[string]interface{}{
			"id":        payload.CategoryId,
			"companyID": claim.CompanyID,
		})
		if err != nil {
			return response.Error(http.StatusInternalServerError, err.Error())
		}
		if ticketCategory == nil {
			return response.Error(http.StatusBadRequest, "ticket category not found")
		}

		ticketCategoryFK = &model.TicketCategoryFK{
			ID:   ticketCategory.ID.Hex(),
			Name: ticketCategory.Name,
		}
	}

	now := time.Now()

	// get from config
	config := u._CacheConfig(ctx)

	// check is need balance
	if customer.IsNeedBalance {
		if !(now.After(customer.Subscription.StartAt) && now.Before(customer.Subscription.EndAt)) || customer.Subscription.Status != model.Active {
			return response.Error(http.StatusBadRequest, "you don't have active subscription")
		}

		if customer.Subscription.Balance == nil || customer.Subscription.Balance.Time.Remaining.Total < config.MinimumCredit {
			return response.Error(http.StatusBadRequest, "you don't have enough time balance")
		}
	}

	// count ticket
	countTicket := u.mongodbRepo.CountTicket(ctx, map[string]interface{}{
		"companyID":        claim.CompanyID,
		"companyProductID": claim.CompanyProductID,
		"today":            true,
	})

	// generate random char
	randomChar := helpers.RandomChar(len(claim.Company.Code))

	// generate ticket code
	ticketCode := helpers.GenerateFormattedCode(claim.Company.Code, countTicket+1, randomChar)

	// get detail cover media
	ticketAttachments := make([]model.AttachmentFK, 0)
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
				ID:   attachment.ID.Hex(),
				Name: attachment.Name,
				Size: attachment.Size,
				URL:  attachment.URL,
				// ExpiredUrlAt: nil,
				Type:        attachment.Type,
				ProviderKey: attachment.ProviderKey,
				IsPrivate:   attachment.IsPrivate,
			}

			ticketAttachments = append(ticketAttachments, attach)
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

	// create ticket
	ticket := &model.Ticket{
		ID:      primitive.NewObjectID(),
		Company: claim.Company,
		// Product:  claim.CompanyProduct,
		Project:  projectFK,
		Category: ticketCategoryFK,
		Customer: model.CustomerFK{
			ID:    claim.User.ID,
			Name:  claim.User.Name,
			Email: claim.User.Email,
		},
		Subject:     payload.Subject,
		Content:     payload.Content,
		Code:        ticketCode,
		Attachments: ticketAttachments,
		LogTime: model.LogTime{
			StartAt:           nil,
			EndAt:             nil,
			DurationInSeconds: 0,
			Status:            model.NotStarted,
		},
		Status:    model.Open,
		Priority:  model.TicketPriority(payload.Priority),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// set Detail time
	year, month, day := time.Now().Date()
	ticket.DetailTime = model.DetailTime{
		Year:    year,
		Month:   int(month),
		Day:     day,
		DayName: strings.ToLower(time.Now().Weekday().String()),
	}

	// check parent
	if payload.ParentId != "" {
		parentTicket, err := u.mongodbRepo.FetchOneTicket(ctx, map[string]interface{}{
			"companyID":  claim.CompanyID,
			"customerID": claim.UserID,
			"code":       payload.ParentId,
		})

		if err != nil {
			return response.Error(http.StatusInternalServerError, err.Error())
		}

		if parentTicket == nil {
			return response.Error(http.StatusBadRequest, "parent ticket not found")
		}

		ticket.Parent = &model.TicketNested{
			ID:       parentTicket.ID.Hex(),
			Subject:  parentTicket.Subject,
			Content:  parentTicket.Content,
			Priority: parentTicket.Priority,
		}
	}

	if err := u.mongodbRepo.CreateTicket(ctx, ticket); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// update customer last activity
	t := time.Now()
	customer.UpdatedAt = t
	customer.LastActivityAt = &t

	if err := u.mongodbRepo.UpdateOneCustomer(
		ctx,
		map[string]interface{}{"id": customer.ID},
		map[string]interface{}{
			"updatedAt":      customer.UpdatedAt,
			"lastActivityAt": customer.LastActivityAt,
		}); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// get company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": customer.Company.ID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if company == nil {
		return response.Error(http.StatusBadRequest, "company not found")
	}

	// get company agent
	companyAgents, err := u.mongodbRepo.FetchAgentList(ctx, map[string]interface{}{
		"companyID": company.ID.Hex(),
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// get all agent email
	agentEmails := make([]string, 0)
	for companyAgents.Next(ctx) {
		row := model.Agent{}
		err := companyAgents.Decode(&row)
		if err != nil {
			logrus.Error("Agent Decode ", err)
			return response.Error(http.StatusInternalServerError, err.Error())
		}
		agentEmails = append(agentEmails, row.Email)
	}

	// send notification
	go _sendTicketNotfication(config, ticket, agentEmails, company)

	// update total ticket
	u.mongodbRepo.IncrementOneCompany(ctx, company.ID.Hex(), map[string]int64{
		"ticketTotal": 1,
	})
	// u.mongodbRepo.IncrementOneCompanyProduct(ctx, customer.CompanyProduct.ID, map[string]int64{
	// 	"ticketTotal": 1,
	// })

	return response.Success(ticket)
}

func _sendTicketNotfication(config model.Config, ticket *model.Ticket, agentEmails []string, company *model.Company) {
	// assign helper mailer
	mailer := helpers.NewSMTPMailer(company)

	//setup mail content
	mailer.To(agentEmails)
	mailer.Subject(config.Email.Template.CreateTicket.Title)
	mailer.Body(helpers.StringReplacer(config.Email.Template.CreateTicket.Body, map[string]string{
		"title":           config.Email.Template.CreateTicket.Title,
		"customer_name":   ticket.Customer.Name,
		"ticket_subject":  ticket.Subject,
		"ticket_priority": string(ticket.Priority),
		// "brand_name":      ticket.Product.Name,
	}))

	//send mail
	if err := mailer.Send(); err != nil {
		logrus.WithFields(logrus.Fields{
			"ticketID": ticket.ID.Hex(),
			"subject":  ticket.Subject,
			"receiver": agentEmails,
		}).Errorf("Failed to send email: %s", err.Error())
	}
}

func (u *appUsecase) GetTicketDetail(ctx context.Context, claim domain.JWTClaimUser, ticketID string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check ticket
	ticket, err := u.mongodbRepo.FetchOneTicket(ctx, map[string]interface{}{
		"id":               ticketID,
		"companyID":        claim.CompanyID,
		"companyProductID": claim.CompanyProductID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if ticket == nil {
		return response.Error(http.StatusBadRequest, "ticket not found")
	}

	return response.Success(ticket)
}

func (u *appUsecase) CloseTicket(ctx context.Context, claim domain.JWTClaimUser, payload domain.CloseTicketRequest) response.Base {
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
		"id": payload.TicketId,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if ticket == nil {
		return response.Error(http.StatusBadRequest, "ticket not found")
	}

	// check customer
	customer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"id": ticket.Customer.ID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if customer == nil {
		return response.Error(http.StatusUnauthorized, "customer not found")
	}

	if ticket.Status != model.InProgress && ticket.Status != model.Resolve {
		return response.Error(400, "ticket only can be close if status is in progress or resolve")
	}

	// only create timelogs if ticket is in progress
	if ticket.Status == model.InProgress {
		ticket.LogTime.EndAt = &now
		ticket.LogTime.DurationInSeconds = int(now.Sub(*ticket.LogTime.StartAt).Seconds()) - ticket.LogTime.PauseDurationInSeconds
		ticket.LogTime.TotalDurationInSeconds += ticket.LogTime.DurationInSeconds
		ticket.LogTime.TotalPausedDurationInSeconds += ticket.LogTime.PauseDurationInSeconds
		ticket.LogTime.Status = model.Done

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

		// check ticket timelogs
		timelogs, err := u.mongodbRepo.FetchOneTicketlogs(ctx, map[string]interface{}{
			"ticket.id": ticket.ID.Hex(),
			"sort":      "createdAt",
			"dir":       "desc",
		})

		if err != nil {
			return response.Error(http.StatusInternalServerError, err.Error())
		}

		if timelogs == nil {
			return response.Error(http.StatusBadRequest, "time log not found")
		}

		timelogs.EndAt = &logsEndAt
		timelogs.DurationInSeconds += duration
		timelogs.UpdatedAt = &now

		if err := u.mongodbRepo.UpdateTicketlogs(ctx, timelogs); err != nil {
			return response.Error(http.StatusInternalServerError, err.Error())
		}

		// check need balance
		if customer.IsNeedBalance {
			// calculate time balance
			usedTime := customer.Subscription.Balance.Time.Used + int64(duration)

			// save update time balance customer
			if err := u.mongodbRepo.UpdateOneCustomer(ctx, map[string]interface{}{
				"id": customer.ID,
			}, map[string]interface{}{
				"subscription.balance.time.used": usedTime,
				"updatedAt":                      now,
			}); err != nil {
				return response.Error(http.StatusInternalServerError, err.Error())
			}

			// create customer balance history
			if err := u.mongodbRepo.CreateCustomerBalanceHistory(ctx, &model.CustomerBalanceHistory{
				ID:       primitive.NewObjectID(),
				Customer: ticket.Customer,
				Out:      int64(ticket.LogTime.DurationInSeconds),
				Reference: model.Reference{
					UniqueID: ticket.ID.Hex(),
					Type:     model.TicketReference,
				},
				CreatedAt: now,
				UpdatedAt: now,
			}); err != nil {
				return response.Error(http.StatusInternalServerError, err.Error())
			}
		}
	}

	// update ticket
	ticket.Status = model.Closed
	ticket.ClosedAt = &now
	ticket.UpdatedAt = now

	if err := u.mongodbRepo.UpdateTicket(ctx, ticket); err != nil {
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

	// get config
	config := u._CacheConfig(ctx)

	// check agent assigned
	if len(ticket.Agent) > 0 {
		// send notification to assigned agent
		agentEmails := make([]string, 0)
		for _, agent := range ticket.Agent {
			agentEmails = append(agentEmails, agent.Email)
		}
		go _sendCloseTicketNotification(config, ticket, agentEmails, company)
	} else {
		// send notification to all company agent
		companyAgents, err := u.mongodbRepo.FetchAgentList(ctx, map[string]interface{}{
			"companyID": ticket.Company.ID,
		})
		if err != nil {
			return response.Error(http.StatusInternalServerError, err.Error())
		}

		// get all company agent email
		agentEmails := make([]string, 0)
		for companyAgents.Next(ctx) {
			row := model.Agent{}
			err := companyAgents.Decode(&row)
			if err != nil {
				logrus.Error("Agent Decode ", err)
				return response.Error(http.StatusInternalServerError, err.Error())
			}
			agentEmails = append(agentEmails, row.Email)
		}
		go _sendCloseTicketNotification(config, ticket, agentEmails, company)
	}

	return response.Success(ticket)
}

func (u *appUsecase) CloseTicketByEmail(ctx context.Context, payload domain.CloseTicketbyEmailRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating request
	if payload.Token == "" {
		errValidation["token"] = "token field is required"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check the db
	ticket, err := u.mongodbRepo.FetchOneTicket(ctx, map[string]interface{}{
		"token": payload.Token,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if ticket == nil {
		return response.Error(http.StatusBadRequest, "close ticket token not valid")
	}

	//check ticket status
	if ticket.Status != model.InProgress && ticket.Status != model.Resolve {
		return response.Error(400, "ticket only can be close if status is in progress or resolve")
	}

	//update closed
	ticket.Status = model.Closed
	ticket.Token = ""
	ticket.UpdatedAt = time.Now()
	ticket.ClosedAt = &ticket.UpdatedAt

	//save
	if err := u.mongodbRepo.UpdateTicket(ctx, ticket); err != nil {
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

	// get config
	config := u._CacheConfig(ctx)

	// check agent assigned
	if len(ticket.Agent) > 0 {
		// send notification to assigned agent
		agentEmails := make([]string, 0)
		for _, agent := range ticket.Agent {
			agentEmails = append(agentEmails, agent.Email)
		}
		go _sendCloseTicketNotification(config, ticket, agentEmails, company)
	} else {
		// send notification to all company agent
		companyAgents, err := u.mongodbRepo.FetchAgentList(ctx, map[string]interface{}{
			"companyID": ticket.Company.ID,
		})
		if err != nil {
			return response.Error(http.StatusInternalServerError, err.Error())
		}

		// get all company agent email
		agentEmails := make([]string, 0)
		for companyAgents.Next(ctx) {
			row := model.Agent{}
			err := companyAgents.Decode(&row)
			if err != nil {
				logrus.Error("Agent Decode ", err)
				return response.Error(http.StatusInternalServerError, err.Error())
			}
			agentEmails = append(agentEmails, row.Email)
		}
		go _sendCloseTicketNotification(config, ticket, agentEmails, company)
	}

	return response.Success(ticket)
}

func _sendCloseTicketNotification(config model.Config, ticket *model.Ticket, agentEmails []string, company *model.Company) {
	// assign helper mailer
	mailer := helpers.NewSMTPMailer(company)

	//setup mail content
	mailer.To(agentEmails)
	mailer.Subject(config.Email.Template.CloseTicket.Title)
	mailer.Body(helpers.StringReplacer(config.Email.Template.CloseTicket.Body, map[string]string{
		"title":           config.Email.Template.CloseTicket.Title,
		"customer_name":   ticket.Customer.Name,
		"ticket_subject":  ticket.Subject,
		"ticket_priority": string(ticket.Priority),
		// "brand_name":      ticket.Product.Name,
	}))

	//send mail
	if err := mailer.Send(); err != nil {
		logrus.WithFields(logrus.Fields{
			"ticketID": ticket.ID.Hex(),
			"subject":  ticket.Subject,
			"receiver": agentEmails,
		}).Errorf("Failed to send email: %s", err.Error())
	}
}

func (u *appUsecase) CreateTicketComment(ctx context.Context, claim domain.JWTClaimUser, payload domain.TicketCommentRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating

	if payload.TicketId == "" {
		errValidation["ticketId"] = "ticketId field is required"
	}

	if payload.Content == "" {
		errValidation["content"] = "content field is required"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check ticket

	ticket, err := u.mongodbRepo.FetchOneTicket(ctx, map[string]interface{}{
		"id": payload.TicketId,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if ticket == nil {
		return response.Error(http.StatusBadRequest, "ticked not found")
	}

	// get detail attachments
	commentAttachments := make([]model.AttachmentFK, 0)
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
				ID:   attachment.ID.Hex(),
				Name: attachment.Name,
				Size: attachment.Size,
				URL:  attachment.URL,
				Type: attachment.Type,
				// ExpiredUrlAt: nil,
				ProviderKey: attachment.ProviderKey,
				IsPrivate:   attachment.IsPrivate,
			}

			commentAttachments = append(commentAttachments, attach)
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
		Company: ticket.Company,
		// Product: ticket.Product,
		Customer: model.CustomerFK{
			ID:    claim.User.ID,
			Name:  claim.User.Name,
			Email: claim.User.Name,
		},
		Ticket: model.TicketNested{
			ID:      ticket.ID.Hex(),
			Subject: ticket.Subject,
		},
		Content:     payload.Content,
		Sender:      model.CustomerSender,
		Attachments: commentAttachments,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// update customer last activity
	t := time.Now()

	if err := u.mongodbRepo.UpdateOneCustomer(
		ctx,
		map[string]interface{}{"id": claim.UserID},
		map[string]interface{}{
			"updatedAt":      t,
			"lastActivityAt": t,
		}); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if err := u.mongodbRepo.CreateTicketComment(ctx, ticketComment); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(ticketComment)
}

func (u *appUsecase) GetTicketCommentList(ctx context.Context, claim domain.JWTClaimUser, ticketId string, query url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, limit, offset := yurekahelpers.GetLimitOffset(query)

	fetchOptions := map[string]interface{}{
		"limit":            limit,
		"offset":           offset,
		"ticketID":         ticketId,
		"companyID":        claim.CompanyID,
		"companyProductID": claim.CompanyProductID,
	}

	// filtering
	if query.Get("sort") != "" {
		fetchOptions["sort"] = query.Get("sort")
	}

	if query.Get("dir") != "" {
		fetchOptions["dir"] = query.Get("dir")
	}

	if query.Get("q") != "" {
		fetchOptions["q"] = query.Get("q")
	}

	// count first
	totalDocuments := u.mongodbRepo.CountTicketComment(ctx, fetchOptions)

	if totalDocuments == 0 {
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

func (u *appUsecase) GetTicketCommentDetail(ctx context.Context, claim domain.JWTClaimUser, commentId string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check ticket
	ticket, err := u.mongodbRepo.FetchOneTicketComment(ctx, map[string]interface{}{
		"id":               commentId,
		"companyID":        claim.CompanyID,
		"companyProductID": claim.CompanyProductID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if ticket == nil {
		return response.Error(http.StatusBadRequest, "ticket not found")
	}

	return response.Success(ticket)
}

func (u *appUsecase) ReopenTicket(ctx context.Context, claim domain.JWTClaimUser, payload domain.ReopenTicketRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating
	if payload.TicketId == "" {
		errValidation["id"] = "id field is required"
	}

	// check ticket
	ticket, err := u.mongodbRepo.FetchOneTicket(ctx, map[string]interface{}{
		"id":               payload.TicketId,
		"companyID":        claim.CompanyID,
		"companyProductID": claim.CompanyProductID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if ticket == nil {
		return response.Error(http.StatusBadRequest, "ticket not found")
	}

	if ticket.Status != model.Closed {
		return response.Error(http.StatusBadRequest, "ticket only can be reopen if status is closed")
	}

	// check customer
	customer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"id": claim.UserID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if customer == nil {
		return response.Error(http.StatusBadRequest, "customer not found")
	}

	helpers.CustomerBalanceFormat(customer)

	now := time.Now()

	// get from config
	config := u._CacheConfig(ctx)

	// check is need balance
	if customer.IsNeedBalance {
		if !(now.After(customer.Subscription.StartAt) && now.Before(customer.Subscription.EndAt)) || customer.Subscription.Status != model.Active {
			return response.Error(http.StatusBadRequest, "you don't have active subscription")
		}

		if customer.Subscription.Balance == nil || customer.Subscription.Balance.Time.Remaining.Total < config.MinimumCredit {
			return response.Error(http.StatusBadRequest, "you don't have enough time balance")
		}
	}

	// update ticket
	ticket.Status = model.Open
	ticket.UpdatedAt = time.Now()
	ticket.ClosedAt = nil

	if err := u.mongodbRepo.UpdateTicket(ctx, ticket); err != nil {
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

	// check agent assigned
	if len(ticket.Agent) > 0 {
		// send notification to assigned agent
		agentEmails := make([]string, 0)
		for _, agent := range ticket.Agent {
			agentEmails = append(agentEmails, agent.Email)
		}
		go _sendReopenTicketNotification(config, ticket, agentEmails, company)
	} else {
		// send notification to all company agent
		companyAgents, err := u.mongodbRepo.FetchAgentList(ctx, map[string]interface{}{
			"companyID": ticket.Company.ID,
		})
		if err != nil {
			return response.Error(http.StatusInternalServerError, err.Error())
		}

		// get all company agent email
		agentEmails := make([]string, 0)
		for companyAgents.Next(ctx) {
			row := model.Agent{}
			err := companyAgents.Decode(&row)
			if err != nil {
				logrus.Error("Agent Decode ", err)
				return response.Error(http.StatusInternalServerError, err.Error())
			}
			agentEmails = append(agentEmails, row.Email)
		}
		go _sendReopenTicketNotification(config, ticket, agentEmails, company)
	}

	return response.Success(ticket)
}

func _sendReopenTicketNotification(config model.Config, ticket *model.Ticket, agentEmails []string, company *model.Company) {
	// assign helper mailer
	mailer := helpers.NewSMTPMailer(company)

	//setup mail content
	mailer.To(agentEmails)
	mailer.Subject(config.Email.Template.ReopenTicket.Title)
	mailer.Body(helpers.StringReplacer(config.Email.Template.ReopenTicket.Body, map[string]string{
		"title":           config.Email.Template.ReopenTicket.Title,
		"customer_name":   ticket.Customer.Name,
		"ticket_subject":  ticket.Subject,
		"ticket_priority": string(ticket.Priority),
		// "brand_name":      ticket.Product.Name,
	}))

	//send mail
	if err := mailer.Send(); err != nil {
		logrus.WithFields(logrus.Fields{
			"ticketID": ticket.ID.Hex(),
			"subject":  ticket.Subject,
			"receiver": agentEmails,
		}).Errorf("Failed to send email: %s", err.Error())
	}
}

func (u *appUsecase) CancelTicket(ctx context.Context, claim domain.JWTClaimUser, payload domain.CancelTicketRequest) response.Base {
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
		"id":               payload.TicketId,
		"companyID":        claim.CompanyID,
		"companyProductID": claim.CompanyProductID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if ticket == nil {
		return response.Error(http.StatusBadRequest, "ticked not found")
	}

	if ticket.Status != model.Open {
		return response.Error(http.StatusBadRequest, "ticket only can be cancel if status is open")
	}

	// update ticket
	ticket.Status = model.Cancel
	ticket.UpdatedAt = now

	if err := u.mongodbRepo.UpdateTicket(ctx, ticket); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(ticket)
}
