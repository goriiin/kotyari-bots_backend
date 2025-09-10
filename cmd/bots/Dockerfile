WORKDIR /var/backend

COPY cmd/profiles/main.go main.go
COPY . .

RUN go mod tidy
RUN go build -o main main.go

FROM alpine:edge as prod

RUN apk add bash

COPY --from=build /var/backend/main /app/main
COPY --from=build /var/backend/.env /app/.env


RUN mkdir -p /app/configs
COPY configs app/configs

WORKDIR /app
EXPOSE 8002

ENTRYPOINT exec ./main --env=local --config="./configs/bots-local-config.yaml"