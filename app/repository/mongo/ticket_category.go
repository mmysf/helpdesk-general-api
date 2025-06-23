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

func generateQueryFilterTicketCategory(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	if name, ok := options["name"].(string); ok {
		query["name"] = name
	}

	if companyID, ok := options["companyID"].(string); ok {
		query["company.id"] = companyID
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

func (r *mongoDBRepo) FetchTicketCategoryList(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {

	query, findOptions := generateQueryFilterTicketCategory(options, true)

	cur, err = r.Conn.Collection(r.TicketCategoryCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchTicketCategoryList Find:", err)
		return
	}

	return
}

func (r *mongoDBRepo) FetchOneTicketCategory(ctx context.Context, options map[string]interface{}) (row *model.TicketCategory, err error) {
	query, _ := generateQueryFilterTicketCategory(options, false)

	err = r.Conn.Collection(r.TicketCategoryCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneTicketCategory FindOne:", err)
		return
	}

	return
}
func (r *mongoDBRepo) CreateTicketCategory(ctx context.Context, ticketsCategory *model.TicketCategory) (err error) {
	_, err = r.Conn.Collection(r.TicketCategoryCollection).InsertOne(ctx, ticketsCategory)
	if err != nil {
		logrus.Error("CreateTicketCategory InsertOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) CountTicketCategory(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterTicketCategory(options, true)

	total, err := r.Conn.Collection(r.TicketCategoryCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountTicketCategory CountDocuments:", err)
		return 0
	}
	return
}

func (r *mongoDBRepo) UpdateOneTicketCategory(ctx context.Context, ticketsCategory *model.TicketCategory) (err error) {
	_, err = r.Conn.Collection(r.TicketCategoryCollection).UpdateOne(ctx, bson.M{"_id": ticketsCategory.ID}, bson.M{"$set": ticketsCategory})
	if err != nil {
		logrus.Error("UpdateTicketCategory UpdateOne:", err)
		return
	}
	return
}
