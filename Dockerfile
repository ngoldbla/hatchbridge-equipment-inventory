# Node dependencies stage
FROM public.ecr.aws/docker/library/node:22-alpine AS frontend-dependencies
WORKDIR /app

# Install pnpm globally (caching layer)
RUN npm install -g pnpm

# Copy package.json and lockfile to leverage caching
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

# Build Nuxt (frontend) stage
FROM public.ecr.aws/docker/library/node:22-alpine AS frontend-builder
WORKDIR /app

# Install pnpm globally again (it can reuse the cache if not changed)
RUN npm install -g pnpm

# Copy over source files and node_modules from dependencies stage
COPY frontend . 
COPY --from=frontend-dependencies /app/node_modules ./node_modules
RUN pnpm build

# Go dependencies stage
FROM public.ecr.aws/docker/library/golang:alpine AS builder-dependencies
WORKDIR /go/src/app

# Copy go.mod and go.sum for better caching
COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download

# Build API stage
FROM public.ecr.aws/docker/library/golang:alpine AS builder
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG BUILD_TIME
ARG COMMIT
ARG VERSION

# Install necessary build tools
RUN apk update && \
    apk upgrade && \
    apk add --no-cache git build-base gcc g++ && \
    if [ "$TARGETARCH" != "arm" ] && [ "$TARGETARCH" != "riscv64" ]; then apk --no-cache add libwebp libavif libheif libjxl; fi

WORKDIR /go/src/app

# Copy Go modules (from dependencies stage) and source code
COPY --from=builder-dependencies /go/pkg/mod /go/pkg/mod
COPY ./backend .

# Clear old public files and copy new ones from frontend build
RUN rm -rf ./app/api/public
COPY --from=frontend-builder /app/.output/public ./app/api/static/public

# Build Go binary (cache mount removed for Railway compatibility)
RUN if [ "$TARGETARCH" = "arm" ] || [ "$TARGETARCH" = "riscv64" ];  \
    then echo "nodynamic" $TARGETOS $TARGETARCH; CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
        -ldflags "-s -w -X main.commit=$COMMIT -X main.buildTime=$BUILD_TIME -X main.version=$VERSION" \
        -tags nodynamic -o /go/bin/api -v ./app/api/*.go; \
    else \
         echo $TARGETOS $TARGETARCH; CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
        -ldflags "-s -w -X main.commit=$COMMIT -X main.buildTime=$BUILD_TIME -X main.version=$VERSION" \
        -o /go/bin/api -v ./app/api/*.go; \
    fi

# Production stage
FROM public.ecr.aws/docker/library/alpine:latest
ARG TARGETARCH=amd64
ENV HBOX_MODE=production
ENV HBOX_STORAGE_CONN_STRING=file:///?no_tmp_dir=true
ENV HBOX_STORAGE_PREFIX_PATH=data
ENV HBOX_DATABASE_SQLITE_PATH=/data/homebox.db?_pragma=busy_timeout=2000&_pragma=journal_mode=WAL&_fk=1&_time_format=sqlite

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates wget && \
    if [ "$TARGETARCH" != "arm" ] && [ "$TARGETARCH" != "riscv64" ]; then apk --no-cache add libwebp libavif libheif libjxl; fi

# Create application directories and copy over built Go binary
RUN mkdir -p /app /data
COPY --from=builder /go/bin/api /app
RUN chmod +x /app/api

# Runtime entrypoint (maps Railway's $PORT -> HBOX_WEB_PORT)
COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Labels and configuration for the final image
LABEL Name=homebox Version=0.0.1
LABEL org.opencontainers.image.source="https://github.com/sysadminsmedia/homebox"

# Expose necessary ports for Homebox
EXPOSE 7745
WORKDIR /app

# Healthcheck configuration
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD sh -c 'port="${HBOX_WEB_PORT:-}"; if [ -z "$port" ]; then port="${PORT:-7745}"; fi; wget --no-verbose --tries=1 -O - "http://localhost:${port}/api/v1/status" >/dev/null 2>&1 || exit 1'

# Data lives under /data. Do not declare Docker VOLUME here (Railway bans it);
# configure a Railway volume mounted at /data instead.

# Entrypoint and CMD
ENTRYPOINT [ "/entrypoint.sh" ]
CMD [ "/data/config.yml" ]
