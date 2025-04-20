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

func generateQueryFilterCompanyProduct(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	if companyID, ok := options["companyID"].(string); ok {
		query["company.id"] = companyID
	}

	if code, ok := options["code"].(string); ok {
		query["code"] = code
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

func (r *mongoDBRepo) FetchOneCompanyProduct(ctx context.Context, options map[string]interface{}) (row *model.CompanyProduct, err error) {
	query, _ := generateQueryFilterCompanyProduct(options, false)

	err = r.Conn.Collection(r.CompanyProductCollection).FindOne(ctx, query).Decode(&row)
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

func (r *mongoDBRepo) FetchCompanyProductList(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {

	query, findOptions := generateQueryFilterCompanyProduct(options, true)

	cur, err = r.Conn.Collection(r.CompanyProductCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchProductList Find:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CountCompanyProduct(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterCompanyProduct(options, true)

	total, err := r.Conn.Collection(r.CompanyProductCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountTicket", err)
		return 0
	}
	return
}

func (r *mongoDBRepo) CreateCompanyProduct(ctx context.Context, companyProduct *model.CompanyProduct) (err error) {
	_, err = r.Conn.Collection(r.CompanyProductCollection).InsertOne(ctx, companyProduct)
	if err != nil {
		logrus.Error("CreateCompanyProduct:", err)
	}

	return
}

func (r *mongoDBRepo) UpdateCompanyProduct(ctx context.Context, companyProduct *model.CompanyProduct) (err error) {
	_, err = r.Conn.Collection(r.CompanyProductCollection).UpdateOne(ctx, bson.M{"_id": companyProduct.ID}, bson.M{"$set": companyProduct})
	if err != nil {
		logrus.Error("UpdateCompanyProduct UpdateOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) UpdatePartialCompanyProduct(ctx context.Context, options map[string]interface{}, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterCompanyProduct(options, false)
	_, err = r.Conn.Collection(r.CompanyProductCollection).UpdateOne(context.TODO(), query, bson.M{
		"$set": field,
	})

	if err != nil {
		logrus.Error("UpdateCompanyProduct UpdateOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) IncrementOneCompanyProduct(ctx context.Context, id string, payload map[string]int64) (err error) {
	obj, _ := primitive.ObjectIDFromHex(id)
	_, err = r.Conn.Collection(r.CompanyProductCollection).UpdateOne(context.Background(), map[string]any{
		"_id": obj,
	}, bson.M{
		"$inc": payload,
	})
	if err != nil {
		logrus.Error("IncrementOneCompanyProduct UpdateOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) DeleteCompanyProduct(ctx context.Context, companyProduct *model.CompanyProduct) (err error) {
	_, err = r.Conn.Collection(r.CompanyProductCollection).UpdateOne(ctx, bson.M{"_id": companyProduct.ID}, bson.M{"$set": bson.M{"deletedAt": time.Now()}})
	if err != nil {
		logrus.Error("DeleteCompanyProduct DeleteOne:", err)
		return
	}
	return
}
