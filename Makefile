container_redis=trash-hunt-redis
container_postgres=trash-hunt-pg
container_socket_server=trash-hunt-socker-server
container_web_server=trash-hunt-web-server

image_path=./env/container
image_redis=redis
image_postgres=trash-hunt-image-pg
image_socket_server=trash-hunt-image-tcp
image_web_server=trash-hunt-image-api

dev: kill
	docker run -p 5432:5432 --name $(container_postgres) -d $(image_postgres)
	docker run -p 6379:6379 --name $(container_redis) -d $(image_redis)
	docker run -p 3001:8080 --name $(container_socket_server) -d $(image_socket_server)
	docker run -p 3000:8080 --name $(container_web_server) -d $(image_socket_server)

clean:
	- docker rmi $(image_web_server)
	- docker rmi $(image_socket_server)
	- docker rmi $(image_postgres)

images:
	docker build -f $(image_path)/Dockerfile_postgres -t $(image_postgres) $(image_path)
	docker build -f $(image_path)/Dockerfile_tcp -t $(image_socket_server) $(image_path)
	docker build -f $(image_path)/Dockerfile_api -t $(image_web_server) $(image_path)

kill:
	- docker kill $(container_web_server)
	- docker kill $(container_socket_server)
	- docker kill $(container_postgres)
	- docker kill $(container_redis)

remove: kill
	- docker rm $(container_web_server)
	- docker rm $(container_socket_server)
	- docker rm $(container_postgres)
	- docker rm $(container_redis)

pg:
	docker exec -it $(container_postgres) psql -U dev -W dev_secret dev_secret
