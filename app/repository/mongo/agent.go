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

func generateQueryFilterAgent(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
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

func (r *mongoDBRepo) FetchOneAgent(ctx context.Context, options map[string]interface{}) (row *model.Agent, err error) {
	query, _ := generateQueryFilterAgent(options, false)

	err = r.Conn.Collection(r.AgentCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneAgent FindOne:", err)
		return
	}

	return
}

func (r *mongoDBRepo) FetchAgentList(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {
	query, findOptions := generateQueryFilterAgent(options, true)

	cur, err = r.Conn.Collection(r.AgentCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchAgentList Find:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CountAgent(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterAgent(options, true)

	total, err := r.Conn.Collection(r.AgentCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountAgent", err)
		return 0
	}
	return
}

func (r *mongoDBRepo) CreateAgent(ctx context.Context, row *model.Agent) (err error) {
	_, err = r.Conn.Collection(r.AgentCollection).InsertOne(ctx, row)
	if err != nil {
		logrus.Error("CreateAgent InsertOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) UpdateAgent(ctx context.Context, row *model.Agent) (err error) {
	_, err = r.Conn.Collection(r.AgentCollection).UpdateOne(ctx, bson.M{"_id": row.ID}, bson.M{"$set": row})
	if err != nil {
		logrus.Error("UpdateAgent UpdateOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) DeleteAgent(ctx context.Context, row *model.Agent) (err error) {
	_, err = r.Conn.Collection(r.AgentCollection).DeleteOne(ctx, bson.M{"_id": row.ID})
	if err != nil {
		logrus.Error("DeleteAgent DeleteOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) UpdateOneAgent(ctx context.Context, options, payload map[string]interface{}) (err error) {
	query, _ := generateQueryFilterAgent(options, false)
	_, err = r.Conn.Collection(r.AgentCollection).UpdateOne(ctx, query, bson.M{
		"$set": payload,
	})
	if err != nil {
		logrus.Error("UpdateOneAgent UpdateOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) UpdatePartialAgent(ctx context.Context, options, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterAgent(options, false)
	_, err = r.Conn.Collection(r.AgentCollection).UpdateOne(context.Background(), query, bson.M{
		"$set": field,
	})
	if err != nil {
		logrus.Error("UpdateOneAgent UpdateOne:", err)
		return
	}
	return
}
