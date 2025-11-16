FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY common/ ./common/
COPY auth/pkg/ ./auth/pkg/
COPY chat_server/ ./chat_server/

RUN cd chat_server && go mod download

RUN cd chat_server && CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server/main.go

FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/server .

COPY chat_server/.env .env

EXPOSE 50052

CMD ["./server"]
