# ─────────────────────────────────────────────
# Stage 1: Builder
# ─────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# ── Layer 1: workspace descriptor + ALL go.mod files ──────────────────────────
COPY go.work go.work.sum ./
COPY common/go.mod common/go.sum  ./common/
COPY gateway/go.mod gateway/go.sum ./gateway/
COPY bookings/go.mod               ./bookings/
COPY users/go.mod users/go.sum     ./users/
COPY seats/go.mod                  ./seats/

# ── Layer 2: download deps ────────────────────────────────────────────────────
RUN go work sync && cd users && go mod download

# ── Layer 3: source code ──────────────────────────────────────────────────────
COPY common/ ./common/
COPY users/  ./users/

# ── Layer 4: build ────────────────────────────────────────────────────────────
WORKDIR /app/users
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/bin/users \
    ./cmd

# ─────────────────────────────────────────────
# Stage 2: Runtime
# ─────────────────────────────────────────────
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/bin/users .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER appuser

EXPOSE 50051

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD nc -z localhost 50051 || exit 1

ENTRYPOINT ["./users"]
