package usecase_member

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

func (u *appUsecase) GetProjectList(ctx context.Context, claim domain.JWTClaimUser, query url.Values) response.Base {
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
	if query.Get("q") != "" {
		fetchOptions["q"] = query.Get("q")
	}

	// only own project if b2c
	if claim.Company.Type == "B2C" {
		fetchOptions["createdBy"] = claim.UserID
	}

	// count first
	totalDocuments := u.mongodbRepo.CountProject(ctx, fetchOptions)

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

	// check project list
	cur, err := u.mongodbRepo.FetchProjectList(ctx, fetchOptions)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	defer cur.Close(ctx)

	list := make([]interface{}, 0)
	for cur.Next(ctx) {
		row := model.Project{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("Project Decode ", err)
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

func (u *appUsecase) GetProjectDetail(ctx context.Context, claim domain.JWTClaimUser, projectID string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// filter
	fetchOptions := map[string]interface{}{
		"id":               projectID,
		"companyID":        claim.CompanyID,
		"companyProductID": claim.CompanyProductID,
	}

	// only own project if b2c
	if claim.Company.Type == "B2C" {
		fetchOptions["createdBy"] = claim.UserID
	}

	// check project
	project, err := u.mongodbRepo.FetchOneProject(ctx, fetchOptions)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if project == nil {
		return response.Error(http.StatusBadRequest, "project not found")
	}

	return response.Success(project)
}

func (u *appUsecase) CreateProject(ctx context.Context, claim domain.JWTClaimUser, payload domain.ProjectRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validating
	errValidation := make(map[string]string)
	if payload.Name == "" {
		errValidation["name"] = "name field is required"
	}
	if payload.Description == "" {
		errValidation["description"] = "description field is required"
	}
	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// create project
	project := model.Project{
		ID:             primitive.NewObjectID(),
		Company:        claim.Company,
		CompanyProduct: claim.CompanyProduct,
		Name:           payload.Name,
		Description:    payload.Description,
		CreatedAt:      time.Now(),
		CreatedBy: model.CustomerFK{
			ID:    claim.User.ID,
			Name:  claim.User.Name,
			Email: claim.User.Email,
		},
		UpdatedAt: time.Now(),
	}

	// save
	if err := u.mongodbRepo.CreateProject(ctx, &project); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(project)
}

func (u *appUsecase) UpdateProject(ctx context.Context, claim domain.JWTClaimUser, projectID string, payload domain.ProjectRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// jwt claim
	UserID := claim.UserID

	// validating
	errValidation := make(map[string]string)
	if payload.Name == "" {
		errValidation["name"] = "name field is required"
	}
	if payload.Description == "" {
		errValidation["description"] = "description field is required"
	}
	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// filter
	fetchOptions := map[string]interface{}{
		"id":               projectID,
		"companyID":        claim.CompanyID,
		"companyProductID": claim.CompanyProductID,
	}

	// only own project if b2c
	if claim.Company.Type == "B2C" {
		fetchOptions["createdBy"] = claim.UserID
	}

	// check project
	project, err := u.mongodbRepo.FetchOneProject(ctx, fetchOptions)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if project == nil {
		return response.Error(http.StatusBadRequest, "project not found")
	}

	// check customer
	customer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"id": UserID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if customer == nil {
		return response.Error(http.StatusUnauthorized, "customer not found")
	}

	// check company
	if customer.Company.ID != project.Company.ID {
		return response.Error(http.StatusBadRequest, "invalid customer company")
	}

	// check company product
	if customer.CompanyProduct.ID != project.CompanyProduct.ID {
		return response.Error(http.StatusBadRequest, "invalid customer company product")
	}

	// update project
	project.Description = payload.Description
	project.Name = payload.Name
	project.CreatedBy = model.CustomerFK{
		ID:    claim.User.ID,
		Name:  claim.User.Name,
		Email: claim.User.Email,
	}
	project.UpdatedAt = time.Now()

	// save
	if err := u.mongodbRepo.UpdateOneProject(ctx, project); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(project)
}

func (u *appUsecase) DeleteProject(ctx context.Context, claim domain.JWTClaimUser, projectID string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// filter
	fetchOptions := map[string]interface{}{
		"id":               projectID,
		"companyID":        claim.CompanyID,
		"companyProductID": claim.CompanyProductID,
	}

	// only own project if b2c
	if claim.Company.Type == "B2C" {
		fetchOptions["createdBy"] = claim.UserID
	}

	// check project
	project, err := u.mongodbRepo.FetchOneProject(ctx, fetchOptions)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if project == nil {
		return response.Error(http.StatusBadRequest, "project not found")
	}

	now := time.Now()

	// delete project
	project.DeletedAt = &now

	// save
	if err := u.mongodbRepo.UpdateOneProject(ctx, project); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(project)
}
