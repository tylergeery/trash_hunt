container_redis=trash-hunt-redis
container_postgres=trash-hunt-pg
container_tcp_server=trash-hunt-tcp-server
container_api_server=trash-hunt-api-server

build_path=.

image_path=./env/container
image_go_dev=trash-hunt-image-go-dev
image_redis=redis
image_postgres=trash-hunt-image-pg
image_tcp_server=trash-hunt-image-tcp
image_api_server=trash-hunt-image-api

.PHONY: clean dev kill test

dev:
	docker run -p 5432:5432 --name $(container_postgres) -d $(image_postgres)
	docker run -p 6379:6379 --name $(container_redis) -d $(image_redis)
	docker run -p 3001:8080 -v $(shell pwd):/go/src/github.com/tylergeery/trash_hunt/ --name $(container_tcp_server) -d $(image_tcp_server)
	docker run -p 3000:8080 -v $(shell pwd):/go/src/github.com/tylergeery/trash_hunt/ --name $(container_api_server) -d $(image_api_server)

clean:
	- docker rmi $(image_api_server) $(image_tcp_server) $(image_go_dev) $(image_postgres)

images:
	docker build -f $(image_path)/Dockerfile_go_dev -t $(image_go_dev) $(build_path)
	docker build -f $(image_path)/Dockerfile_postgres -t $(image_postgres) $(build_path)
	docker build -f $(image_path)/Dockerfile_tcp -t $(image_tcp_server) $(build_path)
	docker build -f $(image_path)/Dockerfile_api -t $(image_api_server) $(build_path)

kill:
	- docker kill $(container_api_server) $(container_tcp_server) $(container_postgres) $(container_redis)

remove: kill
	- docker rm $(container_api_server) $(container_tcp_server) $(container_postgres) $(container_redis)

pg:
	docker exec -it $(container_postgres) psql -U dev -W dev_secret dev_secret

test:
	docker exec -it $(container_api_server) go test ./...
