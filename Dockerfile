FROM golang:1.19-alpine3.18 as buildbase

WORKDIR /go/src/github.com/rum-people-preseed/telegram-weather-svc

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/telegram-weather-svc main.go

FROM alpine:3.18

COPY --from=buildbase /usr/local/bin/telegram-weather-svc /usr/local/bin/telegram-weather-svc
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["telegram-weather-svc"]