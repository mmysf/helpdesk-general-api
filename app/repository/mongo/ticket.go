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

func generateQueryFilterTicket(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	customQuery := helpers.CustomCommonFilter(options)
	for key, value := range customQuery {
		query[key] = value
	}

	return query, mongoOptions
}

func (r *mongoDBRepo) FetchTicketList(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {

	query, findOptions := generateQueryFilterTicket(options, true)

	cur, err = r.Conn.Collection(r.TicketCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchTicketList Find:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CountTicket(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterTicket(options, true)

	total, err := r.Conn.Collection(r.TicketCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountTicket", err)
		return 0
	}
	return
}

func (r *mongoDBRepo) FetchOneTicket(ctx context.Context, options map[string]interface{}) (row *model.Ticket, err error) {
	query, _ := generateQueryFilterTicket(options, false)

	err = r.Conn.Collection(r.TicketCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneTicket FindOne:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CreateTicket(ctx context.Context, ticket *model.Ticket) (err error) {
	_, err = r.Conn.Collection(r.TicketCollection).InsertOne(ctx, ticket)
	if err != nil {
		logrus.Error("CreateTicket InsertOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) UpdateTicket(ctx context.Context, ticket *model.Ticket) (err error) {
	_, err = r.Conn.Collection(r.TicketCollection).UpdateOne(ctx, bson.M{"_id": ticket.ID}, bson.M{"$set": ticket})
	if err != nil {
		logrus.Error("UpdateTicket UpdateOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) UpdateTicketPartial(ctx context.Context, id primitive.ObjectID, field map[string]interface{}) (err error) {
	_, err = r.Conn.Collection(r.TicketCollection).UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{
		"$set": field,
	})
	if err != nil {
		logrus.Error("UpdateTicketPartial:", err)
		return err
	}
	return nil
}
