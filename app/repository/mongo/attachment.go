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

func generateQueryFilterAttachment(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	if companyID, ok := options["company_id"].(string); ok {
		query["company.id"] = companyID
	}

	return query, mongoOptions
}

func (r *mongoDBRepo) FetchAttachmentList(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {

	query, findOptions := generateQueryFilterAttachment(options, true)

	cur, err = r.Conn.Collection(r.AttachmentCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchAttachmentList Find:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CreateAttachment(ctx context.Context, Attachment *model.Attachment) (err error) {
	_, err = r.Conn.Collection(r.AttachmentCollection).InsertOne(ctx, Attachment)
	if err != nil {
		logrus.Error("CreateAttachment InsertOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) UpdateManyAttachmentPartial(ctx context.Context, ids []primitive.ObjectID, field map[string]interface{}) (err error) {
	filter := bson.M{"_id": bson.M{"$in": ids}}

	_, err = r.Conn.Collection(r.AttachmentCollection).UpdateMany(context.TODO(), filter, bson.M{
		"$set": field,
	})
	if err != nil {
		logrus.Error("UpdateManyAttachmentPartial:", err)
		return err
	}
	return nil
}

func (r *mongoDBRepo) FetchOneAttachment(ctx context.Context, options map[string]interface{}) (row *model.Attachment, err error) {
	query, _ := generateQueryFilterAttachment(options, false)

	err = r.Conn.Collection(r.AttachmentCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneAttachment FindOne:", err)
		return
	}

	return
}

func (r *mongoDBRepo) UpdateAttachmentPartial(ctx context.Context, ids primitive.ObjectID, field map[string]interface{}) (err error) {
	_, err = r.Conn.Collection(r.AttachmentCollection).UpdateOne(context.TODO(), bson.M{"_id": ids}, bson.M{
		"$set": field,
	})
	if err != nil {
		logrus.Error("UpdateAttachmentPartial:", err)
		return err
	}
	return nil
}
