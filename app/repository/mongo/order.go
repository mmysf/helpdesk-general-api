package mongorepo

import (
	"app/domain/model"
	"app/helpers"
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

func generateQueryFilterOrder(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// filter
	// orderNumber
	if orderNumber, ok := options["orderNumber"].(string); ok {
		query["orderNumber"] = orderNumber
	}

	if customerID, ok := options["customerID"].(string); ok {
		query["customer.id"] = customerID
	}

	if packageId, ok := options["packageId"].(string); ok {
		query["package.id"] = packageId
	}

	if status, ok := options["status"].([]string); ok {
		query["status"] = bson.M{
			"$in": status,
		}
	}

	if types, ok := options["types"].(string); ok {
		query["type"] = types
	}

	if q, ok := options["q"].(string); ok {
		regex := bson.M{
			"$regex": primitive.Regex{
				Pattern: q,
				Options: "i",
			},
		}
		query["orderNumber"] = regex
	}

	// created at
	if today, ok := options["today"].(bool); ok {
		if today {
			startOfToday := time.Now().UTC().Truncate(24 * time.Hour)
			endOfToday := startOfToday.Add(24 * time.Hour).Truncate(24 * time.Hour)
			query["createdAt"] = bson.M{
				"$gte": startOfToday,
				"$lte": endOfToday,
			}
		}
	}

	return query, mongoOptions
}

func (r *mongoDBRepo) FetchOrderList(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {

	query, findOptions := generateQueryFilterOrder(options, true)

	cur, err = r.Conn.Collection(r.OrderCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchOrderList Find:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CountOrder(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterOrder(options, true)

	total, err := r.Conn.Collection(r.OrderCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountOrder", err)
		return 0
	}
	return
}

func (r *mongoDBRepo) FetchOneOrder(ctx context.Context, options map[string]interface{}) (row *model.Order, err error) {
	query, _ := generateQueryFilterOrder(options, false)

	err = r.Conn.Collection(r.OrderCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneOrder FindOne:", err)
		return
	}

	return
}

func (r *mongoDBRepo) UpdateOrder(ctx context.Context, order *model.Order) (err error) {
	_, err = r.Conn.Collection(r.OrderCollection).UpdateOne(ctx, bson.M{"_id": order.ID}, bson.M{"$set": order})
	if err != nil {
		logrus.Error("UpdateOrder UpdateOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) CreateOrder(ctx context.Context, row *model.Order) (err error) {
	_, err = r.Conn.Collection(r.OrderCollection).InsertOne(ctx, row)
	if err != nil {
		logrus.Error("CreateOrder InsertOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) UpdateOneOrder(ctx context.Context, order *model.Order) (err error) {
	_, err = r.Conn.Collection(r.OrderCollection).UpdateOne(ctx, bson.M{"_id": order.ID}, bson.M{"$set": order})
	if err != nil {
		logrus.Error("UpdateOneOrder UpdateOne:", err)
		return
	}
	return
}
