package cron

import (
	"app/domain/model"
	"app/helpers"
	"context"
	"time"
	"github.com/sirupsen/logrus"
)

func (cs *cronScheduler) ConfirmNotification(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, cs.contextTimeout)
	defer cancel()

	//get 2x24 hours before now
	twoDaysAgo := time.Now().Add(-48 * time.Hour)

	//filter ticket
	fetchOptions := map[string]interface{}{
		"status":         []string{"resolve"},
		"reminderSent":   false,
		"logTime.status": "done",
		"logTime.endAt":  twoDaysAgo,
	}

	// get tickets
	cur, err := cs.mongodbRepo.FetchTicketList(ctx, fetchOptions)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch tickets")
	}

	defer cur.Close(ctx)

	tickets := []model.Ticket{}
	err = cur.All(ctx, &tickets)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch all tickets")
		return
	}

	//loop every ticket and send notification
	for _, ticket := range tickets {
		// find company
		company, err := cs.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
			"id": ticket.Company.ID,
		})
		if err != nil {
			logrus.WithError(err).Error("Failed to fetch company")
			continue
		}
		if company == nil {
			logrus.Error("Company not found")
			continue
		}

		mailer := helpers.NewSMTPMailer(company)

		if ticket.Customer.Email != "" {
			//default token
			defaultToken := helpers.RandomString(64)
			//mail content
			mailer.To([]string{ticket.Customer.Email})
			mailer.Subject("Action Required: Ticket " + ticket.Subject + " Resolution Confirmation")
			mailer.Body(
				"Dear " + ticket.Customer.Name + ",<br><br>" +
					"We hope this message finds you well. We are writing to inform you that your ticket, '" + ticket.Subject +
					"', was resolved two days ago. We kindly ask you to review the resolution, and if the issue has been fully addressed, " +
					"please confirm by closing the ticket by clicking the button below.<br><br>" +
					"<a href=\"https://ticket.solutionlab.id/close-ticket-by-email/" + defaultToken + "\" " +
					"style=\"background-color: #FF5733; color: white; padding: 10px 20px; text-align: center; text-decoration: none; display: inline-block; border-radius: 5px;\">Confirm Ticket Closure</a><br>" +
					"<p>or Copy paste the link below in your brower</p>" +
					"https://ticket.solutionlab.id/close-ticket-by-email/" + defaultToken + "<br><br>" +
					"Thank you for using our service, and we look forward to your confirmation.<br><br>",
			)

			//send mail
			if err := mailer.Send(); err != nil {
				logrus.WithFields(logrus.Fields{
					"email":         ticket.Customer.Email,
					"ticketSubject": ticket.Subject,
				}).Errorf("Failed to send email: %s", err.Error())
			} else {
				//log success
				logrus.WithFields(logrus.Fields{
					"email":         ticket.Customer.Email,
					"ticketSubject": ticket.Subject,
				}).Info("Email sent successfully")

				//update ticket
				ticket.ReminderSent = true
				ticket.Token = defaultToken
				if err := cs.mongodbRepo.UpdateTicket(ctx, &ticket); err != nil {
					logrus.WithFields(logrus.Fields{
						"ticketId":      ticket.ID.Hex(),
						"ticketSubject": ticket.Subject,
					}).Errorf("Failed to update ticket: %s", err.Error())
				}
			}
		}
	}
}
