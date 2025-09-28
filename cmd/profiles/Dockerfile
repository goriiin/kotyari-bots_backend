FROM golang:1.24-alpine AS build

WORKDIR /var/backend

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o migrate ./cmd/profiles/migrate/main.go
RUN go build -o main ./cmd/profiles/main.go

FROM alpine:edge AS migrator
RUN apk add --no-cache postgresql-client
WORKDIR /app
COPY --from=build /var/backend/migrate /app/migrate
COPY --from=build /var/backend/assets/migrations/profiles /app/migrations
COPY ./configs/profiles-local.yaml /app/configs/profiles-local.yaml
COPY .env /app/.env
CMD ["./migrate", "--env=local", "--config=./configs/profiles-local.yaml", "up"]

FROM alpine:edge AS prod

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=build /var/backend/main /app/main
COPY ./configs/profiles-local.yaml /app/configs/profiles-local.yaml
COPY .env /app/.env

RUN chown -R appuser:appgroup /app

USER appuser
EXPOSE 8003
EXPOSE 8002

ENTRYPOINT ["/app/main", "--env=local", "--config=./configs/profiles-local.yaml"]