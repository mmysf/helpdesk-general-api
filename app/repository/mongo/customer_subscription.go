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

func generateQueryFilterCustomerSubscription(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// filter
	if customerID, ok := options["customerID"].(string); ok {
		query["customer.id"] = customerID
	}

	if serverPackageId, ok := options["serverPackageId"].(string); ok {
		query["serverPackage.id"] = serverPackageId
	}

	if orderType, ok := options["orderType"].(string); ok {
		query["order.type"] = orderType
	}

	if orderType, ok := options["orderType"].(model.OrderType); ok {
		query["order.type"] = orderType
	}

	if status, ok := options["status"].(model.SubscriptionStatus); ok {
		query["status"] = status
	}

	if status, ok := options["status"].(string); ok {
		query["status"] = status
	}

	customQuery := helpers.CustomCommonFilter(options)
	for key, value := range customQuery {
		query[key] = value
	}

	return query, mongoOptions
}

func (r *mongoDBRepo) FetchCustomerSubscriptionList(ctx context.Context, options map[string]interface{}, withOptions bool) (cur *mongo.Cursor, err error) {

	query, findOptions := generateQueryFilterCustomerSubscription(options, withOptions)

	cur, err = r.Conn.Collection(r.CustomerSubscriptionCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchCustomerSubscriptionList Find:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CountCustomerSubscription(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterCustomerSubscription(options, true)

	total, err := r.Conn.Collection(r.CustomerSubscriptionCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountCustomerSubscription", err)
		return 0
	}
	return
}

func (r *mongoDBRepo) FetchOneCustomerSubscription(ctx context.Context, options map[string]interface{}) (row *model.CustomerSubscription, err error) {
	query, _ := generateQueryFilterCustomerSubscription(options, false)

	err = r.Conn.Collection(r.CustomerSubscriptionCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneCustomerSubscription FindOne:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CreateCustomerSubscription(ctx context.Context, row *model.CustomerSubscription) (err error) {
	_, err = r.Conn.Collection(r.CustomerSubscriptionCollection).InsertOne(ctx, row)
	if err != nil {
		logrus.Error("CreateCustomerSubscription InsertOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) UpdateOneCustomerSubscription(ctx context.Context, customerSubscription *model.CustomerSubscription) (err error) {
	_, err = r.Conn.Collection(r.CustomerSubscriptionCollection).UpdateOne(ctx, bson.M{"_id": customerSubscription.ID}, bson.M{"$set": customerSubscription})
	if err != nil {
		logrus.Error("UpdateOneCustomerSubscription UpdateOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) UpdateManyPartialCustomerSubscription(ctx context.Context, ids []primitive.ObjectID, field map[string]interface{}) (err error) {
	filter := bson.M{"_id": bson.M{"$in": ids}}

	_, err = r.Conn.Collection(r.CustomerSubscriptionCollection).UpdateMany(context.TODO(), filter, bson.M{
		"$set": field,
	})
	if err != nil {
		logrus.Error("UpdateManyPartialCustomerSubscription:", err)
		return err
	}
	return nil
}

func (r *mongoDBRepo) UpdatePartialCustomerSubscription(ctx context.Context, options, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterCustomerSubscription(options, false)
	_, err = r.Conn.Collection(r.CustomerSubscriptionCollection).UpdateOne(context.Background(), query, bson.M{
		"$set": field,
	})
	if err != nil {
		logrus.Error("UpdateOneUser UpdateOne:", err)
		return
	}
	return
}
