FROM golang:1.25-alpine3.22 AS build

WORKDIR /var/backend

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o migrate ./cmd/posts_command_consumer/migrate/main.go
RUN go build -o main ./cmd/posts_command_consumer/main.go


FROM alpine:edge AS migrator
RUN apk add --no-cache postgresql-client
WORKDIR /app
COPY --from=build /var/backend/migrate /app/migrate
COPY --from=build /var/backend/assets/migrations/posts /app/migrations
COPY ./configs/posts-local.yaml /app/configs/posts-local.yaml
COPY .env /app/.env
CMD ["./migrate","--env=local", "--config=./configs/posts-local.yaml", "up"]

FROM alpine:edge AS prod

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=build /var/backend/main /app/main
COPY ./configs/posts-local.yaml /app/configs/posts-local.yaml
COPY .env /app/.env

RUN chown -R appuser:appgroup /app

USER appuser

ENTRYPOINT ["/app/main", "--env=local", "--config=./configs/posts-local.yaml"]