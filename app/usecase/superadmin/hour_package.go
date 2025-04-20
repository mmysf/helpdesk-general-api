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

func (u *superadminUsecase) GetHourPackages(ctx context.Context, claim domain.JWTClaimSuperadmin, query url.Values) response.Base {
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
	totalpackages := u.mongodbRepo.CountHourPackage(ctx, fetchOptions)
	if totalpackages == 0 {
		return response.Success(domain.ResponseList{
			List: response.List{
				List:  []interface{}{},
				Page:  page,
				Limit: limit,
				Total: totalpackages,
			},
			TotalPage: helpers.GetTotalPage(totalpackages, limit),
		})
	}

	// check hour package list
	packages, err := u.mongodbRepo.FetchHourPackageList(ctx, fetchOptions)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	defer packages.Close(ctx)

	list := make([]interface{}, 0)
	for packages.Next(ctx) {
		row := model.HourPackage{}
		err := packages.Decode(&row)
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
			Total: totalpackages,
		},
		TotalPage: helpers.GetTotalPage(totalpackages, limit),
	})
}

func (u *superadminUsecase) CreateHourPackage(ctx context.Context, claim domain.JWTClaimSuperadmin, payload domain.HourPackageRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)

	// validating
	if payload.DurationHours < 1 {
		errValidation["durationHours"] = "durationHours field is required"
	}

	if payload.Name == "" {
		errValidation["name"] = "name field is required"
	}

	if len(payload.Benefit) == 0 {
		errValidation["benefit"] = "benefit field is required"
	}

	if payload.Price <= 0 {
		errValidation["price"] = "price field cannot be 0"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// count duration in seconds
	totalInSeconds := payload.DurationHours * 60 * 60

	// create package
	newpackages := &model.HourPackage{
		ID:      primitive.NewObjectID(),
		Name:    payload.Name,
		Benefit: payload.Benefit,
		Price:   payload.Price,
		Duration: model.HourPackageDuration{
			Hours:          payload.DurationHours,
			TotalinSeconds: totalInSeconds,
		},
		Status:    model.HourPackageActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.mongodbRepo.CreateHourPackage(ctx, newpackages); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(newpackages)
}

func (u *superadminUsecase) GetHourPackageDetail(ctx context.Context, packageId string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if packageId == "" {
		return response.Error(http.StatusBadRequest, "packages id is required")
	}

	packages, err := u.mongodbRepo.FetchOneHourPackage(ctx, map[string]interface{}{
		"id": packageId,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if packages == nil {
		return response.Error(http.StatusBadRequest, "hour package not found")
	}

	return response.Success(packages)
}

func (u *superadminUsecase) UpdateHourPackage(ctx context.Context, claim domain.JWTClaimSuperadmin, packageId string, payload domain.HourPackageUpdate) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if packageId == "" {
		return response.Error(http.StatusBadRequest, "hour packages id is required")
	}

	errValidation := make(map[string]string)

	// validating
	if payload.DurationHours < 1 {
		errValidation["durationHours"] = "durationHours field is required"
	}

	if payload.Name == "" {
		errValidation["name"] = "name field is required"
	}

	if len(payload.Benefit) == 0 {
		errValidation["benefit"] = "benefit field is required"
	}

	if payload.Price <= 0 {
		errValidation["price"] = "price field cannot be 0"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check hour package
	packages, err := u.mongodbRepo.FetchOneHourPackage(ctx, map[string]interface{}{
		"id": packageId,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if packages == nil {
		return response.Error(http.StatusBadRequest, "hour package not found")
	}

	// count duration in seconds
	totalInSeconds := payload.DurationHours * 60 * 60

	// update package
	packages.Name = payload.Name
	packages.Benefit = payload.Benefit
	packages.Price = payload.Price
	packages.Duration.Hours = payload.DurationHours
	packages.Duration.TotalinSeconds = totalInSeconds

	packages.UpdatedAt = time.Now()

	if err := u.mongodbRepo.UpdateHourPackage(ctx, packages); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(packages)
}

func (u *superadminUsecase) UpdateStatusHourPackage(ctx context.Context, claim domain.JWTClaimSuperadmin, packageId string, payload domain.HourPackageStatusUpdate) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if packageId == "" {
		return response.Error(http.StatusBadRequest, "packages id is required")
	}

	errValidation := make(map[string]string)

	// validating

	if !helpers.InArrayString(payload.Status, []string{string(model.HourPackageActive), string(model.HourPackageInactive)}) {
		errValidation["status"] = "status only can be " + string(model.HourPackageActive) + " or " + string(model.HourPackageInactive)
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check hour package
	packages, err := u.mongodbRepo.FetchOneHourPackage(ctx, map[string]interface{}{
		"id": packageId,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if packages == nil {
		return response.Error(http.StatusBadRequest, "hour package not found")
	}

	// update package
	packages.Status = model.HourPackageStatus(payload.Status)

	packages.UpdatedAt = time.Now()

	if err := u.mongodbRepo.UpdateHourPackage(ctx, packages); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(packages)
}

func (u *superadminUsecase) DeleteHourPackage(ctx context.Context, claim domain.JWTClaimSuperadmin, packageId string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if packageId == "" {
		return response.Error(http.StatusBadRequest, "packages id is required")
	}

	packages, err := u.mongodbRepo.FetchOneHourPackage(ctx, map[string]interface{}{
		"id": packageId,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if packages == nil {
		return response.Error(http.StatusBadRequest, "packages not found")
	}

	now := time.Now()
	packages.UpdatedAt = now
	packages.DeletedAt = &now

	if err := u.mongodbRepo.UpdateHourPackage(ctx, packages); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(packages)
}
