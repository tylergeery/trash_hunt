# Trash Hunt
Multiplayer maze racing game. You're a raccoon and you're hungry. Find the trash can before your opponent.

## Architecture
### HTTP Server
A simple Go HTTP server will exist for game CRUD operations around setup, results, players etc...
See the http_server directory for more information.

### TCP Server
A Go TCP Server will handle socket connections between clients for actual game play. This will handle move validation and sending movement events to each of the active clients.
See the socket_server directory for more information.

### Storage
The storage directory handles some wrapper operations around databases that will be used by the various servers.

## Setup
### Set up docker env
```bash
make images
make dev
```
The dev server should then be available at localhost:3000

### Connecting to postgres container
```bash
make pg
```
