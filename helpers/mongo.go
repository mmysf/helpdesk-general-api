package helpers

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

func CommonFilter(options map[string]any) map[string]any {
	query := map[string]any{
		"deleted_at": bson.M{
			"$eq": nil,
		},
	}

	if id, ok := options["id"].(primitive.ObjectID); ok {
		query["_id"] = id
	} else if id, ok := options["id"].(string); ok {
		obj, _ := primitive.ObjectIDFromHex(id)
		query["_id"] = obj
	}

	if ids, ok := options["ids"].([]primitive.ObjectID); ok {
		query["_id"] = bson.M{
			"$in": ids,
		}
	} else if ids, ok := options["ids"].([]string); ok {
		objIDs := make([]primitive.ObjectID, 0)
		for _, id := range ids {
			if obID, err := primitive.ObjectIDFromHex(strings.TrimSpace(id)); err == nil {
				objIDs = append(objIDs, obID)
			}
		}
		query["_id"] = bson.M{
			"$in": objIDs,
		}
	}

	return query
}

func CommonMongoFindOptions(options map[string]any) *moptions.FindOptions {
	// limit, offset & sort
	mongoOptions := moptions.Find()
	if offset, ok := options["offset"].(int64); ok {
		mongoOptions.SetSkip(offset)
	} else if offset, ok := options["offset"].(int); ok {
		mongoOptions.SetSkip(int64(offset))
	}

	if limit, ok := options["limit"].(int64); ok {
		mongoOptions.SetLimit(limit)
	} else if limit, ok := options["limit"].(int); ok {
		mongoOptions.SetLimit(int64(limit))
	}

	if sortBy, ok := options["sort"].(string); ok {
		sortDir, ok := options["dir"].(string)
		if !ok {
			sortDir = "asc"
		}

		sortQ := bson.D{}
		sortDirMongo := int(1)
		if strings.ToLower(sortDir) == "desc" {
			sortDirMongo = -1
		}
		sortQ = append(sortQ, bson.E{
			Key:   sortBy,
			Value: sortDirMongo,
		})
		mongoOptions.SetSort(sortQ)
	} else if sortBy, ok := options["sort"].(map[string]int); ok {
		sortQ := bson.D{}
		for k, sort := range sortBy {
			sortQ = append(sortQ, bson.E{
				Key:   k,
				Value: sort,
			})
		}
		mongoOptions.SetSort(sortQ)
	}

	if projection, ok := options["projection"].(map[string]int); ok {
		mongoOptions.SetProjection(projection)
	}

	return mongoOptions
}
