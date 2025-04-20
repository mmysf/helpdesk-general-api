package usecase_superadmin

import (
	"app/domain"
	"app/domain/model"
	"app/helpers"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	yurekahelpers "github.com/Yureka-Teknologi-Cipta/yureka/helpers"
	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *superadminUsecase) GetTicketCommentList(ctx context.Context, claim domain.JWTClaimSuperadmin, ticketId string, query url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, limit, offset := yurekahelpers.GetLimitOffset(query)

	fetchOptions := map[string]interface{}{
		"limit":    limit,
		"offset":   offset,
		"ticketID": ticketId,
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
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	defer cur.Close(ctx)

	list := make([]interface{}, 0)
	for cur.Next(ctx) {
		row := model.TicketComment{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("Ticket comment decode", err)
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

func (u *superadminUsecase) GetTicketCommentDetail(ctx context.Context, claim domain.JWTClaimSuperadmin, commentId string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check ticket
	ticket, err := u.mongodbRepo.FetchOneTicketComment(ctx, map[string]interface{}{
		"id": commentId,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if ticket == nil {
		return response.Error(http.StatusBadRequest, "ticket not found")
	}

	return response.Success(ticket)
}

func (u *superadminUsecase) CreateTicketComment(ctx context.Context, claim domain.JWTClaimSuperadmin, payload domain.SuperadminTicketCommentRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	now := time.Now()

	errValidation := make(map[string]string)
	// validating

	if payload.AgentId == "" {
		errValidation["agentId"] = "agentId field is required"
	}
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

	// check agent
	agent, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{
		"id": payload.AgentId,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if agent == nil {
		return response.Error(http.StatusBadRequest, "agent not found")
	}

	// check ticket
	ticket, err := u.mongodbRepo.FetchOneTicket(ctx, map[string]interface{}{
		"id":        payload.TicketId,
		"companyID": agent.Company.ID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if ticket == nil {
		return response.Error(http.StatusBadRequest, "ticket not found")
	}

	// check ticket company
	if ticket.Company.ID != agent.Company.ID {
		return response.Error(http.StatusBadRequest, "ticket is not in agent's company")
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

	// check company type for assigned agent
	if company.Type == "B2C" {
		assigned := false
		for _, agentAssigned := range ticket.Agent {
			if agentAssigned.ID == agent.ID.Hex() {
				assigned = true
				break
			}
		}
		if len(ticket.Agent) == 0 {
			return response.Error(http.StatusBadRequest, "the agent is not assigned to this ticket")
		}
		if !assigned {
			return response.Error(http.StatusBadRequest, "the agent is not assigned to this ticket")
		}
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

	// validate status
	switch ticket.Status {
	case model.Open:
		if !helpers.InArrayString(string(payload.Status), []string{string(model.Open), string(model.InProgress)}) {
			return response.Error(http.StatusBadRequest, "status allowed only open or in progress")
		}
	case model.InProgress:
		if !helpers.InArrayString(string(payload.Status), []string{string(model.InProgress), string(model.Resolve)}) {
			return response.Error(http.StatusBadRequest, "status allowed only in progress or resolve")
		}
	case model.Resolve:
		if !helpers.InArrayString(string(payload.Status), []string{string(model.InProgress), string(model.Resolve)}) {
			return response.Error(http.StatusBadRequest, "status allowed only in progress or resolve")
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
			"company_id": ticket.Company.ID,
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

	// set agent nested
	agentNested := model.AgentNested{
		ID:   agent.ID.Hex(),
		Name: agent.Name,
	}

	// create ticket comment
	ticketComment := &model.TicketComment{
		ID: primitive.NewObjectID(),
		Company: model.CompanyNested{
			ID:    ticket.Company.ID,
			Name:  ticket.Company.Name,
			Image: ticket.Company.Image,
			Type:  ticket.Company.Type,
		},
		Product: model.CompanyProductNested{
			ID:    ticket.Product.ID,
			Name:  ticket.Product.Name,
			Image: ticket.Product.Image,
			Code:  ticket.Product.Code,
		},
		Agent: agentNested,
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
	if err := u._updateTicketAndTimelog(ctx, ticket, payload.Status, agentNested); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// get from config
	config := u._CacheConfig(ctx)

	// notification
	go _sendCommentNotfication(config, ticket, ticketComment, &agentNested, company)

	return response.Success(ticketComment)
}

func _sendCommentNotfication(config model.Config, ticket *model.Ticket, ticketComment *model.TicketComment, agent *model.AgentNested, company *model.Company) {
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

func (u *superadminUsecase) _updateTicketAndTimelog(ctx context.Context, ticket *model.Ticket, status model.TicketStatus, agent model.AgentNested) (err error) {
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
				Product:  ticket.Product,
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
				Product:  ticket.Product,
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
			//check pause
			if len(ticket.LogTime.PauseHistory) == 0 {
				ticket.LogTime.DurationInSeconds = int(now.Sub(*ticket.LogTime.StartAt).Seconds())
			} else { //count active duration after last pause
				lastPause := ticket.LogTime.PauseHistory[len(ticket.LogTime.PauseHistory)-1]
				if lastPause.ResumedAt != nil {
					duration := int(now.Sub(*lastPause.ResumedAt).Seconds())
					ticket.LogTime.DurationInSeconds += duration
				}
			}
			ticket.LogTime.TotalDurationInSeconds += ticket.LogTime.DurationInSeconds
			ticket.LogTime.TotalPausedDurationInSeconds += ticket.LogTime.PauseDurationInSeconds
			ticket.LogTime.Status = model.Done
			ticket.ReminderSent = true
			ticket.Token = defaultToken
			ticket.Status = model.Resolve
			ticket.UpdatedAt = now

			if err := u.mongodbRepo.UpdateTicket(ctx, ticket); err != nil {
				return err
			}

			if err := u._updateTimeBalance(ctx, ticket); err != nil {
				return err
			}

			// find company
			company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
				"id": ticket.Company.ID,
			})
			if err != nil {
				return err
			}
			if company == nil {
				return err
			}

			// get from config
			config := u._CacheConfig(ctx)

			go _sendConfirmCloseTicketNotification(config, ticket, company)

			// create ticket timelogs
			u.mongodbRepo.CreateTicketTimelogs(ctx, &model.TicketTimeLogs{
				ID:       primitive.NewObjectID(),
				Company:  ticket.Company,
				Customer: ticket.Customer,
				Product:  ticket.Product,
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
				ActivityTtype:          "ticket resolve",
				CreatedAt:              now,
			})
		}
	}

	return nil
}

func _sendConfirmCloseTicketNotification(config model.Config, ticket *model.Ticket, company *model.Company) {
	//check customer email
	if ticket.Customer.Email == "" {
		logrus.Error("Customer email not found on ticket")
		return
	}

	closeTicketLik := helpers.StringReplacer(config.CloseTicketLink, map[string]string{
		"token": ticket.Token,
	})

	//mail content
	mailer := helpers.NewSMTPMailer(company)
	mailer.To([]string{ticket.Customer.Email})
	mailer.Subject(config.Email.Template.ConfirmCloseTicket.Title)
	mailer.Body(helpers.StringReplacer(config.Email.Template.ConfirmCloseTicket.Body, map[string]string{
		"title":          config.Email.Template.ConfirmCloseTicket.Title,
		"customer_name":  ticket.Customer.Name,
		"ticket_subject": ticket.Subject,
		"close_link":     closeTicketLik,
	}))

	//send mail
	if err := mailer.Send(); err != nil {
		logrus.Errorf("Send Email to %s error %v", ticket.Customer.Email, err)
	}
}

func (u *superadminUsecase) _updateTimeBalance(ctx context.Context, ticket *model.Ticket) (err error) {
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

	// count balance
	usedTime := int64(customer.Subscription.Balance.Time.Used) + int64(ticket.LogTime.DurationInSeconds)

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
