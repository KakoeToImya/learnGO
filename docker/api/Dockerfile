FROM golang:1.26.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/api ./cmd/api

FROM alpine:3.23

RUN apk --no-cache add ca-certificates

WORKDIR /app

RUN adduser -D -g '' appuser

COPY --from=builder /app/api .

USER appuser

EXPOSE 8080

CMD ["./api"]