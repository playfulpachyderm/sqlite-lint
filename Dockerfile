from golang:alpine as builder

run apk add sqlite-dev build-base

copy . /code
workdir /code

env CGO_ENABLED=1
run go build -ldflags="-w -s -linkmode=external -extldflags=-static" -o sqlite_lint ./cmd/sqlite_lint/main.go

# ---

from alpine:3.20

COPY --from=builder /code/sqlite_lint /

entrypoint ["/sqlite_lint"]
