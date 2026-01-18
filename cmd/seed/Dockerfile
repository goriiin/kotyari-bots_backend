# cmd/seed/Dockerfile
FROM golang:1.25-alpine3.22 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/seed ./cmd/seed

FROM alpine:edge

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=build /out/seed /app/seed
COPY .env /app/.env

ENTRYPOINT ["/app/seed"]
