FROM golang:1.22-alpine AS builder

WORKDIR /app

# SQLite requires CGO
ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev

# Kopya
COPY go.mod go.sum* ./
RUN go mod download || true

COPY . .

# Modülleri tidy yap
ENV GOPROXY=https://proxy.golang.org,direct
RUN go mod tidy || go mod tidy

# Swagger setup ve init
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init --parseDependency --parseInternal

# Derle
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"
RUN go build -o golearn main.go

FROM alpine:latest
WORKDIR /root/

# SQLite C kütüphanesi için deps (Gerekebilir)
RUN apk add --no-cache tzdata

COPY --from=builder /app/golearn .
COPY --from=builder /app/docs ./docs

# Portu aç
EXPOSE 8090

CMD ["./golearn"]
