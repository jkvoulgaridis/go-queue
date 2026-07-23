FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/ ./examples/cmd/...

FROM alpine:3.20 AS runner

RUN adduser -D -g '' appuser

WORKDIR /app
COPY --from=builder /app .

USER appuser

EXPOSE 18018