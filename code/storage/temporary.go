package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gomodule/redigo/redis"
)

var connection redis.Conn
var tempOnce sync.Once

func ensureRedisConnection() {
	var err error

	tempOnce.Do(func() {
		connection, err = redis.DialURL(
			fmt.Sprintf(
				"redis://%s:%s",
				os.Getenv("REDIS_HOST"),
				os.Getenv("REDIS_PORT"),
			),
		)

		if err != nil {
			log.Fatal("Redis connection error: " + err.Error())
		}
	})
}

// QueryKey will fetch a redis key and return the value as a string
func QueryKey(key string) (string, error) {
	ensureRedisConnection()

	return redis.String(connection.Do("GET", key))
}

// QueryNumericKey will fetch a redis key and return the value as an int64
func QueryNumericKey(key string) (int64, error) {
	ensureRedisConnection()

	return redis.Int64(connection.Do("GET", key))
}

// QueryJSONKey will fetch a redis key and jsonUnmarshal
func QueryJSONKey(key string, v interface{}) error {
	ensureRedisConnection()

	bytes, err := redis.Bytes(connection.Do("GET", key))

	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, v)
}

// QueryHash will return a map of all properties from redis hash
func QueryHash(key string) (map[string]string, error) {
	ensureRedisConnection()

	return redis.StringMap(connection.Do("HGETALL", key))
}

// QueryHashKey will get a single property from redis hash
func QueryHashKey(key, prop string) (string, error) {
	ensureRedisConnection()

	return redis.String(connection.Do("HGET", key, prop))
}

// SetKey in redis
func SetKey(key string, value string) error {
	output, err := redis.String(connection.Do("SET", key, value))

	if err != nil {
		return err
	}

	if output != "OK" {
		return errors.New("Could not set value, unknown error")
	}

	return nil
}

// SetNumericKey in redis
func SetNumericKey(key string, value int64) error {
	output, err := redis.String(connection.Do("SET", key, value))

	if err != nil {
		return err
	}

	if output != "OK" {
		return errors.New("Could not set numeric value, unknown error")
	}

	return nil
}

// SetHash - Set redis hash
func SetHash(key string, value map[string]interface{}) error {
	args := []interface{}{key}
	for key, val := range value {
		args = append(args, key, val)
	}
	output, err := redis.String(connection.Do("HMSET", args...))

	if err != nil {
		return err
	}

	if output != "OK" {
		return errors.New("Could not set hash value, unknown error")
	}

	return nil
}

// SetHashKey - Set Redis Hash
func SetHashKey(key, prop, value string) error {
	output, err := redis.Int64(connection.Do("HSET", key, prop, value))

	if err != nil {
		return err
	}

	if output != 0 {
		return errors.New("Could not set hash value, unknown error")
	}

	return nil
}
