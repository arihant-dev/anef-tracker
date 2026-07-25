# Stage 1: Build binary
FROM golang:1.24-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X github.com/arihant-dev/anef-tracker/internal/version.Version=v0.9.0" -o /app/bin/anef ./cmd/anef

# Stage 2: Minimal runtime image
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/bin/anef /usr/local/bin/anef

VOLUME ["/app/data"]
ENTRYPOINT ["/usr/local/bin/anef"]
CMD ["watch", "--interval", "360"]
