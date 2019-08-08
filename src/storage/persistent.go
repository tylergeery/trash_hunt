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

func init() {
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
			os.Getenv("DB_NAME"),
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

func getColumns(cols []string) string {
	return strings.Join(cols, ",")
}

// Exec executes an aribitrary SQL command
func Exec(command string) (sql.Result, error) {
	return db.Exec(command)
}

// Insert a new DB record
func Insert(table string, values []interface{}, columns []string) (int64, error) {
	var insertID int64

	query := fmt.Sprintf(
		"INSERT INTO %s (%s, created_at, updated_at) VALUES (%s, NOW(), NOW()) RETURNING id",
		table,
		getColumns(columns),
		valuesStub(len(values)),
	)
	row := db.QueryRow(query, values...)
	err := row.Scan(&insertID)

	return insertID, err
}

// Update an existing DB record
func Update(table string, values []interface{}, columns []string, ID int64) error {
	query := fmt.Sprintf(
		"UPDATE %s SET (%s, updated_at) = (%s, NOW()) WHERE id = %d",
		table,
		getColumns(columns),
		valuesStub(len(values)),
		ID,
	)
	_, err := db.Exec(query, values...)

	return err
}

// Remove an existing DB record
func Remove(table string, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", table)
	_, err := db.Exec(query, id)

	return err
}

// FetchRows for an arbitrary db query
func FetchRows(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}

// FetchRow returns a single record
func FetchRow(query string, args ...interface{}) *sql.Row {
	return db.QueryRow(query, args...)
}

func valuesStub(count int) string {
	var stub string

	for i := 1; i <= count; i++ {
		stub += fmt.Sprintf("$%s, ", strconv.Itoa(i))
	}

	return strings.TrimRight(stub, ", ")
}
