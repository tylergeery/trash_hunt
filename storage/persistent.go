package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	// wrapper around sql package for postgres
	_ "github.com/lib/pq"
)

var db *sql.DB
var persistentOnce sync.Once

func ensureConnection() {
	persistentOnce.Do(func() {
		openConnection()
	})
}

func openConnection() {
	var err error

	db, err = sql.Open(
		"postgres",
		fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=%s",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_TABLE"),
			os.Getenv("DB_SSL_MODE"),
		),
	)

	if err != nil {
		log.Fatalf("Postgres connection error: %s", err)
	}
}

func closeConnection() {
	db.Close()
}

// Insert a new DB record
func Insert(table string, insert map[string]interface{}, types map[string]string) (int64, error) {
	var insertID int64

	ensureConnection()

	values := getValues(insert, types)
	query := fmt.Sprintf("INSERT INTO %s (%s, created_at, updated_at) VALUES (%s, NOW(), NOW()) RETURNING id", table, getColumns(insert), valuesStub(len(insert)))
	row := db.QueryRow(query, values...)
	err := row.Scan(&insertID)

	return insertID, err
}

// Update an existing DB record
func Update(table string, update map[string]interface{}, types map[string]string) error {
	ensureConnection()

	values := getValues(update, types)
	query := fmt.Sprintf("UPDATE %s SET (%s, updated_at) = (%s, NOW())", table, getColumns(update), valuesStub(len(update)))
	_, err := db.Exec(query, values...)

	return err
}

// Remove an existing DB record
func Remove(table string, id int64) error {
	ensureConnection()

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", table)
	_, err := db.Exec(query, id)

	return err
}

// FetchRows for an arbitrary db query
func FetchRows(query string, args ...interface{}) (*sql.Rows, error) {
	ensureConnection()

	return db.Query(query, args...)
}

// FetchRow returns a single record
func FetchRow(query string, args ...interface{}) *sql.Row {
	ensureConnection()

	return db.QueryRow(query, args...)
}

func getColumns(m map[string]interface{}) string {
	cols := make([]string, 0, len(m))

	for k := range m {
		cols = append(cols, k)
	}

	return strings.Join(cols, ",")
}

func getValues(values map[string]interface{}, types map[string]string) []interface{} {
	vals := make([]interface{}, 0, len(values))

	for k, v := range values {
		if v == nil {
			vals = append(vals, nil)
			continue
		}

		switch types[k] {
		case "int":
			vals = append(vals, v.(int64))
		default:
			vals = append(vals, v.(string))
		}
	}

	return vals
}

func valuesStub(count int) string {
	var stub string

	for i := 1; i <= count; i++ {
		stub += fmt.Sprintf("$%s, ", strconv.Itoa(i))
	}

	return strings.TrimRight(stub, ", ")
}
