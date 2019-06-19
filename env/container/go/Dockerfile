FROM golang:latest as base

ENV TOKEN_SIGNING_VALUE 'hello_world'
ENV GOPATH '/go'

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/tylergeery/trash_hunt/src

COPY ./src .

RUN dep ensure -vendor-only



FROM base as http_server

WORKDIR /go/src/github.com/tylergeery/trash_hunt/src/api_server

EXPOSE 8080



FROM http_server as http_server_dev

RUN go get -u github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -o /go/bin/http ." -command="/go/bin/http"



FROM base as tcp_server

WORKDIR /go/src/github.com/tylergeery/trash_hunt/src/tcp_server

EXPOSE 8080



FROM tcp_server as tcp_server_dev

RUN go get -u github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -o /go/bin/tcp ." -command="/go/bin/tcp"
