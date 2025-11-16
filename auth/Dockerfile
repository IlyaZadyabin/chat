FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY common/ ./common/
COPY auth/ ./auth/

RUN cd auth && go mod download

RUN cd auth && CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server/main.go

FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/server .

COPY auth/.env .env

EXPOSE 50051 8081 8082 2112

CMD ["./server"]
