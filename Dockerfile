FROM golang:1.23-alpine AS builder

## just be case
RUN mkdir -p /app
WORKDIR /app

RUN mkdir -p /app/keys

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

ENV RUNNING_IS_DOCKER=true

COPY --from=builder /app/main /app/

COPY --from=builder /app/keys/ /app/keys/

# SSLは要らない
COPY --from=builder /app/ssl/server.crt /app/ssl/
COPY --from=builder /app/ssl/server.key /app/ssl/

EXPOSE 8080

CMD ["/app/main"]
