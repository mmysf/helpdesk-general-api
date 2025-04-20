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

func generateQueryFilterCustomer(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// your
	if email, ok := options["email"].(string); ok {
		query["email"] = email
	}

	if token, ok := options["token"].(string); ok {
		query["token"] = token
	}

	if companyID, ok := options["companyID"].(string); ok {
		query["company.id"] = companyID

	}

	if companyType, ok := options["type"].(string); ok {
		query["company.type"] = companyType
	}

	if customerID, ok := options["customerID"].(string); ok {
		objectID, err := primitive.ObjectIDFromHex(customerID)
		if err != nil {
			return nil, nil
		}
		query["_id"] = bson.M{
			"$ne": objectID,
		}
	}

	if companyProductID, ok := options["companyProductID"].(string); ok {
		query["companyProduct.id"] = companyProductID
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

	if passwordResetToken, ok := options["passwordResetToken"].(string); ok {
		query["passwordResetToken"] = passwordResetToken
	}

	return query, mongoOptions
}

func (r *mongoDBRepo) FetchCustomerList(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {

	query, findOptions := generateQueryFilterCustomer(options, true)

	cur, err = r.Conn.Collection(r.CustomerCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchCustomerList Find:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CountCustomer(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterCustomer(options, true)

	total, err := r.Conn.Collection(r.CustomerCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountCustomer", err)
		return 0
	}
	return
}

func (r *mongoDBRepo) FetchOneCustomer(ctx context.Context, options map[string]interface{}) (row *model.Customer, err error) {
	query, _ := generateQueryFilterCustomer(options, false)

	err = r.Conn.Collection(r.CustomerCollection).FindOne(ctx, query).Decode(&row)
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

func (r *mongoDBRepo) CreateCustomer(ctx context.Context, row *model.Customer) (err error) {
	_, err = r.Conn.Collection(r.CustomerCollection).InsertOne(ctx, row)
	if err != nil {
		logrus.Error("CreateUser InsertOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) UpdateOneCustomer(ctx context.Context, options, payload map[string]interface{}) (err error) {
	query, _ := generateQueryFilterCustomer(options, false)
	_, err = r.Conn.Collection(r.CustomerCollection).UpdateOne(ctx, query, bson.M{
		"$set": payload,
	})
	if err != nil {
		logrus.Error("UpdateOneUser UpdateOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) UpdatePartialCustomer(ctx context.Context, options, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterCustomer(options, false)
	_, err = r.Conn.Collection(r.CustomerCollection).UpdateOne(context.Background(), query, bson.M{
		"$set": field,
	})
	if err != nil {
		logrus.Error("UpdateOneUser UpdateOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) UpdateManyPartialCustomer(ctx context.Context, ids []primitive.ObjectID, field map[string]interface{}) (err error) {
	filter := bson.M{"_id": bson.M{"$in": ids}}

	_, err = r.Conn.Collection(r.CustomerCollection).UpdateMany(context.TODO(), filter, bson.M{
		"$set": field,
	})
	if err != nil {
		logrus.Error("UpdateManyPartialCustomer:", err)
		return err
	}
	return nil
}

func (r *mongoDBRepo) CreateManyCustomer(ctx context.Context, rows []*model.Customer) (err error) {
	docs := make([]interface{}, len(rows))
	for i, row := range rows {
		docs[i] = row
	}
	_, err = r.Conn.Collection(r.CustomerCollection).InsertMany(ctx, docs)
	if err != nil {
		logrus.Error("CreateManyCustomer InsertMany:", err)
		return
	}
	return
}
