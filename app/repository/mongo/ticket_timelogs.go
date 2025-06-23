package mongorepo

import (
	"app/domain/model"
	"app/helpers"
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

func generateQueryFilterTicketTimelogs(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
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

func generateQueryFilterOneTicketTimelogs(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOneOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOneOptions(options)
	}

	customQuery := helpers.CustomCommonFilter(options)
	for key, value := range customQuery {
		query[key] = value
	}

	return query, mongoOptions
}

func (r *mongoDBRepo) FetchTicketTimelogsList(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {

	query, findOptions := generateQueryFilterTicketTimelogs(options, true)

	cur, err = r.Conn.Collection(r.TicketTimelogsCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchTicketTimelogsList Find:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CountTicketTimelogs(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterTicketTimelogs(options, true)

	total, err := r.Conn.Collection(r.TicketTimelogsCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountTicketTimelogs", err)
		return 0
	}
	return
}

func (r *mongoDBRepo) CreateTicketTimelogs(ctx context.Context, TicketTimelogs *model.TicketTimeLogs) (err error) {
	_, err = r.Conn.Collection(r.TicketTimelogsCollection).InsertOne(ctx, TicketTimelogs)
	if err != nil {
		logrus.Error("CreateTicketTimelogs InsertOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) FetchOneTicketlogs(ctx context.Context, options map[string]interface{}) (row *model.TicketTimeLogs, err error) {
	query, optionOne := generateQueryFilterOneTicketTimelogs(options, true)

	err = r.Conn.Collection(r.TicketTimelogsCollection).FindOne(ctx, query, optionOne).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneTicketlogs FindOne:", err)
		return
	}

	return
}

func (r *mongoDBRepo) UpdateTicketlogs(ctx context.Context, ticket *model.TicketTimeLogs) (err error) {
	_, err = r.Conn.Collection(r.TicketTimelogsCollection).UpdateOne(ctx, bson.M{"_id": ticket.ID}, bson.M{"$set": ticket})
	if err != nil {
		logrus.Error("UpdateTicketlogs UpdateOne:", err)
		return
	}
	return
}
