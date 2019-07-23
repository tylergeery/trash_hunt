module github.com/tylergeery/trash_hunt/api_server

go 1.12

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-ozzo/ozzo-routing v2.1.4+incompatible
	github.com/golang/gddo v0.0.0-20190419222130-af0f2af80721 // indirect
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.3.0 // indirect
	github.com/goware/emailx v0.2.0
	github.com/lib/pq v1.2.0
	github.com/stretchr/testify v1.3.0 // indirect
	github.com/tylergeery/trash_hunt/auth v0.0.0
	github.com/tylergeery/trash_hunt/game v0.0.0
	github.com/tylergeery/trash_hunt/storage v0.0.0
	github.com/tylergeery/trash_hunt/test v0.0.0
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
)

replace github.com/tylergeery/trash_hunt/auth => ../auth

replace github.com/tylergeery/trash_hunt/game => ../game

replace github.com/tylergeery/trash_hunt/storage => ../storage

replace github.com/tylergeery/trash_hunt/test => ../test
