package storage

import (
	"database/sql"
	"testing"

	_ "github.com/tylergeery/trash_hunt/src/test"
)

func TestInsertUpdateAndRemove(t *testing.T) {
	var fetchID, rowFetchID int64
	var createdAt, rowCreatedAt, updatedAt, rowUpdatedAt string

	table := "test"
	query := "SELECT id, created_at, updated_at FROM test WHERE id = $1"
	insert := []interface{}{"Hello!"}
	cols := []string{"keyword"}

	defer func() {
		closeConnection()
	}()

	id, err := Insert(table, insert, cols)

	if err != nil {
		t.Fatalf("Unexpected insert error: %s", err)
	}

	if id <= 0 {
		t.Fatalf("Received invalid id: %d", id)
	}

	row := FetchRow(query, id)
	row.Scan(&fetchID, &createdAt, &updatedAt)

	if id != fetchID {
		t.Fatalf("Received invalid id on fetch: %d", fetchID)
	}
	if createdAt == "" {
		t.Fatalf("Received invalid created_at from fetch: %s", createdAt)
	}
	if updatedAt == "" {
		t.Fatalf("Received invalid updated_at from fetch: %s", updatedAt)
	}

	update := []interface{}{"Yup:::!"}

	err = Update(table, update, cols, id)

	if err != nil {
		t.Fatalf("Unexpected update error: %s", err)
	}

	rows, err := FetchRows(query, id)

	if err != nil {
		t.Fatalf("Unexpected FetchRows error: %s", err)
	}

	defer rows.Close()

	rows.Next()
	rows.Scan(&rowFetchID, &rowCreatedAt, &rowUpdatedAt)

	if id != rowFetchID {
		t.Fatalf("Received invalid id on fetchRows: %d", rowFetchID)
	}
	if createdAt != rowCreatedAt {
		t.Fatalf("Received invalid created_at from fetchRows: %s", rowCreatedAt)
	}
	if updatedAt == rowUpdatedAt || rowUpdatedAt == "" {
		t.Fatalf("Received invalid updated_at from fetchRows: %s", rowUpdatedAt)
	}

	err = Remove(table, id)

	if err != nil {
		t.Fatalf("Unexpected Remove error: %s", err)
	}

	row = FetchRow(query, id)
	err = row.Scan(&fetchID, &createdAt, &updatedAt)

	if err != sql.ErrNoRows {
		t.Fatalf("Expected no rows error")
	}
}
