package mongorepo

import (
	"app/domain/model"
	"app/helpers"
	"context"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

func generateQueryFilterNotification(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)

	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	if isRead, ok := options["isRead"].(bool); ok {
		query["isRead"] = isRead
	}

	if typ, ok := options["type"].(string); ok {
		query["type"] = typ
	}

	if IsLastTwo, ok := options["isLastTwo"].(bool); ok {
		query["isLastTwo"] = IsLastTwo
	}

	if userRole := options["userRole"]; userRole != nil {
		query["userRole"] = userRole
	}

	if userID, ok := options["userID"].(string); ok {
		query["user.id"] = userID
	}

	if category, ok := options["category"].(string); ok {
		query["category.name"] = category
	}

	return query, mongoOptions
}

func (r *mongoDBRepo) CountNotification(ctx context.Context, options map[string]interface{}) int64 {
	query, _ := generateQueryFilterNotification(options, false)
	count, err := r.Conn.Collection(r.NotificationCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountNotification Find:", err)
		return 0
	}
	return count
}

func (r *mongoDBRepo) CreateNotification(ctx context.Context, row *model.Notification) (err error) {
	if row.Title == "" {
		return
	}
	_, err = r.Conn.Collection(r.NotificationCollection).InsertOne(ctx, row)
	if err != nil {
		logrus.Error("CreateNotification InsertOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) DeleteNotification(ctx context.Context, notification *model.Notification) (err error) {
	_, err = r.Conn.Collection(r.NotificationCollection).UpdateOne(ctx, bson.M{"_id": notification.ID}, bson.M{"$set": bson.M{"deletedAt": time.Now().UTC()}})
	if err != nil {
		logrus.Error("DeleteNotification UpdateOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) UpdateNotification(ctx context.Context, notification *model.Notification) (err error) {
	_, err = r.Conn.Collection(r.NotificationCollection).UpdateOne(ctx, bson.M{"_id": notification.ID}, bson.M{"$set": notification})

	if err != nil {
		logrus.Error("UpdateNotification UpdateOne:", err)
		return
	}
	return
}

func (r *mongoDBRepo) FetchNotificationList(ctx context.Context, options map[string]interface{}) (cursor *mongo.Cursor, err error) {
	query, mongoOptions := generateQueryFilterNotification(options, true)

	jsonBytes, err := json.Marshal(mongoOptions)
	if err != nil {
		logrus.Error("FetchNotificationList Marshal:", err)
		return
	}
	logrus.Info("query: ", string(jsonBytes))

	cursor, err = r.Conn.Collection(r.NotificationCollection).Find(ctx, query, mongoOptions)
	if err != nil {
		logrus.Error("FetchNotificationList Find:", err)
	}

	return cursor, err

}

func (r *mongoDBRepo) FetchOneNotification(ctx context.Context, options map[string]interface{}) (*model.Notification, error) {
	query, _ := generateQueryFilterNotification(options, false)
	var notification *model.Notification
	err := r.Conn.Collection(r.NotificationCollection).FindOne(ctx, query).Decode(&notification)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return nil, err
		}
		return nil, err
	}
	return notification, nil
}

func (r *mongoDBRepo) ReadAllNotification(ctx context.Context, userID string) (err error) {
	_, err = r.Conn.Collection(r.NotificationCollection).UpdateMany(ctx, bson.M{}, bson.M{"$set": bson.M{"isRead": true}})
	if err != nil {
		logrus.Error("ReadAllNotification UpdateMany:", err)
		return
	}
	return
}
