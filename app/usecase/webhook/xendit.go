package usecase_webhook

import (
	"app/domain"
	"app/domain/model"
	"app/helpers"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (u *webhookUsecase) HandleWebhook(ctx context.Context, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// payload
	payload := options["payload"].(domain.SnapWebhookRequest)

	// check order
	order, err := u.mongodbRepo.FetchOneOrder(ctx, map[string]interface{}{
		"id": payload.ExternalID,
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

	// add webhook order to history
	order.Payment.Webhook.History = append(order.Payment.Webhook.History, payload)

	// update latest webhook to detail
	order.Payment.Webhook.Detail = payload

	// update status
	order.Status = model.STATUS_EXPIRED
	order.UpdatedAt = time.Now()
	if payload.Status == "PAID" {
		order.Status = model.STATUS_PAID
		order.PaidAt = &payload.PaidAt
		order.Invoice.PaymentMethod = payload.PaymentMethod
		order.Invoice.MerchantName = payload.MerchantName
		order.Invoice.BankCode = payload.BankCode
		order.Invoice.PaymentChannel = payload.PaymentChannel
		order.Invoice.PaymentDestination = payload.PaymentDestination

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

func (u *webhookUsecase) _createSubscription(ctx context.Context, order *model.Order, customer *model.Customer, oneHourPackage *model.HourPackage, _ *model.ServerPackage) (err error) {
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
	expiredAt := now

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

		expiredAt = now.AddDate(0, helpers.GetSubscriptionDuration(), 0)

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
			StartAt: now,
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
		expiredAt = now.AddDate(0, 0, int(order.ServerPackage.Validity*order.Amount))
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

func (u *webhookUsecase) _sendEmailCustomerOrder(config model.Config, order *model.Order, customer *model.Customer, company *model.Company, customerSubscription *model.CustomerSubscription) (err error) {
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
	// send email
	mail := helpers.NewSMTPMailer(company)
	mail.To([]string{customer.Email})
	mail.Subject(config.Email.Template.PackageActivated.Title)
	mail.Body(helpers.StringReplacer(config.Email.Template.PackageActivated.Body, map[string]string{
		"package_type":   cases.Title(language.English).String(string(order.Type)),
		"package_name":   packageName,
		"price":          strconv.FormatFloat(order.GrandTotal, 'f', 2, 64),
		"payment_method": order.Invoice.PaymentMethod,
		"order_number":   order.OrderNumber,
		"purchase_date":  order.CreatedAt.Format("2006-01-02"),
		"expire_date":    customerSubscription.ExpiredAt.Format("2006-01-02"),
		"login_link":     loginLink,
	}))

	// send
	if err := mail.Send(); err != nil {
		logrus.Errorf("Send Email to %s error %v", customer.Email, err)
	}

	return
}
