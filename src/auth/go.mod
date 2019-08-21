module github.com/tylergeery/trash_hunt/auth

go 1.12

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/tylergeery/trash_hunt/model v0.0.0
	github.com/tylergeery/trash_hunt/storage v0.0.0
	github.com/tylergeery/trash_hunt/test v0.0.0
)

replace github.com/tylergeery/trash_hunt/model => ../model

replace github.com/tylergeery/trash_hunt/storage => ../storage

replace github.com/tylergeery/trash_hunt/test => ../test
