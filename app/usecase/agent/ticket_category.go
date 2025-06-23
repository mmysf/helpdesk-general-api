package usecase_agent

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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *agentUsecase) GetTicketCategoriesList(ctx context.Context, claim domain.JWTClaimAgent, query url.Values) response.Base {
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
	if query.Get("name") != "" {
		fetchOptions["name"] = query.Get("name")
	}

	// count first
	totalDocuments := u.mongodbRepo.CountTicketCategory(ctx, fetchOptions)
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

	// check ticket category list
	cur, err := u.mongodbRepo.FetchTicketCategoryList(ctx, fetchOptions)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	defer cur.Close(ctx)

	list := make([]interface{}, 0)
	for cur.Next(ctx) {
		row := model.TicketCategory{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("Ticket Category Decode ", err)
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

func (u *agentUsecase) GetTicketCategoryDetail(ctx context.Context, claim domain.JWTClaimAgent, id string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check ticket category
	ticketCategory, err := u.mongodbRepo.FetchOneTicketCategory(ctx, map[string]interface{}{
		"id":        id,
		"companyID": claim.CompanyID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if ticketCategory == nil {
		return response.Error(http.StatusBadRequest, "ticket category not found")
	}

	return response.Success(ticketCategory)
}

func (u *agentUsecase) CreateTicketCategory(ctx context.Context, claim domain.JWTClaimAgent, payload domain.TicketCategoryRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating request
	if payload.Name == "" {
		errValidation["name"] = "name field is required"
	}
	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	now := time.Now()

	// create ticket category
	ticketCategory := model.TicketCategory{
		ID:        primitive.NewObjectID(),
		Company:   claim.Company,
		Name:      payload.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.mongodbRepo.CreateTicketCategory(ctx, &ticketCategory); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(ticketCategory)
}

func (u *agentUsecase) UpdateTicketCategory(ctx context.Context, claim domain.JWTClaimAgent, id string, payload domain.TicketCategoryRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating request
	if payload.Name == "" {
		errValidation["name"] = "name field is required"
	}
	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check ticket category
	ticketCategory, err := u.mongodbRepo.FetchOneTicketCategory(ctx, map[string]interface{}{
		"id":        id,
		"companyID": claim.CompanyID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if ticketCategory == nil {
		return response.Error(http.StatusBadRequest, "ticket category not found")
	}

	// update ticket category
	ticketCategory.Name = payload.Name
	ticketCategory.UpdatedAt = time.Now()

	if err := u.mongodbRepo.UpdateOneTicketCategory(ctx, ticketCategory); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(ticketCategory)
}

func (u *agentUsecase) DeleteTicketCategory(ctx context.Context, claim domain.JWTClaimAgent, id string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check ticket category
	ticketCategory, err := u.mongodbRepo.FetchOneTicketCategory(ctx, map[string]interface{}{
		"id":        id,
		"companyID": claim.CompanyID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if ticketCategory == nil {
		return response.Error(http.StatusBadRequest, "ticket category not found")
	}

	now := time.Now()

	// delete ticket category
	ticketCategory.DeletedAt = &now

	if err := u.mongodbRepo.UpdateOneTicketCategory(ctx, ticketCategory); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(nil)
}
