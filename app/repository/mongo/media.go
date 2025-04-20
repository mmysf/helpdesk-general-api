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

func generateQueryFilterMedia(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// filter
	if category, ok := options["category"].(model.MediaCategory); ok {
		query["category"] = category
	}

	return query, mongoOptions
}

func (r *mongoDBRepo) FetchMediaList(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {

	query, findOptions := generateQueryFilterMedia(options, true)

	cur, err = r.Conn.Collection(r.MediaCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchMediaList Find:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CreateMedia(ctx context.Context, Media *model.Media) (err error) {
	_, err = r.Conn.Collection(r.MediaCollection).InsertOne(ctx, Media)
	if err != nil {
		logrus.Error("CreateMedia InsertOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) UpdateManyMediaPartial(ctx context.Context, ids []primitive.ObjectID, field map[string]interface{}) (err error) {
	filter := bson.M{"_id": bson.M{"$in": ids}}

	_, err = r.Conn.Collection(r.MediaCollection).UpdateMany(context.TODO(), filter, bson.M{
		"$set": field,
	})
	if err != nil {
		logrus.Error("UpdateManyMediaPartial:", err)
		return err
	}
	return nil
}

func (r *mongoDBRepo) FetchOneMedia(ctx context.Context, options map[string]interface{}) (row *model.Media, err error) {
	query, _ := generateQueryFilterMedia(options, false)

	err = r.Conn.Collection(r.MediaCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneMedia FindOne:", err)
		return
	}

	return
}

func (r *mongoDBRepo) UpdateMediaPartial(ctx context.Context, ids primitive.ObjectID, field map[string]interface{}) (err error) {
	_, err = r.Conn.Collection(r.MediaCollection).UpdateOne(context.TODO(), bson.M{"_id": ids}, bson.M{
		"$set": field,
	})
	if err != nil {
		logrus.Error("UpdateMediaPartial:", err)
		return err
	}
	return nil
}
