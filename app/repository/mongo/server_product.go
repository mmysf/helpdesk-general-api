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

func generateQueryFilterServerPackage(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
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

func (r *mongoDBRepo) FetchServerPackageList(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {

	query, findOptions := generateQueryFilterServerPackage(options, true)

	cur, err = r.Conn.Collection(r.ServerPackageCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchServerPackageList Find:", err)
		return
	}

	return
}

func (r *mongoDBRepo) FetchOneServerPackage(ctx context.Context, options map[string]interface{}) (row *model.ServerPackage, err error) {
	query, _ := generateQueryFilterServerPackage(options, false)

	err = r.Conn.Collection(r.ServerPackageCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneServerPackage FindOne:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CountServerPackage(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterServerPackage(options, true)

	total, err := r.Conn.Collection(r.ServerPackageCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountServerPackage", err)
		return 0
	}
	return
}

func (r *mongoDBRepo) CreateServerPackage(ctx context.Context, ServerPackages *model.ServerPackage) (err error) {
	_, err = r.Conn.Collection(r.ServerPackageCollection).InsertOne(ctx, ServerPackages)
	if err != nil {
		logrus.Error("CreateServerPackage InsertOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) UpdateServerPackage(ctx context.Context, ServerPackages *model.ServerPackage) (err error) {
	_, err = r.Conn.Collection(r.ServerPackageCollection).UpdateOne(ctx, bson.M{"_id": ServerPackages.ID}, bson.M{"$set": ServerPackages})
	if err != nil {
		logrus.Error("UpdateServerPackage UpdateOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) DeleteServerPackage(ctx context.Context, ServerPackages *model.ServerPackage) (err error) {
	_, err = r.Conn.Collection(r.ServerPackageCollection).DeleteOne(ctx, bson.M{"_id": ServerPackages.ID})
	if err != nil {
		logrus.Error("DeleteServerPackage DeleteOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) UpdatePartialServerPackage(ctx context.Context, options map[string]interface{}, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterServerPackage(options, false)
	_, err = r.Conn.Collection(r.ServerPackageCollection).UpdateOne(context.TODO(), query, bson.M{
		"$set": field,
	})

	if err != nil {
		logrus.Error("UpdateServerPackage UpdateOne:", err)
		return
	}
	return
}
