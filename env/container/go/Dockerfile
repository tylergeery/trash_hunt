FROM golang:1.12.7 as base

ENV TOKEN_SIGNING_VALUE 'hello_world'
ENV GOPATH '/go/path'
ENV PATH="/go/path/bin:${PATH}"

WORKDIR /go/src

COPY ./src .



FROM base as api_server

WORKDIR /go/src/api_server
RUN go build -o /go/bin/trash_hunt_api_server

EXPOSE 8080

ENTRYPOINT /go/bin/trash_hunt_api_server



FROM api_server as api_server_dev

RUN go get -u github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -o /go/bin/trash_hunt_api_server ." -command="/go/bin/trash_hunt_api_server"



FROM base as tcp_server

WORKDIR /go/src/tcp_server
RUN go build -o /go/bin/trash_hunt_tcp_server

EXPOSE 8080

ENTRYPOINT /go/bin/trash_hunt_tcp_server



FROM tcp_server as tcp_server_dev

RUN go get -u github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -o /go/bin/trash_hunt_tcp_server ." -command="/go/bin/trash_hunt_tcp_server"
