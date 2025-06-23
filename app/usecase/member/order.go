package usecase_member

import (
	"app/domain"
	"app/domain/model"
	"app/helpers"
	"context"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"time"

	yurekahelpers "github.com/Yureka-Teknologi-Cipta/yureka/helpers"
	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *appUsecase) GetOrderList(ctx context.Context, claim domain.JWTClaimUser, query url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, limit, offset := yurekahelpers.GetLimitOffset(query)

	// get from config
	config := u._CacheConfig(ctx)

	fetchOptions := map[string]interface{}{
		"limit":      limit,
		"offset":     offset,
		"customerID": claim.UserID,
	}

	// filtering
	if query.Get("sort") != "" {
		fetchOptions["sort"] = query.Get("sort")
	}
	if query.Get("dir") != "" {
		fetchOptions["dir"] = query.Get("dir")
	}
	if query.Get("types") != "" {
		fetchOptions["types"] = query.Get("types")
	}

	// count first
	totalDocuments := u.mongodbRepo.CountOrder(ctx, fetchOptions)

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

	// check order list
	cur, err := u.mongodbRepo.FetchOrderList(ctx, fetchOptions)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	defer cur.Close(ctx)

	list := make([]interface{}, 0)
	for cur.Next(ctx) {
		row := model.Order{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("Order Decode ", err)
			return response.Error(http.StatusInternalServerError, err.Error())
		}

		list = append(list, row.Format(&config))
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

func (u *appUsecase) GetOrderDetail(ctx context.Context, claim domain.JWTClaimUser, orderID string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check order
	order, err := u.mongodbRepo.FetchOneOrder(ctx, map[string]interface{}{
		"id": orderID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if order == nil {
		return response.Error(http.StatusBadRequest, "order not found")
	}

	// update order status if expired
	if order.Status == model.STATUS_PENDING && order.ExpiredAt.Before(time.Now()) {
		order.Status = model.STATUS_EXPIRED
		order.Payment.Status = string(model.STATUS_EXPIRED)
		order.UpdatedAt = time.Now()

		// update in background
		go func() {
			u.mongodbRepo.UpdateOneOrder(context.Background(), order)
		}()

	}

	// get from config
	config := u._CacheConfig(ctx)

	return response.Success(order.Format(&config))
}

func (u *appUsecase) CreateHourOrder(ctx context.Context, claim domain.JWTClaimUser, payload domain.OrderRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validating
	errValidation := make(map[string]string)
	if payload.PackageID == "" {
		errValidation["packageId"] = "packageId field is required"
	}
	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check hour package
	oneHourPackage, err := u.mongodbRepo.FetchOneHourPackage(ctx, map[string]interface{}{
		"id":     payload.PackageID,
		"status": string(model.HourPackageActive),
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if oneHourPackage == nil {
		return response.Error(http.StatusBadRequest, "hour package not found")
	}

	// count order
	count := u.mongodbRepo.CountOrder(ctx, map[string]interface{}{
		"today": true,
	})

	// get from config
	config := u._CacheConfig(ctx)

	// generate random char
	randomChar := helpers.RandomChar(3)

	// generate order number
	orderNumber := helpers.GenerateFormattedCode("TRX", count+1, randomChar)

	tax := 0
	adminFee := 0
	discount := 0

	subTotal := oneHourPackage.Price * 1
	grandTotal := subTotal + float64(adminFee) + float64(tax) - float64(discount)

	// create order
	order := &model.Order{
		ID: primitive.NewObjectID(),
		HourPackage: &model.HourPackageFK{
			ID:      oneHourPackage.ID.Hex(),
			Name:    oneHourPackage.Name,
			Hours:   oneHourPackage.Duration.Hours,
			Benefit: oneHourPackage.Benefit,
			Price:   oneHourPackage.Price,
		},
		Customer: model.CustomerFK{
			ID:    claim.User.ID,
			Name:  claim.User.Name,
			Email: claim.User.Email,
		},
		OrderNumber: orderNumber,
		Amount:      1,
		Type:        model.HOUR_TYPE,
		Tax:         float64(tax),
		AdminFee:    float64(adminFee),
		Discount:    float64(discount),
		SubTotal:    subTotal,
		GrandTotal:  grandTotal,
		Status:      model.STATUS_PENDING,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if !config.ManualPayment.IsActive {
		// do generete snap link
		result, err := u.xenditRepo.GenereteSnapLink(context.Background(), oneHourPackage, nil, *order)
		if err != nil || result.Status != 200 {
			if result.Status != 0 {
				return response.Error(400, result.Message)
			}

			return response.Error(400, err.Error())
		}
		respDataXendit, _ := result.Data.(domain.XenditGenereteSnapLinkResponseSuccess)
		// update snap order
		order.Payment.Status = respDataXendit.Status
		order.Invoice.InvoiceXenditId = respDataXendit.ID
		order.Invoice.InvoiceURL = respDataXendit.InvoiceURL
		order.Invoice.InvoiceExternalId = respDataXendit.ExternalID
		order.Payment.Snap = respDataXendit
		order.UpdatedAt = time.Now()
		order.ExpiredAt = respDataXendit.ExpiryDate
	} else {
		order.Payment.Status = string(model.STATUS_PENDING)
		order.Invoice.BankCode = config.ManualPayment.BankName
		order.Invoice.MerchantName = config.ManualPayment.AccountName
		order.Invoice.PaymentMethod = "MANUAL_PAYMENT"
		order.Invoice.PaymentChannel = config.ManualPayment.BankName
		order.Invoice.PaymentDestination = config.ManualPayment.AccountNumber
		order.Invoice.SwiftCode = config.ManualPayment.SwiftCode
		order.ExpiredAt = time.Now().Add(time.Second * time.Duration(config.ManualPayment.DurationInSecond))
	}

	// save
	if err := u.mongodbRepo.CreateOrder(ctx, order); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(order)
}

func (u *appUsecase) CreateServerOrder(ctx context.Context, claim domain.JWTClaimUser, payload domain.OrderRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validating
	errValidation := make(map[string]string)
	if payload.PackageID == "" {
		errValidation["packageId"] = "packageId field is required"
	}
	if payload.Amount < 1 {
		errValidation["amount"] = "amount field must be greater than 0"
	}
	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check server package
	oneServerPackage, err := u.mongodbRepo.FetchOneServerPackage(ctx, map[string]interface{}{
		"id":     payload.PackageID,
		"status": string(model.ServerPackageActive),
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if oneServerPackage == nil {
		return response.Error(http.StatusBadRequest, "server package not found")
	}
	if oneServerPackage.Customizable {
		return response.Error(http.StatusBadRequest, "server package is customizable, please contact our support")
	}

	// count order
	count := u.mongodbRepo.CountOrder(ctx, map[string]interface{}{
		"today": true,
	})

	// get from config
	config := u._CacheConfig(ctx)

	// generate random char
	randomChar := helpers.RandomChar(3)

	// generate order number
	orderNumber := helpers.GenerateFormattedCode("TRX", count+1, randomChar)

	tax := 0
	adminFee := 0
	discount := 0

	subTotal := oneServerPackage.Price * float64(payload.Amount)
	grandTotal := subTotal + float64(adminFee) + float64(tax) - float64(discount)

	// create order
	order := &model.Order{
		ID: primitive.NewObjectID(),
		ServerPackage: &model.ServerPackageFK{
			ID:           oneServerPackage.ID.Hex(),
			Name:         oneServerPackage.Name,
			Customizable: oneServerPackage.Customizable,
			Validity:     oneServerPackage.Validity,
			Benefit:      oneServerPackage.Benefit,
			Price:        oneServerPackage.Price,
		},
		Customer: model.CustomerFK{
			ID:    claim.User.ID,
			Name:  claim.User.Name,
			Email: claim.User.Email,
		},
		OrderNumber: orderNumber,
		Type:        model.SERVER_TYPE,
		Amount:      payload.Amount,
		Tax:         float64(tax),
		AdminFee:    float64(adminFee),
		Discount:    float64(discount),
		SubTotal:    subTotal,
		GrandTotal:  grandTotal,
		Status:      model.STATUS_PENDING,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if !config.ManualPayment.IsActive {
		// do generete snap link
		result, err := u.xenditRepo.GenereteSnapLink(context.Background(), nil, oneServerPackage, *order)
		if err != nil || result.Status != 200 {
			if result.Status != 0 {
				return response.Error(400, result.Message)
			}

			return response.Error(400, err.Error())
		}
		respDataXendit, _ := result.Data.(domain.XenditGenereteSnapLinkResponseSuccess)
		// update snap order
		order.Payment.Status = respDataXendit.Status
		order.Invoice.InvoiceXenditId = respDataXendit.ID
		order.Invoice.InvoiceURL = respDataXendit.InvoiceURL
		order.Invoice.InvoiceExternalId = respDataXendit.ExternalID
		order.Payment.Snap = respDataXendit
		order.UpdatedAt = time.Now()
		order.ExpiredAt = respDataXendit.ExpiryDate
	} else {
		order.Payment.Status = string(model.STATUS_PENDING)
		order.Invoice.BankCode = config.ManualPayment.BankName
		order.Invoice.MerchantName = config.ManualPayment.AccountName
		order.Invoice.PaymentMethod = "MANUAL_PAYMENT"
		order.Invoice.PaymentChannel = config.ManualPayment.BankName
		order.Invoice.SwiftCode = config.ManualPayment.SwiftCode
		order.Invoice.PaymentDestination = config.ManualPayment.AccountNumber
		order.ExpiredAt = time.Now().Add(time.Second * time.Duration(config.ManualPayment.DurationInSecond))
	}

	// save
	if err := u.mongodbRepo.CreateOrder(ctx, order); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(order)
}

func (u *appUsecase) UploadAttachmentOrder(ctx context.Context, claim domain.JWTClaimUser, payload domain.UploadAttachment, request *http.Request) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// declare variables
	var err error
	validation := make(map[string]string)
	var typeDocument, objName string
	var file multipart.File
	var uploadedFile *multipart.FileHeader
	var fileSize int64

	// validation
	file, uploadedFile, err = request.FormFile("file")
	if err != nil {
		validation["file"] = "file field is required"
	}

	typeDocument = uploadedFile.Header.Get("Content-Type")
	if !helpers.InArrayString(typeDocument, domain.AllowedImgMimeTypes) {
		validation["file"] = "field file is not valid type"
	}

	fileSize = uploadedFile.Size
	maxFileSize := int64(1 * 1024 * 1024) // 1 MB in bytes

	if fileSize > maxFileSize {
		validation["file"] = "file size exceeds the maximum limit of 1 MB"
	}
	defer file.Close()

	if len(validation) > 0 {
		return response.ErrorValidation(validation, "error validation")
	}

	category := helpers.GetCategoryByContentType(typeDocument, true)
	year, month, _ := time.Now().Date()
	objName = "orders/" + category + "/" + strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/" + helpers.GenerateCleanName(uploadedFile.Filename)

	uploadData, err := u.s3Repo.UploadFilePublic(objName, file, typeDocument)

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	typeDocument = helpers.GetCategoryByContentType(typeDocument, false)

	now := time.Now().UTC()
	media := model.Media{
		ID:           primitive.NewObjectID(),
		Name:         payload.Title,
		Provider:     "s3",
		ProviderKey:  objName,
		Type:         typeDocument,
		Category:     model.Other,
		Size:         fileSize,
		URL:          helpers.GeneratePublicURL(uploadData.URL),
		ExpiredUrlAt: nil,
		IsPrivate:    false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// insert into db
	err = u.mongodbRepo.CreateMedia(ctx, &media)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	return response.Success(media)
}

func (u *appUsecase) ConfirmOrder(ctx context.Context, claim domain.JWTClaimUser, payload domain.ConfrimOrderRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validating
	errValidation := make(map[string]string)
	if payload.OrderID == "" {
		errValidation["orderId"] = "orderId field is required"
	}

	if payload.AttachmentId == "" {
		errValidation["attachmentId"] = "attachmentId field is required"
	}

	if payload.AccountName == "" {
		errValidation["accountName"] = "accountName field is required"
	}

	if payload.AccountNumber == "" {
		errValidation["accountNumber"] = "accountNumber field is required"
	} else if !helpers.IsValidNumeric(payload.AccountNumber) {
		errValidation["accountNumber"] = "accountNumber field is not valid"
	}

	if payload.BankName == "" {
		errValidation["bankName"] = "bankName field is required"
	}

	if payload.Note == "" {
		errValidation["note"] = "note field is required"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check hour package
	order, err := u.mongodbRepo.FetchOneOrder(ctx, map[string]interface{}{
		"id":         payload.OrderID,
		"customerID": claim.UserID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if order == nil {
		return response.Error(http.StatusBadRequest, "order not found")
	}
	if order.Status != model.STATUS_PENDING {
		return response.Error(http.StatusBadRequest, "order not valid")
	}
	if order.Invoice.PaymentMethod != "MANUAL_PAYMENT" {
		return response.Error(http.StatusBadRequest, "payment method is not manual payment")
	}
	if order.ExpiredAt.Before(time.Now()) {
		return response.Error(http.StatusBadRequest, "order has been expired")

	}

	// check attachment
	attachment, err := u.mongodbRepo.FetchOneMedia(ctx, map[string]interface{}{
		"id":       payload.AttachmentId,
		"category": model.Other,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if attachment == nil {
		return response.Error(http.StatusBadRequest, "attachment media not found")
	}

	// check customer
	customer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"id": order.Customer.ID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if customer == nil {
		return response.Error(http.StatusUnauthorized, "customer not found")
	}

	manualPaid := &model.ManualPaid{
		BankName:      payload.BankName,
		AccountName:   payload.AccountName,
		AccountNumber: payload.AccountNumber,
		Note:          payload.Note,
		Attachment: model.MediaFK{
			ID:          attachment.ID.Hex(),
			Name:        attachment.Name,
			Size:        attachment.Size,
			URL:         attachment.URL,
			Type:        attachment.Type,
			Category:    attachment.Category,
			ProviderKey: attachment.ProviderKey,
			IsPrivate:   attachment.IsPrivate,
		},
		Approval: nil,
	}

	order.Status = model.STATUS_WAITING_APPROVAL
	order.Payment.Status = string(model.STATUS_WAITING_APPROVAL)
	order.UpdatedAt = time.Now()
	order.Payment.ManualPaid = manualPaid

	go u.mongodbRepo.UpdateManyMediaPartial(ctx, []primitive.ObjectID{attachment.ID}, map[string]interface{}{
		"isUsed": true,
	})

	// save
	if err := u.mongodbRepo.UpdateOneOrder(ctx, order); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(order)
}
