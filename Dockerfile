FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o forum ./cmd

FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache sqlite


COPY --from=builder /app/forum .
COPY --from=builder /app/certs ./certs
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/web ./web
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./forum"]