package usecase_superadmin

import (
	"app/domain"
	"app/domain/model"
	"app/helpers"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	yurekahelpers "github.com/Yureka-Teknologi-Cipta/yureka/helpers"
	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (u *superadminUsecase) GetOrderList(ctx context.Context, claim domain.JWTClaimSuperadmin, query url.Values) response.Base {
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
	if query.Get("packageId") != "" {
		fetchOptions["packageId"] = query.Get("packageId")
	}
	if query.Get("status") != "" {
		fetchOptions["status"] = strings.Split(query.Get("status"), ",")
	}
	if query.Get("q") != "" {
		fetchOptions["q"] = query.Get("q")
	}
	if query.Get("types") != "" {
		fetchOptions["types"] = query.Get("types")
	}

	// count first
	totalOrder := u.mongodbRepo.CountOrder(ctx, fetchOptions)
	if totalOrder == 0 {
		return response.Success(domain.ResponseList{
			List: response.List{
				List:  []interface{}{},
				Page:  page,
				Limit: limit,
				Total: totalOrder,
			},
			TotalPage: helpers.GetTotalPage(totalOrder, limit),
		})
	}

	// check order
	order, err := u.mongodbRepo.FetchOrderList(ctx, fetchOptions)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	defer order.Close(ctx)

	// get from config
	config := u._CacheConfig(ctx)

	list := make([]interface{}, 0)
	for order.Next(ctx) {
		row := model.Order{}
		err := order.Decode(&row)
		if err != nil {
			return response.Error(http.StatusInternalServerError, err.Error())
		}
		list = append(list, row.Format(&config))
	}

	return response.Success(domain.ResponseList{
		List: response.List{
			List:  list,
			Page:  page,
			Limit: limit,
			Total: totalOrder,
		},
		TotalPage: helpers.GetTotalPage(totalOrder, limit),
	})

}

func (u *superadminUsecase) GetOrderDetail(ctx context.Context, orderID string) response.Base {
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

func (u *superadminUsecase) UploadAttachmentOrder(ctx context.Context, claim domain.JWTClaimSuperadmin, payload domain.UploadAttachment, request *http.Request) response.Base {
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

func (u *superadminUsecase) UpdateManualPayment(ctx context.Context, claim domain.JWTClaimSuperadmin, orderID string, payload domain.UpdateManualPaymentRequest) response.Base {
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

	// validating
	errValidation := make(map[string]string)
	if payload.Status == "" {
		errValidation["status"] = "status field is required"
	}
	if !helpers.InArrayString(payload.Status, []string{string(model.STATUS_REJECT), string(model.STATUS_PAID)}) {
		errValidation["status"] = "status only can be " + string(model.STATUS_REJECT) + " or " + string(model.STATUS_PAID)
	}
	if payload.Note == "" {
		errValidation["note"] = "note field is required"
	}
	if payload.Status == string(model.STATUS_PAID) && order.Status != model.STATUS_WAITING_APPROVAL {
		if payload.AccountName == "" {
			errValidation["accountName"] = "account name field is required"
		}
		if payload.AccountNumber == "" {
			errValidation["accountNumber"] = "account number field is required"
		} else if !helpers.IsValidNumeric(payload.AccountNumber) {
			errValidation["accountNumber"] = "accountNumber field is not valid"
		}
		if payload.BankName == "" {
			errValidation["bankName"] = "bank name field is required"
		}
		if payload.AttachmentId == "" {
			errValidation["attachmentId"] = "attachment id field is required"
		}
	}
	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	if order.Status == model.STATUS_PAID {
		return response.Error(http.StatusBadRequest, "order already paid")
	}
	if order.Invoice.PaymentMethod != "MANUAL_PAYMENT" {
		return response.Error(http.StatusBadRequest, "payment method is not manual payment")
	}

	// check customer
	customer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"id": order.Customer.ID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if customer == nil {
		return response.Error(http.StatusBadRequest, "customer not found")
	}

	// update order
	now := time.Now()
	if order.Payment.ManualPaid == nil {
		order.Payment.ManualPaid = &model.ManualPaid{}
		if payload.Status == string(model.STATUS_PAID) {
			order.Payment.ManualPaid.BankName = payload.BankName
			order.Payment.ManualPaid.AccountName = payload.AccountName
			order.Payment.ManualPaid.AccountNumber = payload.AccountNumber
			order.Payment.ManualPaid.Note = payload.Note
		}
	}

	order.Status = model.OrderStatus(payload.Status)
	order.Payment.Status = string(model.OrderStatus(payload.Status))
	order.Payment.ManualPaid.Approval = &model.Approval{
		User: claim.User,
		Note: payload.Note,
		At:   now,
	}
	order.UpdatedAt = now
	if payload.Status == string(model.STATUS_PAID) {
		if payload.AttachmentId != "" && order.Status != model.STATUS_WAITING_APPROVAL {
			attachmentFk := &model.MediaFK{}
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

			attachmentFk = &model.MediaFK{
				ID:          attachment.ID.Hex(),
				Name:        attachment.Name,
				Size:        attachment.Size,
				URL:         attachment.URL,
				Type:        attachment.Type,
				Category:    attachment.Category,
				ProviderKey: attachment.ProviderKey,
				IsPrivate:   attachment.IsPrivate,
			}

			go u.mongodbRepo.UpdateManyMediaPartial(ctx, []primitive.ObjectID{attachment.ID}, map[string]interface{}{
				"isUsed": true,
			})
			order.Payment.ManualPaid.Attachment = *attachmentFk
		}

		order.Status = model.STATUS_PAID
		order.Payment.Status = string(model.STATUS_PAID)
		order.Payment.PaidAt = &now
		order.PaidAt = &now

		if order.Type == model.HOUR_TYPE {
			pkg := &model.HourPackage{
				ID: func() primitive.ObjectID {
					id, _ := primitive.ObjectIDFromHex(order.HourPackage.ID)
					return id
				}(),
				Name:    order.HourPackage.Name,
				Benefit: order.HourPackage.Benefit,
				Price:   order.HourPackage.Price,
				Duration: model.HourPackageDuration{
					TotalinSeconds: order.HourPackage.Hours * 60 * 60,
					Hours:          order.HourPackage.Hours,
				},
			}
			// create customer subscription
			if err := u._createSubscription(ctx, order, customer, pkg, nil); err != nil {
				return response.Error(http.StatusInternalServerError, err.Error())
			}
			// create customer balance history
			if err := u.mongodbRepo.CreateCustomerBalanceHistory(ctx, &model.CustomerBalanceHistory{
				ID:       primitive.NewObjectID(),
				Customer: order.Customer,
				In:       pkg.Duration.TotalinSeconds,
				Reference: model.Reference{
					UniqueID: order.ID.Hex(),
					Type:     model.OrderReference,
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}); err != nil {
				return response.Error(http.StatusInternalServerError, err.Error())
			}
		} else {
			pkg := &model.ServerPackage{
				ID: func() primitive.ObjectID {
					id, _ := primitive.ObjectIDFromHex(order.ServerPackage.ID)
					return id
				}(),
				Name:         order.ServerPackage.Name,
				Benefit:      order.ServerPackage.Benefit,
				Price:        order.ServerPackage.Price,
				Customizable: order.ServerPackage.Customizable,
				Validity:     order.ServerPackage.Validity,
			}

			// create customer subscription
			if err := u._createSubscription(ctx, order, customer, nil, pkg); err != nil {
				return response.Error(http.StatusInternalServerError, err.Error())
			}
		}

	}

	// save
	if err := u.mongodbRepo.UpdateOneOrder(ctx, order); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(order)
}

func (u *superadminUsecase) _createSubscription(ctx context.Context, order *model.Order, customer *model.Customer, oneHourPackage *model.HourPackage, _ *model.ServerPackage) (err error) {
	// get from config
	config := u._CacheConfig(ctx)

	// check company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": customer.Company.ID,
	})
	if err != nil {
		return err
	}
	if company == nil {
		return fmt.Errorf("company not found")
	}

	// get current time
	now := time.Now()

	// check active subscription
	fetchOptions := map[string]interface{}{
		"customerID": customer.ID.Hex(),
		"orderType":  order.Type,
	}

	if order.Type == model.SERVER_TYPE {
		fetchOptions["serverPackageId"] = order.ServerPackage.ID
	} else {
		fetchOptions["status"] = model.Active
	}

	activeSubscription, err := u.mongodbRepo.FetchOneCustomerSubscription(ctx, fetchOptions)
	if err != nil {
		return err
	}

	//set expiry date
	expiredAt := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)

	if order.Type == model.HOUR_TYPE {
		// set active subscription to expired
		if activeSubscription != nil {
			activeSubscription.Status = model.Expired
			activeSubscription.UpdatedAt = now

			// save
			if err = u.mongodbRepo.UpdateOneCustomerSubscription(ctx, activeSubscription); err != nil {
				return err
			}
		}

		expiredAt = expiredAt.AddDate(0, helpers.GetSubscriptionDuration(), 0)

		// add time balance
		timeBalance := oneHourPackage.Duration.TotalinSeconds
		if balance := customer.Subscription.Balance; balance != nil {
			timeBalance += balance.Time.Total
			balance.Time.Total = timeBalance
		} else {
			customer.Subscription.Balance = &model.Balance{
				Time: model.TimeBalance{
					Total: timeBalance,
				},
			}
		}

		// set subscription to customer
		customer.Subscription = &model.Subscription{
			Status: model.Active,
			HourPackage: &model.HourPackageFK{
				ID:      oneHourPackage.ID.Hex(),
				Name:    oneHourPackage.Name,
				Hours:   oneHourPackage.Duration.Hours,
				Benefit: oneHourPackage.Benefit,
				Price:   oneHourPackage.Price,
			},
			Balance: customer.Subscription.Balance,
			StartAt: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			EndAt:   expiredAt,
		}
		customer.UpdatedAt = now

		// save
		if err = u.mongodbRepo.UpdateOneCustomer(ctx, map[string]interface{}{
			"id": customer.ID.Hex(),
		}, map[string]interface{}{
			"subscription": customer.Subscription,
			"updatedAt":    customer.UpdatedAt,
		}); err != nil {
			return
		}
	} else {
		expiredAt = expiredAt.AddDate(0, 0, int(order.ServerPackage.Validity*order.Amount))
		if activeSubscription != nil {
			if activeSubscription.ExpiredAt.Before(now) {
				expiredAt = now.AddDate(0, 0, int(order.ServerPackage.Validity*order.Amount))
			} else {
				expiredAt = activeSubscription.ExpiredAt.AddDate(0, 0, int(order.ServerPackage.Validity*order.Amount))
			}
			activeSubscription.Status = model.Active
			activeSubscription.ExpiredAt = expiredAt
			activeSubscription.UpdatedAt = now

			// save
			if err = u.mongodbRepo.UpdateOneCustomerSubscription(ctx, activeSubscription); err != nil {
				return err
			}
		}
	}

	if order.Type == model.HOUR_TYPE || activeSubscription == nil {
		// create customer subscription
		newCustomerSubscription := &model.CustomerSubscription{
			ID: primitive.NewObjectID(),
			Customer: model.CustomerFK{
				ID:    customer.ID.Hex(),
				Name:  customer.Name,
				Email: customer.Email,
			},
			HourPackage:   order.HourPackage,
			ServerPackage: order.ServerPackage,
			Order: model.OrderFK{
				ID:          order.ID.Hex(),
				OrderNumber: order.OrderNumber,
				Type:        order.Type,
			},
			Status:    model.Active,
			ExpiredAt: expiredAt,
			CreatedAt: now,
			UpdatedAt: now,
		}

		// save
		if err = u.mongodbRepo.CreateCustomerSubscription(ctx, newCustomerSubscription); err != nil {
			return
		}

		// send email
		go u._sendEmailCustomerOrder(config, order, customer, company, newCustomerSubscription)
	} else {

		// send email
		go u._sendEmailCustomerOrder(config, order, customer, company, activeSubscription)
	}

	return
}

func (u *superadminUsecase) _sendEmailCustomerOrder(config model.Config, order *model.Order, customer *model.Customer, company *model.Company, customerSubscription *model.CustomerSubscription) (err error) {
	loginLink := helpers.StringReplacer(config.LoginLink, map[string]string{
		"base_url_frontend": company.Settings.Domain.FullUrl,
	})

	//get package name
	var packageName string
	if order.Type == model.HOUR_TYPE {
		packageName = order.HourPackage.Name
	} else {
		packageName = order.ServerPackage.Name
	}

	// Format price
	rupiah := order.GrandTotal * config.DollarInIdr
	formattedPrice := helpers.FormatFloat("#,###.##", order.GrandTotal)
	formattedPriceRp := helpers.FormatFloat("#,###.##", rupiah)
	// send email
	mail := helpers.NewSMTPMailer(company)
	mail.To([]string{customer.Email})
	mail.Subject(config.Email.Template.PackageActivated.Title)
	mail.Body(helpers.StringReplacer(config.Email.Template.PackageActivated.Body, map[string]string{
		"package_type":   cases.Title(language.English).String(string(order.Type)),
		"package_name":   packageName,
		"price":          "USD " + formattedPrice + " (IDR " + formattedPriceRp + ")",
		"payment_method": order.Invoice.PaymentMethod,
		"order_number":   order.OrderNumber,
		"purchase_date":  order.CreatedAt.Format("2006-01-02"),
		"expire_date":    customerSubscription.ExpiredAt.Format("2006-01-02"),
		"login_link":     loginLink,
		"customer_name":  order.Customer.Name,
	}))

	// send
	if err := mail.Send(); err != nil {
		logrus.Errorf("Send Email to %s error %v", customer.Email, err)
	}

	return
}
