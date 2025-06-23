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

func generateQueryFilterHourPackage(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	if price, ok := options["price"].(float64); ok {
		query["price"] = price
	}

	if name, ok := options["name"].(string); ok {
		query["name"] = name
	}

	if status, ok := options["status"].(string); ok {
		query["status"] = status
	}

	if q, ok := options["q"].(string); ok {
		regex := bson.M{
			"$regex": primitive.Regex{
				Pattern: q,
				Options: "i",
			},
		}
		query["name"] = regex
	}

	return query, mongoOptions
}

func (r *mongoDBRepo) FetchHourPackageList(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {

	query, findOptions := generateQueryFilterHourPackage(options, true)

	cur, err = r.Conn.Collection(r.HourPackageCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchHourPackageList Find:", err)
		return
	}

	return
}

func (r *mongoDBRepo) FetchOneHourPackage(ctx context.Context, options map[string]interface{}) (row *model.HourPackage, err error) {
	query, _ := generateQueryFilterHourPackage(options, false)

	err = r.Conn.Collection(r.HourPackageCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneHourPackage FindOne:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CountHourPackage(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterHourPackage(options, true)

	total, err := r.Conn.Collection(r.HourPackageCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountHourPackage", err)
		return 0
	}
	return
}

func (r *mongoDBRepo) CreateHourPackage(ctx context.Context, hourPackages *model.HourPackage) (err error) {
	_, err = r.Conn.Collection(r.HourPackageCollection).InsertOne(ctx, hourPackages)
	if err != nil {
		logrus.Error("CreateHourPackage InsertOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) UpdateHourPackage(ctx context.Context, hourPackages *model.HourPackage) (err error) {
	_, err = r.Conn.Collection(r.HourPackageCollection).UpdateOne(ctx, bson.M{"_id": hourPackages.ID}, bson.M{"$set": hourPackages})
	if err != nil {
		logrus.Error("UpdateHourPackage UpdateOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) DeleteHourPackage(ctx context.Context, hourPackages *model.HourPackage) (err error) {
	_, err = r.Conn.Collection(r.HourPackageCollection).DeleteOne(ctx, bson.M{"_id": hourPackages.ID})
	if err != nil {
		logrus.Error("DeleteHourPackage DeleteOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) UpdatePartialHourPackage(ctx context.Context, options map[string]interface{}, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterHourPackage(options, false)
	_, err = r.Conn.Collection(r.HourPackageCollection).UpdateOne(context.TODO(), query, bson.M{
		"$set": field,
	})

	if err != nil {
		logrus.Error("UpdateHourPackage UpdateOne:", err)
		return
	}
	return
}
