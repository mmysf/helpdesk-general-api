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

func generateQueryFilterSuperuser(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
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

func (r *mongoDBRepo) FetchOneSuperuser(ctx context.Context, options map[string]interface{}) (row *model.Superuser, err error) {
	query, _ := generateQueryFilterSuperuser(options, false)

	err = r.Conn.Collection(r.SuperuserCollection).FindOne(ctx, query).Decode(&row)
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
