module github.com/tylergeery/trash_hunt/tcp_server

go 1.12

require (
	github.com/tylergeery/trash_hunt/auth v0.0.0
	github.com/tylergeery/trash_hunt/game v0.0.0
	github.com/tylergeery/trash_hunt/storage v0.0.0
	github.com/tylergeery/trash_hunt/test v0.0.0
)

replace github.com/tylergeery/trash_hunt/auth => ../auth

replace github.com/tylergeery/trash_hunt/game => ../game

replace github.com/tylergeery/trash_hunt/storage => ../storage

replace github.com/tylergeery/trash_hunt/test => ../test
