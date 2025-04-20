package mongorepo

import (
	"app/domain/model"
	"app/helpers"
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

func generateQueryFilterTicketComment(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	if ticketID, ok := options["ticketID"].(string); ok {
		query["ticket.id"] = ticketID
	}

	if companyID, ok := options["companyID"].(string); ok {
		query["company.id"] = companyID
	}

	if customerID, ok := options["customerID"].(string); ok {
		query["customer.id"] = customerID
	}

	if q, ok := options["q"].(string); ok {
		regex := bson.M{
			"$regex": primitive.Regex{
				Pattern: q,
				Options: "i",
			},
		}
		query["content"] = regex
	}

	return query, mongoOptions
}

func (r *mongoDBRepo) CreateTicketComment(ctx context.Context, ticket *model.TicketComment) (err error) {
	_, err = r.Conn.Collection(r.TicketCommentCollection).InsertOne(ctx, ticket)
	if err != nil {
		logrus.Error("CreateTicketComment InsertOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) FetchTicketCommentList(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {
	query, findOptions := generateQueryFilterTicketComment(options, true)

	cur, err = r.Conn.Collection(r.TicketCommentCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchTicketCommentList Find:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CountTicketComment(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterTicketComment(options, true)

	total, err := r.Conn.Collection(r.TicketCommentCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountTicketComment", err)
		return 0
	}
	return
}

func (r *mongoDBRepo) FetchOneTicketComment(ctx context.Context, options map[string]interface{}) (row *model.TicketComment, err error) {
	query, _ := generateQueryFilterTicketComment(options, false)

	err = r.Conn.Collection(r.TicketCommentCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneTicketComment FindOne:", err)
		return
	}

	return
}

func (r *mongoDBRepo) UpdateTicketCommentPartial(ctx context.Context, id primitive.ObjectID, field map[string]interface{}) (err error) {
	_, err = r.Conn.Collection(r.TicketCommentCollection).UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{
		"$set": field,
	})
	if err != nil {
		logrus.Error("UpdateTicketPartialComment:", err)
		return err
	}
	return nil
}
