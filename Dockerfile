# Build stage
FROM golang:1.26-alpine AS builder
WORKDIR /app 
COPY . .
RUN go build -o main main.go

RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.1/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY migrations ./migration

EXPOSE 8082
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]