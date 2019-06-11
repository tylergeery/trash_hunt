FROM golang:latest as base

ENV TOKEN_SIGNING_VALUE 'hello_world'
ENV GOPATH '/go'

RUN go get -u github.com/golang/dep/cmd/dep

COPY ./ /go/src/github.com/tylergeery/trash_hunt



FROM base as http_server

WORKDIR /go/src/github.com/tylergeery/trash_hunt/api_server

RUN dep ensure -vendor-only

EXPOSE 8080



FROM http_server as http_server_dev

RUN go get -u github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -o /go/bin/http ." -command="/go/bin/http"



FROM base as tcp_server

WORKDIR /go/src/github.com/tylergeery/trash_hunt/tcp_server

RUN dep ensure -vendor-only

EXPOSE 8080



FROM tcp_server as tcp_server_dev

RUN go get -u github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -o /go/bin/tcp ." -command="/go/bin/tcp"
