FROM golang:1.25-alpine3.22 AS build

WORKDIR /var/backend

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/posts_command_consumer/main.go

FROM alpine:edge AS prod

WORKDIR /app

COPY --from=build /var/backend/main /app/main

ENTRYPOINT ["/app/main"]