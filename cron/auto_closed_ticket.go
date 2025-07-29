package cronjob

import (
	"app/domain/model"
	"app/helpers"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

func (cj *cronjob) AutoCloseResolvedTickets() {
	cj.cron.AddFunc("*/1 * * * *", func() {
		t := time.Now()
		logrus.Info("AutoCloseResolvedTickets: cron started at ", t)

		// Tentukan batas waktu 7 hari yang lalu
		sevenDaysAgo := time.Now().Add(-7 * 24 * time.Hour)

		// Filter tiket yang statusnya 'resolved' dan 'updated_at' lebih dari 3 hari
		fetchOptions := map[string]interface{}{
			"status":    []string{"resolve"},
			"updatedAt": map[string]interface{}{"$lt": sevenDaysAgo},
		}

		// Ambil tiket yang memenuhi filter
		cur, err := cj.mongodbRepo.FetchTicketList(cj.ctx, fetchOptions)
		if err != nil {
			logrus.WithError(err).Error("Failed to fetch tickets")
			return
		}

		defer cur.Close(cj.ctx)

		for cur.Next(cj.ctx) {
			ticket := model.Ticket{}
			err := cur.Decode(&ticket)
			if err != nil {
				logrus.Error("Ticket Decode ", err)
				continue
			}

			// Pastikan ticket memenuhi kriteria untuk diubah statusnya
			if ticket.UpdatedAt.Before(sevenDaysAgo) && ticket.Status == model.Resolve {
				now := time.Now()
				// Find company
				company, err := cj.mongodbRepo.FetchOneCompany(cj.ctx, map[string]interface{}{"id": ticket.Company.ID})
				if err != nil {
					logrus.WithError(err).Error("Failed to fetch company")
					continue
				}
				if company == nil {
					logrus.Error("Company not found")
					continue
				}

				// Update status ticket menjadi 'closed'
				ticket.Status = model.Closed
				ticket.UpdatedAt = now
				ticket.ClosedAt = &now

				if err := cj.mongodbRepo.UpdateTicket(cj.ctx, &ticket); err != nil {
					logrus.WithFields(logrus.Fields{
						"ticketID":      ticket.ID.Hex(),
						"ticketSubject": ticket.Subject,
					}).Errorf("Failed to update ticket: %s", err.Error())
					continue
				}

				// Send email to the customer notifying them that their ticket has been closed
				mailer := helpers.NewSMTPMailer(company)
				if ticket.Customer.Email != "" {
					// Email content
					mailer.To([]string{ticket.Customer.Email})
					mailer.Subject(fmt.Sprintf("Ticket '%s' has been closed", ticket.Subject))
					mailer.Body(fmt.Sprintf(`
						Dear %s,<br><br>
						We are writing to inform you that your ticket, '%s', which was resolved three days ago, has now been closed.<br><br>
						Thank you for using our service. If you have any further issues, feel free to contact us.<br><br>
						Best regards,<br>
						Your Support Team
					`, ticket.Customer.Name, ticket.Subject))

					// Send the email
					if err := mailer.Send(); err != nil {
						logrus.WithFields(logrus.Fields{
							"email":         ticket.Customer.Email,
							"ticketSubject": ticket.Subject,
						}).Errorf("Failed to send email: %s", err.Error())
					} else {
						logrus.WithFields(logrus.Fields{
							"email":         ticket.Customer.Email,
							"ticketSubject": ticket.Subject,
						}).Info("Email sent successfully")
					}
				}
			}
		}
	})

	logrus.Info("Cron AutoCloseResolvedTickets added")
}