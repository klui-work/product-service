FROM golang:1.25.1-alpine3.22 as Develop
WORKDIR /app

RUN apk update && \
apk add --no-cache tzdata

ENV TZ Asia/Bangkok

RUN go install github.com/air-verse/air@latest

COPY ./code/go.mod ./code/go.sum ./
RUN go mod download && go mod verify

CMD ["air", "-c", ".air.toml"]
EXPOSE 8080

