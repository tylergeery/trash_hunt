#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

cd "${DIR}/api_server" && go test ./...
cd "${DIR}/auth" && go test ./...
cd "${DIR}/game" && go test ./...
cd "${DIR}/storage" && go test ./...