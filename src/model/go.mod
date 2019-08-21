module github.com/tylergeery/trash_hunt/model

go 1.12

require (
	github.com/goware/emailx v0.2.0
	github.com/tylergeery/trash_hunt/storage v0.0.0
	github.com/tylergeery/trash_hunt/test v0.0.0
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
)

replace github.com/tylergeery/trash_hunt/storage => ../storage

replace github.com/tylergeery/trash_hunt/test => ../test
