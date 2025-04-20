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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *superadminUsecase) GetServerPackageList(ctx context.Context, claim domain.JWTClaimSuperadmin, query url.Values) response.Base {
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
	if query.Get("dir") != "" {
		fetchOptions["dir"] = query.Get("dir")
	}
	if query.Get("q") != "" {
		fetchOptions["q"] = query.Get("q")
	}

	// count first
	totalServerPackage := u.mongodbRepo.CountServerPackage(ctx, fetchOptions)
	if totalServerPackage == 0 {
		return response.Success(domain.ResponseList{
			List: response.List{
				List:  []interface{}{},
				Page:  page,
				Limit: limit,
				Total: totalServerPackage,
			},
			TotalPage: helpers.GetTotalPage(totalServerPackage, limit),
		})
	}

	// check package list
	serverPackages, err := u.mongodbRepo.FetchServerPackageList(ctx, fetchOptions)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	list := make([]interface{}, 0)
	for serverPackages.Next(ctx) {
		row := model.ServerPackage{}
		err := serverPackages.Decode(&row)
		if err != nil {
			return response.Error(http.StatusInternalServerError, err.Error())
		}
		list = append(list, row)
	}

	return response.Success(domain.ResponseList{
		List: response.List{
			List:  list,
			Page:  page,
			Limit: limit,
			Total: totalServerPackage,
		},
		TotalPage: helpers.GetTotalPage(totalServerPackage, limit),
	})
}

func (u *superadminUsecase) CreateServerPackage(ctx context.Context, claim domain.JWTClaimSuperadmin, payload domain.ServerPackageRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)

	// validating
	if payload.Name == "" {
		errValidation["name"] = "name field is required"
	}

	if len(payload.Benefit) == 0 {
		errValidation["benefit"] = "benefit field is required"
	}

	if !payload.Customizable && payload.Price <= 0 {
		errValidation["price"] = "price field must be greater than 0"
	}

	if payload.Validity < 1 {
		errValidation["validity"] = "validity field must be greater than 0"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// create server package
	newServerPackages := &model.ServerPackage{
		ID:           primitive.NewObjectID(),
		Name:         payload.Name,
		Benefit:      payload.Benefit,
		Price:        payload.Price,
		Customizable: payload.Customizable,
		Validity:     payload.Validity,
		Status:       model.ServerPackageActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := u.mongodbRepo.CreateServerPackage(ctx, newServerPackages); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(newServerPackages)
}

func (u *superadminUsecase) GetServerPackageDetail(ctx context.Context, ServerPackageId string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validating
	if ServerPackageId == "" {
		return response.Error(http.StatusBadRequest, "packages id is required")
	}

	serverPackages, err := u.mongodbRepo.FetchOneServerPackage(ctx, map[string]interface{}{
		"id": ServerPackageId,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if serverPackages == nil {
		return response.Error(http.StatusBadRequest, "server package not found")
	}

	return response.Success(serverPackages)
}

func (u *superadminUsecase) UpdateServerPackage(ctx context.Context, claim domain.JWTClaimSuperadmin, serverpackageId string, payload domain.ServerPackageUpdate) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)

	// validating
	if payload.Name == "" {
		errValidation["name"] = "name field is required"
	}

	if len(payload.Benefit) == 0 {
		errValidation["benefit"] = "benefit field is required"
	}

	if !payload.Customizable && payload.Price <= 0 {
		errValidation["price"] = "price field must be greater than 0"
	}

	if payload.Validity < 1 {
		errValidation["validity"] = "validity field must be greater than 0"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// fetch server package
	serverPackages, err := u.mongodbRepo.FetchOneServerPackage(ctx, map[string]interface{}{
		"id": serverpackageId,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if serverPackages == nil {
		return response.Error(http.StatusBadRequest, "server package not found")
	}

	// update server package
	serverPackages.Name = payload.Name
	serverPackages.Benefit = payload.Benefit
	serverPackages.Price = payload.Price
	serverPackages.Customizable = payload.Customizable
	serverPackages.Validity = payload.Validity
	serverPackages.UpdatedAt = time.Now()

	if err := u.mongodbRepo.UpdateServerPackage(ctx, serverPackages); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	return response.Success(serverPackages)
}

func (u *superadminUsecase) UpdateStatusServerPackage(ctx context.Context, claim domain.JWTClaimSuperadmin, serverPackageId string, payload domain.ServerPackageStatusUpdate) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if serverPackageId == "" {
		return response.Error(http.StatusBadRequest, "packages id is required")
	}

	errValidation := make(map[string]string)

	// validating

	if !helpers.InArrayString(payload.Status, []string{string(model.ServerPackageActive), string(model.ServerPackageInactive)}) {
		errValidation["status"] = "status only can be " + string(model.ServerPackageActive) + " or " + string(model.ServerPackageInactive)
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check server hour package
	serverPackages, err := u.mongodbRepo.FetchOneServerPackage(ctx, map[string]interface{}{
		"id": serverPackageId,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if serverPackages == nil {
		return response.Error(http.StatusBadRequest, "server package not found")
	}

	serverPackages.Status = model.ServerPackageStatus(payload.Status)
	serverPackages.UpdatedAt = time.Now()

	if err := u.mongodbRepo.UpdateServerPackage(ctx, serverPackages); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(serverPackages)
}

func (u *superadminUsecase) DeleteServerPackage(ctx context.Context, claim domain.JWTClaimSuperadmin, serverPackageId string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if serverPackageId == "" {
		return response.Error(http.StatusBadRequest, "packages id is required")
	}

	serverPackages, err := u.mongodbRepo.FetchOneServerPackage(ctx, map[string]interface{}{
		"id": serverPackageId,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if serverPackages == nil {
		return response.Error(http.StatusBadRequest, "server package not found")
	}

	now := time.Now()
	serverPackages.UpdatedAt = now
	serverPackages.DeletedAt = &now

	if err := u.mongodbRepo.UpdateServerPackage(ctx, serverPackages); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success("server package deleted successfully")
}
