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

func generateQueryFilterCompany(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)

	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	if accessKey, ok := options["accessKey"].(string); ok {
		query["accessKey"] = accessKey
	}

	if subdomain, ok := options["subdomain"].(string); ok {
		query["settings.domain.subdomain"] = subdomain
	}

	if isCustom, ok := options["isCustom"].(bool); ok {
		query["settings.domain.isCustom"] = isCustom
	}

	if fullUrl, ok := options["fullUrl"].(string); ok {
		query["settings.domain.fullUrl"] = fullUrl
	}

	if email, ok := options["email"].(string); ok {
		query["settings.email"] = email
	}

	if name, ok := options["name"].(string); ok {
		query["name"] = name
	}

	if search, ok := options["search"].(string); ok {
		regex := bson.M{
			"$regex": primitive.Regex{
				Pattern: search,
				Options: "i",
			},
		}
		query["name"] = regex
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

func (r *mongoDBRepo) FetchOneCompany(ctx context.Context, options map[string]interface{}) (row *model.Company, err error) {
	query, _ := generateQueryFilterCompany(options, false)

	err = r.Conn.Collection(r.CompanyCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneCompany FindOne:", err)
		return
	}

	return
}

func (r *mongoDBRepo) FetchCompanyList(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {
	query, mongoOptions := generateQueryFilterCompany(options, true)

	cur, err = r.Conn.Collection(r.CompanyCollection).Find(ctx, query, mongoOptions)
	if err != nil {
		logrus.Error("FetchCompanyList Find:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CountCompany(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterCompany(options, false)

	total, err := r.Conn.Collection(r.CompanyCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountCompany:", err)
	}

	return
}

func (r *mongoDBRepo) CreateCompany(ctx context.Context, company *model.Company) (err error) {
	_, err = r.Conn.Collection(r.CompanyCollection).InsertOne(ctx, company)
	if err != nil {
		logrus.Error("CreateCompany:", err)
	}

	return
}

func (r *mongoDBRepo) UpdateCompany(ctx context.Context, company *model.Company) (err error) {
	_, err = r.Conn.Collection(r.CompanyCollection).UpdateOne(ctx, bson.M{"_id": company.ID}, bson.M{"$set": company})
	if err != nil {
		logrus.Error("UpdateCompany UpdateOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) UpdatePartialCompany(ctx context.Context, options map[string]interface{}, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterCompany(options, false)
	_, err = r.Conn.Collection(r.CompanyCollection).UpdateOne(context.TODO(), query, bson.M{
		"$set": field,
	})

	if err != nil {
		logrus.Error("UpdateCompany UpdateOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) IncrementOneCompany(ctx context.Context, id string, payload map[string]int64) (err error) {
	obj, _ := primitive.ObjectIDFromHex(id)
	_, err = r.Conn.Collection(r.CompanyCollection).UpdateOne(context.Background(), map[string]any{
		"_id": obj,
	}, bson.M{
		"$inc": payload,
	})
	if err != nil {
		logrus.Error("IncrementOneCompany UpdateOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) DeleteCompany(ctx context.Context, company *model.Company) (err error) {
	_, err = r.Conn.Collection(r.CompanyCollection).UpdateOne(ctx, bson.M{"_id": company.ID}, bson.M{"$set": bson.M{"deletedAt": time.Now()}})
	if err != nil {
		logrus.Error("DeleteCompany DeleteOne:", err)
		return
	}
	return
}
