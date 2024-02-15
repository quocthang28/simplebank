FROM golang:1.21.6-alpine as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main main.go

RUN apk add curl

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine:latest as server

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/app.env .

COPY --from=builder /app/migrate ./migrate

COPY db/migration ./migration

COPY start.sh .

COPY wait-for.sh .

RUN chmod +x start.sh wait-for.sh

EXPOSE 8080

CMD ["/app/main"]

ENTRYPOINT [ "/app/start.sh" ]

