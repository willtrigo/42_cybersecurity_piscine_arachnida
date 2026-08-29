# syntax=docker/dockerfile:1.26

# =============================================================================
# Build arguments
# =============================================================================
ARG GO_VERSION=1.27.0
ARG DEBIAN_CODENAME=bookworm

# =============================================================================
# Stage: base
# =============================================================================
FROM golang:${GO_VERSION}-${DEBIAN_CODENAME} AS base

# Debian ships without CA certs, git or build tooling required by some
# Go modules.
RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
        git \
        make \
        tzdata \
    && rm -rf /var/lib/apt/lists/* \
    && update-ca-certificates
 
WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download && go mod verify

# =============================================================================
# Stage: gui-build-deps
# Shared X11/OpenGL/Wayland headers required to compile the Fyne GUI
# presenter's cgo/glfw bindings. Factored out so lint, test, builder and
# dev install the exact same set instead of four hand-kept copies.
# =============================================================================
FROM base AS gui-build-deps
 
RUN apt-get update && apt-get install -y --no-install-recommends \
        build-essential \
        pkg-config \
        libgl1-mesa-dev \
        libx11-dev \
        libxcursor-dev \
        libxinerama-dev \
        libxrandr-dev \
        libxi-dev \
        libxext-dev \
        libxfixes-dev \
        libxxf86vm-dev \
        libwayland-dev \
        wayland-protocols \
        libxkbcommon-dev \
    && rm -rf /var/lib/apt/lists/*

# =============================================================================
# Stage: lint
# Static analysis gate. Fails the build if linting fails, keeping the
# pipeline honest without polluting the runtime image with the linter.
# =============================================================================
FROM gui-build-deps AS lint

ARG GOLANGCI_LINT_VERSION=v2.13.1

# Fyne's GUI presenter uses cgo (via glfw/OpenGL bindings), so the
# lint needs a C toolchain and the X11/GL headers to link against.
RUN go install \
        github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_LINT_VERSION} \
    && golangci-lint version

COPY . .

RUN golangci-lint run --timeout=10m ./...
 
# =============================================================================
# Stage: test
# Runs the test suite with race detection and coverage in isolation.
# =============================================================================
FROM gui-build-deps AS test
 
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
FROM gui-build-deps AS builder
 
ARG APP_NAME=scorpion
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown
 
COPY . .
 
ENV CGO_ENABLED=1 \
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
FROM gui-build-deps AS dev
 
ARG APP_UID=10001
ARG APP_GID=10001
 
RUN apt-get update && apt-get install -y --no-install-recommends \
        gosu \
    && rm -rf /var/lib/apt/lists/* \
    && go install github.com/air-verse/air@v1.67.4 \
    && groupadd -g ${APP_GID} appgroup \
    && useradd -m -u ${APP_UID} -g appgroup -s /usr/sbin/nologin appuser \
    && chown -R appuser:appgroup /workspace /go
 
COPY scripts/docker-entrypoint-dev.sh /usr/local/bin/docker-entrypoint-dev.sh
RUN chmod +x /usr/local/bin/docker-entrypoint-dev.sh
 
COPY --chown=appuser:appgroup . .
 
EXPOSE 4242

ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

ENTRYPOINT ["docker-entrypoint-dev.sh"]

CMD ["air", "-c", ".air.toml"]
 
# =============================================================================
# Stage: runtime
# Minimal, distroless-style production image. No shell, no package
# manager, no root user — only the statically linked binary and certs.
# =============================================================================
FROM debian:${DEBIAN_CODENAME}-slim AS runtime

ARG APP_NAME=scorpion
ARG APP_UID=10001
ARG APP_GID=10001

# libgl1, libx11-6 and friends are the runtime counterparts of the -dev
# packages installed in gui-build-deps: the binary is linked against them
# dynamically, so they must be present here even though there's no
# compiler in this image.
RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
        tzdata \
        libgl1 \
        libx11-6 \
        libxcursor1 \
        libxinerama1 \
        libxrandr2 \
        libxi6 \
        libxext6 \
        libxfixes3 \
        libxxf86vm1 \
        libwayland-client0 \
        libxkbcommon0 \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd -g ${APP_GID} appgroup \
    && useradd -m -u ${APP_UID} -g appgroup -s /usr/sbin/nologin appuser

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

ARG APP_NAME=scorpion

COPY --from=builder /bin/${APP_NAME} /${APP_NAME}
