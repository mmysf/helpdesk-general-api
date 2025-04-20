package mongorepo

import (
	"app/domain/model"
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *mongoDBRepo) FetchOneConfig(ctx context.Context, options map[string]interface{}) (row *model.Config, err error) {
	err = r.Conn.Collection(r.ConfigCollection).FindOne(ctx, bson.M{}).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneConfig FindOne:", err)
		return
	}

	return
}
