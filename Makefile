container_redis=trash-hunt-redis
container_postgres=trash-hunt-pg
container_tcp_server=trash-hunt-tcp-server
container_api_server=trash-hunt-api-server

container_redis_dev=trash-hunt-redis-dev
container_postgres_dev=trash-hunt-pg-dev
container_tcp_server_dev=trash-hunt-tcp-server-dev
container_api_server_dev=trash-hunt-api-server-dev

db_name_dev=trash_hunt
db_pass_dev=dev_secret
db_user_dev=dev

build_path=.
image_path=./env/container

image_redis=redis
image_postgres=trash-hunt-image-pg
image_tcp_serve=trash-hunt-image-tcp
image_api_servev=trash-hunt-image-api

image_redis_dev=redis
image_postgres_dev=trash-hunt-image-pg-dev
image_tcp_server_dev=trash-hunt-image-tcp-dev
image_api_server_dev=trash-hunt-image-api-dev

.PHONY: help dev test
.DEFAULT_GOAL := help

dev: dev-run-db dev-run-go ## Get a dev docker environment up and running

dev-run-db:
	docker run -p 5432:5432 --name $(container_postgres_dev) -d $(image_postgres_dev)
	docker run -p 6379:6379 --name $(container_redis_dev) -d $(image_redis_dev)

dev-run-go:
	$(eval $@_DB_HOST := $(shell docker inspect --format '{{ .NetworkSettings.IPAddress }}' $(container_postgres_dev)))
	$(eval $@_REDIS_HOST := $(shell docker inspect --format '{{ .NetworkSettings.IPAddress }}' $(container_redis_dev)))
	$(eval $@_API_ENV := -e DB_HOST=$($@_DB_HOST) -e DB_NAME=$(db_name_dev) -e DB_PASS=$(db_pass_dev) -e DB_USER=$(db_user_dev) -e DB_SSL_MODE=disable -e REDIS_HOST=$($@_REDIS_HOST))
	docker run -p 3001:8080 $($@_API_ENV) -v $(shell pwd)/src:/go/src/github.com/tylergeery/trash_hunt/src/ --name $(container_tcp_server_dev) -d $(image_tcp_server_dev)
	docker run -p 3000:8080 $($@_API_ENV) -v $(shell pwd)/src:/go/src/github.com/tylergeery/trash_hunt/src/ --name $(container_api_server_dev) -d $(image_api_server_dev)

dev-clean: ## Remove local docker images
	- docker rmi $(image_api_server_dev) $(image_tcp_server_dev) $(image_postgres_dev)

dev-images: ## Make dev docker images
	docker build -f $(image_path)/pg/Dockerfile --target dev -t $(image_postgres_dev) $(build_path)
	docker build -f $(image_path)/go/Dockerfile --target tcp_server_dev -t $(image_tcp_server_dev) $(build_path)
	docker build -f $(image_path)/go/Dockerfile --target http_server_dev -t $(image_api_server_dev) $(build_path)

dev-kill:
	- docker kill $(container_api_server_dev) $(container_tcp_server_dev) $(container_postgres_dev) $(container_redis_dev)

dev-rm: dev-kill ## Tear down local dev environment
	- docker rm $(container_api_server_dev) $(container_tcp_server_dev) $(container_postgres_dev) $(container_redis_dev)

dev-pg: ## Exec into local pg instance
	docker exec -it $(container_postgres_dev) psql -U dev -W dev_secret dev_secret

test: ## Run tests with local docker env
	$(eval $@_DB_HOST := $(shell docker inspect --format '{{ .NetworkSettings.IPAddress }}' $(container_postgres_dev)))
	$(eval $@_REDIS_HOST := $(shell docker inspect --format '{{ .NetworkSettings.IPAddress }}' $(container_redis_dev)))
	$(eval $@_API_ENV := -e DB_HOST=$($@_DB_HOST) -e DB_NAME=$(db_name_dev) -e DB_PASS=$(db_pass_dev) -e DB_USER=$(db_user_dev) -e DB_SSL_MODE=disable -e REDIS_HOST=$($@_REDIS_HOST))
	docker exec $($@_API_ENV) -it $(container_api_server_dev) /bin/bash -c "go test ./..."

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
