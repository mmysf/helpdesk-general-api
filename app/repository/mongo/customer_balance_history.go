package mongorepo

import (
	"app/domain/model"
	"context"

	"github.com/sirupsen/logrus"
)

func (r *mongoDBRepo) CreateCustomerBalanceHistory(ctx context.Context, row *model.CustomerBalanceHistory) (err error) {
	_, err = r.Conn.Collection(r.CustomerBalanceHistoryCollection).InsertOne(ctx, row)
	if err != nil {
		logrus.Error("CreateCustomerBalanceHistory InsertOne:", err)
		return
	}
	return
}
