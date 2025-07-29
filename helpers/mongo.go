package helpers

import (
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

func CommonFilter(options map[string]any) map[string]any {
	query := map[string]any{
		"deletedAt": bson.M{
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

func CustomCommonFilter(options map[string]any) (query map[string]any) {
	query = map[string]any{}
	if companyID, ok := options["companyID"].(string); ok {
		query["company.id"] = companyID
	}

	if customerID, ok := options["customerID"].(string); ok {
		query["customer.id"] = customerID
	}

	if companyProductID, ok := options["companyProductID"].(string); ok {
		query["product.id"] = companyProductID
	}

	if projectID, ok := options["projectID"].(string); ok {
		query["project.id"] = projectID
	}

	if ticketId, ok := options["ticketId"].(string); ok {
		query["ticket.id"] = ticketId
	}

	if companyProductName, ok := options["companyProductName"].(string); ok {
		regex := bson.M{
			"$regex": primitive.Regex{
				Pattern: companyProductName,
				Options: "i",
			},
		}
		query["product.name"] = regex
	}

	if status, ok := options["status"].([]string); ok {
		// query["status"] = status
		query["status"] = bson.M{
			"$in": status,
		}
	}

	if token, ok := options["token"].(string); ok {
		query["token"] = token
	}

	if code, ok := options["code"].(string); ok {
		regex := bson.M{
			"$regex": primitive.Regex{
				Pattern: code,
				Options: "i",
			},
		}
		query["code"] = regex
	}

	if subject, ok := options["subject"].(string); ok {
		escapedSubject := regexp.QuoteMeta(subject)
		regex := bson.M{
			"$regex": primitive.Regex{
				Pattern: escapedSubject,
				Options: "i",
			},
		}
		query["subject"] = regex
	}

	// day
	if day, ok := options["day"].(string); ok {
		query["detailTime.dayName"] = day
	}

	// date
	startDate, startOk := options["startDate"].(time.Time)
	endDate, endOk := options["endDate"].(time.Time)
	if !startOk && endOk {
		startDate = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	if !endOk && startOk {
		endDate = time.Now()
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 17, 0, 0, 0, time.UTC)
	}
	if startOk || endOk {
		if startOk {
			startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day()-1, 17, 0, 0, 0, time.UTC)
		}
		if endOk {
			endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 17, 0, 0, 0, time.UTC)
		}
		query["createdAt"] = bson.M{
			"$gte": startDate,
			"$lte": endDate,
		}
	}

	// date expired
	startExpDate, startExpOk := options["startExpDate"].(time.Time)
	endExpDate, endExpOk := options["endExpDate"].(time.Time)
	if !startExpOk && endExpOk {
		startExpDate = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	if !endExpOk && startExpOk {
		endExpDate = time.Now()
		endExpDate = time.Date(endExpDate.Year(), endExpDate.Month(), endExpDate.Day(), 17, 0, 0, 0, time.UTC)
	}
	if startExpOk || endExpOk {
		// if startExpOk {
		// 	startExpDate = time.Date(startExpDate.Year(), startExpDate.Month(), startExpDate.Day()-1, 17, 0, 0, 0, time.UTC)
		// }
		// if endExpOk {
		// 	endExpDate = time.Date(endExpDate.Year(), endExpDate.Month(), endExpDate.Day(), 17, 0, 0, 0, time.UTC)
		// }
		query["expiredAt"] = bson.M{
			"$gte": startExpDate,
			"$lte": endExpDate,
		}
	}

	if date, ok := options["date"].(time.Time); ok {
		query["detailTime.day"] = date.Day()
	}

	if month, ok := options["month"].(time.Month); ok {
		query["detailTime.month"] = month
	}

	if year, ok := options["year"].(int); ok {
		query["detailTime.year"] = year
	}

	// created at
	if today, ok := options["today"].(bool); ok {
		if today {
			startOfToday := time.Now().UTC().Truncate(24 * time.Hour)
			endOfToday := startOfToday.Add(24 * time.Hour).Truncate(24 * time.Hour)
			query["createdAt"] = bson.M{
				"$gte": startOfToday,
				"$lte": endOfToday,
			}
		}
	}

	// expired at
	if expToday, ok := options["expToday"].(bool); ok {
		if expToday {
			startOfToday := time.Now().UTC().Truncate(24 * time.Hour)
			endOfToday := startOfToday.Add(24 * time.Hour).Truncate(24 * time.Hour)
			query["expiredAt"] = bson.M{
				"$gte": startOfToday,
				"$lte": endOfToday,
			}
		}
	}

	//logTime status
	if logStatus, ok := options["logTime.status"].(string); ok {
		query["logTime.status"] = logStatus
	}

	//logTime endAt
	if logEndTime, ok := options["logTime.endAt"].(time.Time); ok {
		query["logTime.endAt"] = bson.M{
			"$lt": logEndTime,
		}
	}

	//reminderSent
	if reminderSent, ok := options["reminderSent"].(bool); ok {
		query["$or"] = []bson.M{
			{"reminderSent": reminderSent},
			{"reminderSent": bson.M{"$exists": false}},
		}
	}

	// agent
	if agentID, ok := options["agentID"].(string); ok {
		query["agents.id"] = agentID
	}

	// priority
	if priority, ok := options["priority"].(string); ok {
		query["priority"] = priority
	}

	// category
	if categoryID, ok := options["categoryID"].(string); ok {
		query["category.id"] = categoryID
	}

	// last month ticket
	if lastMonth, ok := options["lastMonth"].(bool); ok && lastMonth {
		now := time.Now()
		startOfLastMonth := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.UTC)
		endOfLastMonth := startOfLastMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

		query["createdAt"] = bson.M{
			"$gte": startOfLastMonth,
			"$lte": endOfLastMonth,
		}
	}

	// completed by
	if completedBy, ok := options["completedBy"].(string); ok && completedBy != "" {
		query["completedBy.id"] = completedBy
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

func CommonMongoFindOneOptions(options map[string]any) *moptions.FindOneOptions {
	// limit, offset & sort
	mongoOptions := moptions.FindOne()
	if offset, ok := options["offset"].(int64); ok {
		mongoOptions.SetSkip(offset)
	} else if offset, ok := options["offset"].(int); ok {
		mongoOptions.SetSkip(int64(offset))
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
