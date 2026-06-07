FROM golang:1.25-alpine AS builder
LABEL maintainer="Wei Lee <mai@mai0313.com>"

RUN apk add --no-cache make git
WORKDIR /app

# Pre-fetch modules so a source-only change doesn't re-download deps.
# `go.*` covers go.mod and (when present) go.sum.
COPY go.* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make

########################################################################################

FROM alpine:3.21 AS prod
LABEL maintainer="Wei Lee <mai@mai0313.com>" \
    org.label-schema.name="go_template" \
    org.label-schema.vendor="Wei Lee" \
    org.label-schema.schema-version="1.0" \
    com.centurylinklabs.watchtower.stop-signal="SIGINT"

RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/build/go_template /usr/local/bin/go_template
RUN chmod +x /usr/local/bin/go_template

# Run as a non-root user so bind-mounted files are not root-owned on the host.
RUN addgroup -g 1000 app \
    && adduser -D -u 1000 -G app app
USER app
CMD ["go_template"]
