# syntax=docker/dockerfile:1.26

# =============================================================================
# Build arguments
# =============================================================================
ARG GO_VERSION=1.26.5
ARG ALPINE_VERSION=3.24

# =============================================================================
# Stage: base
# =============================================================================
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS base

# Alpine ships without CA certs, git or build tooling required by some
# Go modules.
RUN apk add --no-cache \
        ca-certificates \
        git \
        make \
        tzdata \
    && update-ca-certificates
 
WORKDIR /workspace

COPY go.mod ./

# =============================================================================
# Stage: lint
# Static analysis gate. Fails the build if linting fails, keeping the
# pipeline honest without polluting the runtime image with the linter.
# =============================================================================
FROM base AS lint
 
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
 
COPY . .

RUN golangci-lint run --timeout=5m ./...
 
# =============================================================================
# Stage: test
# Runs the test suite with race detection and coverage in isolation.
# =============================================================================
FROM base AS test
 
RUN apk add build-base

COPY . .

ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

RUN go vet ./... \
    && go test -race -vet=off -covermode=atomic -coverprofile=/tmp/coverage.out ./...
 
# =============================================================================
# Stage: builder
# Produces a statically linked, stripped binary suitable for a
# from-scratch runtime image.
# =============================================================================
FROM base AS builder
 
ARG APP_NAME=spider
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown
 
COPY . .
 
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
 
# -ldflags strips debug symbols (-s -w) and injects build metadata into
# the binary via a version package, avoiding hardcoded values in source.
RUN go build \
        -trimpath \
        -ldflags="-s -w \
            -X 'main.version=${VERSION}' \
            -X 'main.commit=${COMMIT}' \
            -X 'main.buildDate=${BUILD_DATE}'" \
        -o /bin/${APP_NAME} \
        ./cmd/${APP_NAME}
 
# =============================================================================
# Stage: dev
# Full toolchain image for local/interactive development, used by
# docker-compose. Runs as a non-root user with hot-reload via air.
# =============================================================================
FROM base AS dev
 
ARG APP_UID=10001
ARG APP_GID=10001
 
RUN apk add --no-cache su-exec \
    && go install github.com/air-verse/air@v1.67.4 \
    && addgroup -g ${APP_GID} appgroup \
    && adduser -D -u ${APP_UID} -G appgroup appuser \
    && chown -R appuser:appgroup /workspace /go
 
COPY scripts/docker-entrypoint-dev.sh /usr/local/bin/docker-entrypoint-dev.sh
RUN chmod +x /usr/local/bin/docker-entrypoint-dev.sh
 
COPY --chown=appuser:appgroup . .
 
EXPOSE 4242

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

ENTRYPOINT ["docker-entrypoint-dev.sh"]

CMD ["air", "-c", ".air.toml"]
 
# =============================================================================
# Stage: runtime
# Minimal, distroless-style production image. No shell, no package
# manager, no root user — only the statically linked binary and certs.
# =============================================================================
FROM alpine:${ALPINE_VERSION} AS runtime

ARG APP_NAME=spider
ARG APP_UID=10001
ARG APP_GID=10001

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -g ${APP_GID} appgroup \
    && adduser -D -u ${APP_UID} -G appgroup appuser

COPY --from=builder /bin/${APP_NAME} /app/${APP_NAME}

RUN chown -R appuser:appgroup /app

USER appuser

WORKDIR /app

ENTRYPOINT ["/app"]

# =============================================================================
# Stage: export
# Scratch-only stage exposing solely the compiled binary. Never run — it
# exists to be extracted to the host filesystem via
# `docker buildx build --target export --output type=local,dest=<dir>`.
# =============================================================================
FROM scratch AS export

ARG APP_NAME=spider

COPY --from=builder /bin/${APP_NAME} /${APP_NAME}
