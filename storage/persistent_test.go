package storage

import (
	"testing"

	"github.com/tylergeery/trash_hunt/test"
)

func init() {
	test.SetVars()
}

func TestInsert(t *testing.T) {
	table := "test"
	insert := map[string]interface{}{
		"keyword": "Hello!",
	}
	types := map[string]string{
		"keyword": "string",
	}

	defer func() {
		closeConnection()
	}()

	id, err := Insert(table, insert, types)

	if err != nil {
		t.Fatalf("Unexpected insert error: %s", err)
	}

	if id <= 0 {
		t.Fatalf("Received invalid id: %d", id)
	}
}
