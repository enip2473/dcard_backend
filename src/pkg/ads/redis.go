package ads

import (
	"dcard_backend/internal/db"
	"encoding/json"
	"fmt"
	"time"
)

func RedisQuery(query Query) ([]Response, error) {
	queryString := QueryToString(query)
	val, err := db.GetRedisClient().Get(db.Ctx, queryString).Result()

	if err != nil {
		return []Response{}, err
	}

	var ads []Response
	if err := json.Unmarshal([]byte(val), &ads); err == nil {
		currentTime := time.Now()
		index := len(ads)
		for i, ad := range ads {
			if currentTime.Before(ad.EndAt) {
				index = i
				break
			}
		}
		if index != 0 {
			ads = ads[index:]
			err = db.GetRedisClient().Del(db.Ctx, queryString).Err()
			if err != nil {
				return []Response{}, err
			}
		}
		return ads, nil
	} else {
		return []Response{}, err
	}
}

func RedisInsert(query Query, ads []Response) error {
	queryString := QueryToString(query)
	serialized, err := json.Marshal(ads)

	if err != nil {
		return err
	}
	err = db.GetRedisClient().Set(db.Ctx, queryString, serialized, 0).Err()
	return err
}

func QueryToString(query Query) string {
	queryString := fmt.Sprintf(
		"offset=%d&limit=%d&age=%d&gender=%s&country=%s&platform=%s",
		query.Offset,
		query.Limit,
		query.Age,
		query.Gender,
		query.Country,
		query.Platform,
	)
	return queryString
}
