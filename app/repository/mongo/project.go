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

func generateQueryFilterProject(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// filter
	if companyID, ok := options["companyID"].(string); ok {
		query["company.id"] = companyID
	}

	if companyProductID, ok := options["companyProductID"].(string); ok {
		query["companyProduct.id"] = companyProductID
	}

	if createdBy, ok := options["createdBy"].(string); ok {
		query["createdBy.id"] = createdBy
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

func (r *mongoDBRepo) FetchProjectList(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {

	query, findOptions := generateQueryFilterProject(options, true)

	cur, err = r.Conn.Collection(r.ProjectCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchProjectList Find:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CountProject(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterProject(options, true)

	total, err := r.Conn.Collection(r.ProjectCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountProject", err)
		return 0
	}
	return
}

func (r *mongoDBRepo) FetchOneProject(ctx context.Context, options map[string]interface{}) (row *model.Project, err error) {
	query, _ := generateQueryFilterProject(options, false)

	err = r.Conn.Collection(r.ProjectCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneProject FindOne:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CreateProject(ctx context.Context, row *model.Project) (err error) {
	_, err = r.Conn.Collection(r.ProjectCollection).InsertOne(ctx, row)
	if err != nil {
		logrus.Error("CreateProject InsertOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) UpdateOneProject(ctx context.Context, project *model.Project) (err error) {
	_, err = r.Conn.Collection(r.ProjectCollection).UpdateOne(ctx, bson.M{"_id": project.ID}, bson.M{"$set": project})
	if err != nil {
		logrus.Error("UpdateOneProject UpdateOne:", err)
		return
	}
	return
}
