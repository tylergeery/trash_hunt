module github.com/tylergeery/trash_hunt/storage

go 1.12

require (
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/lib/pq v1.2.0
	github.com/tylergeery/trash_hunt/test v0.0.0
)

replace github.com/tylergeery/trash_hunt/test => ../test
