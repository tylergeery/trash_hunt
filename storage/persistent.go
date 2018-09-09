package storage

import (
    "database/sql"
    _ "github.com/lib/pq"
)

func init() {
    var err error
    connection, err = sql.Open(
        "postgres",
        fmt.Sprintf(
            "postgres://%s:%s@%s/%s?sslmode=%s",
			os.Getenv("DB_USER"),
            os.Getenv("DB_PASS"),
            os.Getenv("DB_HOST"),
            os.Getenv("DB_TABLE"),
            os.Getenv("DB_SSL_MODE")
        )
    )

	if err != nil {
		log.Fatal("Postgres connection error" + err.Error())
	}
}

func Insert(table string, insert map[string]interface{}, types map[string][string]) (int64, error) {
    var insertId int64

    values := getValues(insert, types)
    query := fmt.Sprintf("INSERT INTO `%s` (%s, created_at, updated_at) VALUES (%s, NOW(), NOW()) RETURNING id", table, getColumns(insert), valuesStub(len(insert)))
    err := connection.QueryRow(query, values...).Scan(&insertId)

    return insertId, err
}

func Update(table string, update map[string]interface{}, types map[string][string]) error {
    values := getValues(update, types)
    query := fmt.Sprintf("UPDATE `%s` SET (%s, updated_at) = (%s, NOW())", table, getColumns(update), valuesStub(len(update)))
    err := connection.QueryRow(query, values...)

    return err
}

func Remove(table string, id int64) error {
    query := fmt.Sprintf("DELETE FROM `%s` WHERE id (%s, updated_at) = (%s, NOW())", table, getColumns(update), valuesStub(len(update)))
    err := connection.QueryRow(query, id)

    return err
}

func getColumns(m map[string]interface []string {
    cols := make(string[], 0, len(m))

    for k := range m {
        cols = append(cols, k)
    }

    return cols
}

func getValues(values map[string]interface{}, types map[string]string) {
    vals := make(string[], 0, len(values))

    for k, v := range m {
        if v == nil {
            vals = append(vals, nil)
            continue
        }

        switch (types[k]) {
        case "int":
            vals = append(vals, v.(int64))
        default:
            vals = append(vals, v.(string))
        }
    }

    return vals
}

func valuesStub(count int) {
	var stub string

	for i := 1; i <= length; i++ {
		stub += fmt.Sprintf("$%s, ", strconv.Itoa(i))
	}

	return strings.TrimRight(stub, ", ")
}
}
