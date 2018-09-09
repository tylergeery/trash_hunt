package storage

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gomodule/redigo/redis"
)

var connection redis.Conn

func init() {
	var err error
	connection, err = redis.DialURL(
        fmt.Sprintf(
            "redis://%s:%s@%s:%s/%s",
            os.Getenv("REDIS_USER"),
            os.Getenv("REDIS_SECRET"),
            os.Getenv("REDIS_HOST"),
            os.Getenv("REDIS_PORT"),
            os.Getenv("REDIS_DB_NUMBER")
        ),
    )

	if err != nil {
		log.Fatal("Redis connection error: " + err.Error())
	}
}

// QueryKey will fetch a redis key and return the value as a string
func QueryKey(key string) (string, error) {
	return redis.String(connection.Do("GET", key))
}

// QueryNumericKey will fetch a redis key and return the value as an int64
func QueryNumericKey(key string) (int64, error) {
	return redis.Int64(connection.Do("GET", key))
}

// QueryJSONKey will fetch a redis key and jsonUnmarshal
func QueryJSONKey(key string, v interface{}) error {
	bytes, err := redis.Bytes(connection.Do("GET", key))

	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, v)
}

// QueryHash will return a map of all properties from redis hash
func QueryHash(key string) (map[string]string, error) {
	return redis.StringMap(connection.Do("HGETALL", key))
}

// QueryHashKey will get a single property from redis hash
func QueryHashKey(key, prop string) (string, error) {
	return redis.String(connection.Do("HGET", key, prop))
}
