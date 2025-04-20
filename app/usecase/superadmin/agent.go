package usecase_superadmin

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
	"golang.org/x/crypto/bcrypt"
)

func (u *superadminUsecase) GetAgents(ctx context.Context, claim domain.JWTClaimSuperadmin, query url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, limit, offset := yurekahelpers.GetLimitOffset(query)

	fetchOptions := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	if query.Get("name") != "" {
		fetchOptions["name"] = query.Get("name")
	}

	if query.Get("companyID") != "" {
		fetchOptions["companyID"] = query.Get("companyID")
	}

	// count first
	totalAgents := u.mongodbRepo.CountAgent(ctx, fetchOptions)
	if totalAgents == 0 {
		return response.Success(domain.ResponseList{
			List: response.List{
				List:  []interface{}{},
				Page:  page,
				Limit: limit,
				Total: totalAgents,
			},
			TotalPage: helpers.GetTotalPage(totalAgents, limit),
		})
	}

	// check agent list
	agents, err := u.mongodbRepo.FetchAgentList(ctx, fetchOptions)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	defer agents.Close(ctx)

	list := make([]interface{}, 0)
	for agents.Next(ctx) {
		row := model.Agent{}
		err := agents.Decode(&row)
		if err != nil {
			logrus.Error("Agent Decode ", err)
			return response.Error(http.StatusInternalServerError, err.Error())
		}
		list = append(list, row)
	}
	return response.Success(domain.ResponseList{
		List: response.List{
			List:  list,
			Page:  page,
			Limit: limit,
			Total: totalAgents,
		},
		TotalPage: helpers.GetTotalPage(totalAgents, limit),
	})

}

func (u *superadminUsecase) CreateAgent(ctx context.Context, claim domain.JWTClaimSuperadmin, payload domain.AccountRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating request
	if payload.Name == "" {
		errValidation["name"] = "name field is required"
	}
	if payload.CompanyID == "" {
		errValidation["companyId"] = "companyId field is required"
	}
	if payload.Email == "" {
		errValidation["email"] = "email field is required"
	} else if !helpers.IsValidEmail(payload.Email) {
		errValidation["email"] = "email field is invalid"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": payload.CompanyID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if company == nil {
		return response.Error(http.StatusBadRequest, "company not found")
	}

	// check email
	agent, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{
		"email": payload.Email,
	})
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error())
	}
	if agent != nil {
		return response.Error(http.StatusBadRequest, "email already registered")
	}

	// get from config
	config := u._CacheConfig(ctx)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(helpers.RandomString(5)), bcrypt.DefaultCost)

	newAgent := model.Agent{
		ID:        primitive.NewObjectID(),
		Name:      payload.Name,
		Email:     payload.Email,
		Password:  string(hashedPassword),
		Role:      model.AgentRole,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = u.mongodbRepo.CreateAgent(ctx, &newAgent)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error())
	}

	go _sendEmailAgentCrendential(config, newAgent, string(hashedPassword), company)

	return response.Success(newAgent)
}

func (u *superadminUsecase) GetAgentDetail(ctx context.Context, claim domain.JWTClaimSuperadmin, agentId string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if agentId == "" {
		return response.Error(http.StatusBadRequest, "agent id is required")
	}

	agent, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{
		"id": agentId,
	})
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error())
	}

	return response.Success(agent)
}

func (u *superadminUsecase) UpdateAgent(ctx context.Context, claim domain.JWTClaimSuperadmin, agentId string, payload domain.AccountRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	agent, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{
		"id": agentId,
	})
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error())
	}
	if agent == nil {
		return response.Error(http.StatusBadRequest, "agent not found")
	}

	agentEmail, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{
		"email": payload.Email,
	})
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error())
	}

	if agentEmail != nil && agentEmail.ID.Hex() != agent.ID.Hex() {
		return response.Error(http.StatusBadRequest, "email already registered")
	}

	errValidation := make(map[string]string)
	// validating request
	if payload.Name == "" {
		errValidation["name"] = "name field is required"
	}
	if payload.Email == "" {
		errValidation["email"] = "email field is required"
	} else if !helpers.IsValidEmail(payload.Email) {
		errValidation["email"] = "email field is invalid"
	}

	// update agent
	agent.Name = payload.Name
	agent.Email = payload.Email

	err = u.mongodbRepo.UpdateAgent(ctx, agent)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error())
	}

	return response.Success(agent)
}

func (u *superadminUsecase) DeleteAgent(ctx context.Context, claim domain.JWTClaimSuperadmin, agentId string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	agent, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{
		"id": agentId,
	})
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error())
	}

	now := time.Now()
	agent.UpdatedAt = now
	agent.DeletedAt = &now

	if err := u.mongodbRepo.UpdateAgent(ctx, agent); err != nil {
		return response.Error(http.StatusBadRequest, err.Error())
	}

	return response.Success(agent)
}

func (u *superadminUsecase) ResetPasswordAgent(ctx context.Context, claim domain.JWTClaimSuperadmin, agentId string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	agent, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{
		"id": agentId,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if agent == nil {
		return response.Error(http.StatusBadRequest, "agent not found")
	}

	defaultPassword := helpers.RandomString(5)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)

	now := time.Now()
	agent.UpdatedAt = now
	agent.Password = string(hashedPassword)

	err = u.mongodbRepo.UpdateAgent(ctx, agent)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error())
	}

	return response.Success(agent)
}
