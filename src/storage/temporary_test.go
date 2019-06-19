package storage

import (
	"encoding/json"
	"strconv"
	"testing"

	_ "github.com/tylergeery/trash_hunt/src/test"
)

func TestQueryKeyNonexistent(t *testing.T) {
	notExistsKey := "notexistfakekey"
	value, err := QueryKey(notExistsKey)

	if err == nil {
		t.Fatalf("Expected redigo not found error: got value: %s", value)
	}
}

func TestSetAndQueryKey(t *testing.T) {
	key := "temp"
	value := "also temp"

	err := SetKey(key, value)

	if err != nil {
		t.Fatal(err)
	}

	retrieved, err := QueryKey(key)

	if err != nil {
		t.Fatal(err)
	}
	if value != retrieved {
		t.Fatalf("Expected value (%s), received: %s", value, retrieved)
	}
}

func TestSetAndQueryNumericKey(t *testing.T) {
	key := "tempCount"
	var value int64 = 45

	err := SetNumericKey(key, value)

	if err != nil {
		t.Fatal(err)
	}

	retrieved, err := QueryNumericKey(key)

	if err != nil {
		t.Fatal(err)
	}
	if value != retrieved {
		t.Fatalf("Expected value (%d), received: %d", value, retrieved)
	}
}

func TestSetAndQueryJSONKey(t *testing.T) {
	key := "temp_json-key"
	value := map[string]interface{}{
		"name": "Temp",
		"age":  41,
	}
	jsonStr, _ := json.Marshal(value)
	decoded := make(map[string]interface{})

	err := SetKey(key, string(jsonStr))

	if err != nil {
		t.Fatal(err)
	}

	err = QueryJSONKey(key, &decoded)

	if err != nil {
		t.Fatal(err)
	}
	if value["name"] != decoded["name"].(string) {
		t.Fatalf("Expected value (%s), received: %s", value["name"], decoded["name"])
	}
	if value["age"] != int(decoded["age"].(float64)) {
		t.Fatalf("Expected value (%d), received: %d", value["age"], decoded["age"])
	}
}

func TestSetHashAndQuery(t *testing.T) {
	key := "temp hash"
	value := map[string]interface{}{
		"name": "LAT AFNS asd",
		"age":  23,
	}

	err := SetHash(key, value)

	if err != nil {
		t.Fatal(err)
	}

	retrieved, err := QueryHash(key)

	if err != nil {
		t.Fatal(err)
	}
	if value["name"] != retrieved["name"] {
		t.Fatalf("Expected value (%s), received: %s", value["name"], retrieved["name"])
	}
	age, _ := strconv.Atoi(retrieved["age"])
	if value["age"] != int(age) {
		t.Fatalf("Expected value (%d), received: %d", value["age"], age)
	}

	value["age"] = 58
	err = SetHashKey(key, "age", "58")

	if err != nil {
		t.Fatal(err)
	}

	ageStr, err := QueryHashKey(key, "age")

	if err != nil {
		t.Fatal(err)
	}

	age, _ = strconv.Atoi(ageStr)
	if value["age"] != age {
		t.Fatalf("Expected value (%d), received: %d", value["age"], age)
	}
}
