# Trash Hunt
Multiplayer maze racing game. You're a raccoon and you're hungry. Find the trash can before your opponent.

## Architecture
### HTTP Server
A simple Go HTTP server will exist for game CRUD operations around setup, results, players etc...
See the http_server dir for more information.

### TCP Server
A Go TCP Server will handle socket connections between clients for actual gameplay. This will handle move validation and sending movement events to each of the active clients.
See the socket_server dir for more information.

### Storage
The storage dir handles some wrapper operations around redis/postgres that will be used by the various servers.

## Setup
### Set up docker env
docker run -p 5432:5432 --name trash-hunt-pg -d trash-hunt-pg
docker run -p 6379:6379 --name trash-hunt-redis -d redis

### Connecting to postgres container
docker exec -it trash-hunt-pg psql -U dev -W dev_secret dev_secret
