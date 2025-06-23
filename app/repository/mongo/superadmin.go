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

func generateQueryFilterSuperadmin(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// your
	if email, ok := options["email"].(string); ok {
		query["email"] = email
	}

	return query, mongoOptions
}

func (r *mongoDBRepo) FetchSuperadminList(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {

	query, findOptions := generateQueryFilterSuperadmin(options, true)

	cur, err = r.Conn.Collection(r.SuperadminCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchAdminList Find:", err)
		return
	}

	return
}

func (r *mongoDBRepo) FetchOneSuperadmin(ctx context.Context, options map[string]interface{}) (row *model.Superadmin, err error) {
	query, _ := generateQueryFilterSuperadmin(options, false)

	err = r.Conn.Collection(r.SuperadminCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneAdmin FindOne:", err)
		return
	}

	return
}
