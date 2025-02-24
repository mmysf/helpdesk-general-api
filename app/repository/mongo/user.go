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

func generateQueryFilterUser(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// your
	if username, ok := options["username"].(string); ok {
		query["username"] = username
	}

	return query, mongoOptions
}

func (r *mongoDBRepo) FetchOneUser(ctx context.Context, options map[string]interface{}) (row *model.User, err error) {
	query, _ := generateQueryFilterUser(options, false)

	err = r.Conn.Collection(r.UserCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneUser FindOne:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CreateUser(ctx context.Context, row *model.User) (err error) {
	_, err = r.Conn.Collection(r.UserCollection).InsertOne(ctx, row)
	if err != nil {
		logrus.Error("CreateUser InsertOne:", err)
		return
	}
	return
}
