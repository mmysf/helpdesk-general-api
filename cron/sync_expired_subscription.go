package cronjob

import (
	"app/domain/model"
	"app/helpers"
	"time"

	"github.com/sirupsen/logrus"
)

func (cj *cronjob) SyncExpiredSubscription() {
	cj.cron.AddFunc("0 0 * * *", func() {
		t := time.Now()
		logrus.Info("SyncExpiredSubscription: cron started at ", t)

		fetchOptions := map[string]interface{}{
			"endExpDate": time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC),
			"status":     model.Active,
		}

		// check customer subscription list
		cur, err := cj.mongodbRepo.FetchCustomerSubscriptionList(cj.ctx, fetchOptions, false)
		if err != nil {
			logrus.Error("FetchCustomerSubscriptionList: ", err)
		}

		defer cur.Close(cj.ctx)

		for cur.Next(cj.ctx) {
			row := model.CustomerSubscription{}
			err := cur.Decode(&row)
			if err != nil {
				logrus.Error("Subscription Decode ", err)
			}

			go func() {
				// check customer
				customer, err := cj.mongodbRepo.FetchOneCustomer(cj.ctx, map[string]interface{}{
					"id": row.Customer.ID,
				})
				if err == nil && customer != nil {
					// update status
					row.Status = model.Expired
					row.UpdatedAt = time.Now()
					if err := cj.mongodbRepo.UpdatePartialCustomerSubscription(cj.ctx, map[string]interface{}{"id": customer.ID},
						map[string]interface{}{
							"status":    row.Status,
							"updatedAt": time.Now(),
						}); err != nil {
						logrus.Error("UpdatePartialCustomerSubscription: ", err)
					}

					if row.Order.Type == model.HOUR_TYPE {
						// update customer
						customer.Subscription.Status = model.Expired
						customer.UpdatedAt = time.Now()
						customer.Subscription.Balance = nil
						customer.Subscription.HourPackage = nil

						if err = cj.mongodbRepo.UpdateOneCustomer(cj.ctx, map[string]interface{}{"id": customer.ID},
							map[string]interface{}{
								"subscription.status":      customer.Subscription.Status,
								"subscription.balance":     customer.Subscription.Balance,
								"subscription.hourPackage": customer.Subscription.HourPackage,
								"updatedAt":                time.Now(),
							}); err != nil {
							logrus.Error("UpdateOneCustomer: ", err)
						}
					}

					// check company
					company, err := cj.mongodbRepo.FetchOneCompany(cj.ctx, map[string]interface{}{
						"id": customer.Company.ID,
					})
					if err == nil && company != nil {
						// get from config
						config := cj._CacheConfig(cj.ctx)

						// send email
						cj._sendEmailExpSubscription(config, &row, company)
					} else {
						logrus.Error("FetchOneCompany: ", err)
					}
				} else {
					logrus.Error("FetchOneCustomer: ", err)
				}

			}()
		}

	})

	logrus.Info("Cron SyncExpiredSubscription added")
}

func (cj *cronjob) _sendEmailExpSubscription(config model.Config, customerSubscription *model.CustomerSubscription, company *model.Company) (err error) {
	//get package name
	var packageName string
	if customerSubscription.Order.Type == model.HOUR_TYPE {
		packageName = customerSubscription.HourPackage.Name
	} else {
		packageName = customerSubscription.ServerPackage.Name
	}

	// send email
	mail := helpers.NewSMTPMailer(company)
	mail.To([]string{customerSubscription.Customer.Email})
	mail.Subject(config.Email.Template.PackageExpired.Title)
	mail.Body(helpers.StringReplacer(config.Email.Template.PackageExpired.Body, map[string]string{
		"package_name":  packageName,
		"customer_name": customerSubscription.Customer.Name,
	}))

	// send
	if err := mail.Send(); err != nil {
		logrus.Errorf("Send Email to %s error %v", customerSubscription.Customer.Email, err)
	}

	return
}
