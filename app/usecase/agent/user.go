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
	"golang.org/x/crypto/bcrypt"
)

func (u *agentUsecase) GetAgentList(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	paramQuery := options["query"].(url.Values)
	page, limit, offset := yurekahelpers.GetLimitOffset(paramQuery)

	fetchOptions := map[string]interface{}{
		"limit":     limit,
		"offset":    offset,
		"companyID": claim.CompanyID,
	}

	// filtering
	if paramQuery.Get("sort") != "" {
		fetchOptions["sort"] = paramQuery.Get("sort")
	}

	if paramQuery.Get("dir") != "" {
		fetchOptions["dir"] = paramQuery.Get("dir")
	}

	if paramQuery.Get("q") != "" {
		fetchOptions["q"] = paramQuery.Get("q")
	}

	// count first
	totalDocuments := u.mongodbRepo.CountAgent(ctx, fetchOptions)

	if totalDocuments == 0 {
		return response.Success(response.List{
			List:  []interface{}{},
			Page:  page,
			Limit: limit,
			Total: totalDocuments,
		})
	}

	// check agent list
	cur, err := u.mongodbRepo.FetchAgentList(ctx, fetchOptions)

	if err != nil {
		return response.Success(response.List{
			List:  []interface{}{},
			Page:  page,
			Limit: limit,
			Total: totalDocuments,
		})
	}

	defer cur.Close(ctx)

	list := make([]interface{}, 0)
	for cur.Next(ctx) {
		row := model.Agent{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("Agent Decode ", err)
			return response.Success(response.List{
				List:  []interface{}{},
				Page:  page,
				Limit: limit,
				Total: totalDocuments,
			})
		}

		list = append(list, row)
	}

	return response.Success(response.List{
		List:  list,
		Page:  page,
		Limit: limit,
		Total: totalDocuments,
	})
}

func (u *agentUsecase) GetAgentDetail(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check agent
	agent, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{
		"id": options["id"],
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if agent == nil {
		return response.Error(http.StatusUnauthorized, "agent not found")
	}

	return response.Success(agent)
}

func (u *agentUsecase) CreateAgent(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// payload
	payload := options["payload"].(domain.CreateUserRequest)

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

	if payload.JobTitle == "" {
		errValidation["jobTitle"] = "jobTitle field is required"
	}

	if payload.Role == "" {
		errValidation["role"] = "role field is required"
	} else if !helpers.InArrayString(payload.Role, []string{"admin", "agent"}) {
		errValidation["role"] = "role field must be admin or agent"
	}
	if payload.CategoryId == "" {
		errValidation["category"] = "category field is required"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check agent
	agent, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{
		"email": payload.Email,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if agent != nil {
		return response.Error(http.StatusBadRequest, "email already in use for "+agent.Company.Name)
	}

	//check category
	category, err := u.mongodbRepo.FetchOneTicketCategory(ctx, map[string]interface{}{
		"id":        payload.CategoryId,
		"companyID": claim.CompanyID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if category == nil {
		return response.Error(http.StatusBadRequest, "category not found")
	}

	password := helpers.RandomString(5)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// create agent
	newAgent := &model.Agent{
		ID:        primitive.NewObjectID(),
		Name:      payload.Name,
		Email:     payload.Email,
		Password:  string(hashedPassword),
		Company:   claim.Company,
		JobTitle:  payload.JobTitle,
		Role:      model.UserRole(payload.Role),
		Category:  model.TicketCategoryFK{ID: category.ID.Hex(), Name: category.Name},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = u.mongodbRepo.CreateAgent(ctx, newAgent)

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// get from config
	config := u._CacheConfig(ctx)

	// find company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": claim.CompanyID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if company == nil {
		return response.Error(http.StatusUnauthorized, "company not found")
	}

	go _sendEmailAgentCrendential(config, *newAgent, password, company)

	return response.Success(newAgent)
}

func _sendEmailAgentCrendential(config model.Config, agent model.Agent, password string, company *model.Company) {
	loginLink := helpers.StringReplacer(config.LoginLink, map[string]string{
		"base_url_frontend": config.AgentLink,
	})
	// send email
	mail := helpers.NewSMTPMailer(company)
	mail.To([]string{agent.Email})
	mail.Subject(config.Email.Template.DefaultUser.Title)
	mail.Body(helpers.StringReplacer(config.Email.Template.DefaultUser.Body, map[string]string{
		"title":      config.Email.Template.DefaultUser.Title,
		"email":      agent.Email,
		"password":   password,
		"login_link": loginLink,
	}))

	// send
	if err := mail.Send(); err != nil {
		logrus.Errorf("Send Email to %s error %v", agent.Email, err)
	}
}

func (u *agentUsecase) UpdateAgent(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base {
	payload := options["payload"].(domain.CreateUserRequest)

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
	if payload.JobTitle == "" {
		errValidation["jobTitle"] = "jobTitle field is required"
	}
	if payload.Role == "" {
		errValidation["role"] = "role field is required"
	} else if !helpers.InArrayString(payload.Role, []string{"admin", "agent"}) {
		errValidation["role"] = "role field must be admin or agent"
	}
	if payload.CategoryId == "" {
		errValidation["category"] = "category field is required"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check agent
	agent, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{
		"id": options["id"],
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if agent == nil {
		return response.Error(http.StatusUnauthorized, "agent not found")
	}

	// check agent
	existingAgent, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{
		"email": payload.Email,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if existingAgent != nil && existingAgent.ID.Hex() != agent.ID.Hex() {
		return response.Error(http.StatusBadRequest, "email already in use for "+existingAgent.Company.Name)
	}

	//check category
	category, err := u.mongodbRepo.FetchOneTicketCategory(ctx, map[string]interface{}{
		"id":        payload.CategoryId,
		"companyID": claim.CompanyID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if category == nil {
		return response.Error(http.StatusBadRequest, "category not found")
	}

	// update agent
	agent.Name = payload.Name
	agent.Email = payload.Email
	agent.JobTitle = payload.JobTitle
	agent.Role = model.UserRole(payload.Role)
	agent.Category = model.TicketCategoryFK{ID: category.ID.Hex(), Name: category.Name}
	agent.UpdatedAt = time.Now()

	if err := u.mongodbRepo.UpdateOneAgent(
		ctx,
		map[string]interface{}{"id": agent.ID},
		map[string]interface{}{
			"name":      agent.Name,
			"jobTitle":  agent.JobTitle,
			"email":     agent.Email,
			"role":      agent.Role,
			"category":  agent.Category,
			"updatedAt": agent.UpdatedAt,
		}); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(agent)
}

func (u *agentUsecase) DeleteAgent(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base {
	// check agent
	agent, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{
		"id": options["id"],
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if agent == nil {
		return response.Error(http.StatusUnauthorized, "agent not found")
	}

	// validate cant delete self
	if agent.ID.Hex() == claim.UserID {
		return response.Error(http.StatusBadRequest, "cannot delete your own account")
	}

	t := time.Now()

	// update agent
	agent.UpdatedAt = t
	agent.DeletedAt = &t

	if err := u.mongodbRepo.UpdateOneAgent(
		ctx,
		map[string]interface{}{"id": agent.ID},
		map[string]interface{}{
			"updatedAt": agent.UpdatedAt,
			"deletedAt": agent.DeletedAt,
		}); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success("Agent deleted successfully")
}
