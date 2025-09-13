FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/bots-server ./cmd/bots

FROM alpine:3.20 AS prod

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
COPY --from=build /app/bots-server .
COPY configs/ ./configs/

RUN chown -R appuser:appgroup /app
USER appuser
EXPOSE 8001

ENTRYPOINT ["./bots-server"]