# ─────────────────────────────────────────────
# Stage 1: Builder
# ─────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# ── Layer 1: workspace descriptor + ALL go.mod files ──────────────────────────
# Copying every go.mod lets the workspace resolver find all modules without
# needing the actual source code yet. This keeps dependency downloads cached
# even when source changes.
COPY go.work go.work.sum ./
COPY common/go.mod common/go.sum  ./common/
COPY gateway/go.mod gateway/go.sum ./gateway/
COPY bookings/go.mod               ./bookings/
COPY users/go.mod users/go.sum     ./users/
COPY seats/go.mod                  ./seats/

# ── Layer 2: download deps (cached unless a go.mod/go.sum changes) ────────────
RUN go work sync && go mod download -x 2>/dev/null || true
RUN cd gateway && go mod download

# ── Layer 3: source code (invalidates layers below only on source changes) ────
COPY common/  ./common/
COPY gateway/ ./gateway/

# ── Layer 4: build ────────────────────────────────────────────────────────────
WORKDIR /app/gateway
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/bin/gateway \
    ./main.go

# ─────────────────────────────────────────────
# Stage 2: Runtime
# ─────────────────────────────────────────────
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/bin/gateway .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER appuser

EXPOSE 3333

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:3333/health || exit 1

ENTRYPOINT ["./gateway"]
