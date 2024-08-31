FROM golang:1.23-alpine AS builder

## just be case
RUN mkdir -p /app
WORKDIR /app

RUN mkdir -p /app/ssl

COPY go.mod ./
## COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

ENV RUNNING_IS_DOCKER=true

COPY --from=builder /app/main /app/
COPY --from=builder /app/ssl/server.crt /app/ssl/
COPY --from=builder /app/ssl/server.key /app/ssl/

EXPOSE 443

CMD ["/app/main"]
