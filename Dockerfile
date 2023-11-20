FROM golang:1.20 AS build

WORKDIR /trest

COPY go.mod go.sum ./

RUN go mod download && go mod tidy
RUN apt-get update && apt-get install -y sqlite3 libsqlite3-dev


COPY . ./

RUN CGO_ENABLED=1 go build -o ./build/main ./cmd

FROM ubuntu:22.04
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /taskm/build/main /main

EXPOSE 8001
ENTRYPOINT ["/main"]