package ads

import (
	"dcard_backend/internal/db"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetCurrentCacheVersion() (int, error) {
	val, err := db.GetRedisClient().Get(db.Ctx, "cache_version").Result()
	if err == redis.Nil {
		return 1, nil
	} else if err != nil {
		return 0, err
	}

	version, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return version, nil
}

func IncrementCacheVersion() error {
	_, err := db.GetRedisClient().Incr(db.Ctx, "cache_version").Result()
	return err
}

func RedisQuery(query Query) ([]Response, error) {

	cacheVersion, err := GetCurrentCacheVersion()
	if err != nil {
		return []Response{}, err
	}
	queryString := QueryToString(query, cacheVersion)
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
	cacheVersion, err := GetCurrentCacheVersion()

	if err != nil {
		return err
	}

	queryString := QueryToString(query, cacheVersion)
	serialized, err := json.Marshal(ads)

	if err != nil {
		return err
	}
	err = db.GetRedisClient().Set(db.Ctx, queryString, serialized, time.Minute).Err()
	return err
}

func QueryToString(query Query, cacheVersion int) string {
	queryString := fmt.Sprintf(
		"v=%d&offset=%d&limit=%d&age=%d&gender=%s&country=%s&platform=%s",
		cacheVersion,
		query.Offset,
		query.Limit,
		query.Age,
		query.Gender,
		query.Country,
		query.Platform,
	)
	return queryString
}
