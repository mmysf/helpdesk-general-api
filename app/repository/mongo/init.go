package mongorepo

import (
	"app/domain/model"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDBRepo struct {
	Conn           *mongo.Database
	UserCollection string
}

func NewMongodbRepo(Conn *mongo.Database) MongoDBRepo {
	return &mongoDBRepo{
		Conn:           Conn,
		UserCollection: "users",
	}
}

type MongoDBRepo interface {
	FetchOneUser(ctx context.Context, options map[string]interface{}) (*model.User, error)
	CreateUser(ctx context.Context, usermodel *model.User) (err error)
}
