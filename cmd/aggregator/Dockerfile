FROM golang:1.25-alpine3.22 AS build

WORKDIR /var/backend

COPY cmd/aggregator/main.go main.go
COPY . .

#RUN apk add git
#RUN --mount=type=cache,target=/go/pkg/mod/ \
#    --mount=type=bind,source=go.sum,target=go.sum \
#    --mount=type=bind,source=go.mod,target=go.mod \
#    go mod download -x


RUN go mod tidy
RUN go build -o main main.go

FROM alpine:edge as prod

RUN apk add bash

COPY --from=build /var/backend/main /app/main
COPY --from=build /var/backend/.env /app/.env


RUN mkdir -p /app/internal/configs
COPY configs app/configs

WORKDIR /app
EXPOSE 8000

ENTRYPOINT exec ./main --env=local --config="./configs/docker-config.yaml"
