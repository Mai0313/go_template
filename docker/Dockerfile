FROM golang:1.23-alpine AS builder
LABEL maintainer="Wei Lee <mai@mai0313.com>"

RUN apk add --no-cache make git
WORKDIR /app
COPY . .
RUN make build_linux_amd64

########################################################################################

FROM alpine:3.21 AS prod
LABEL maintainer="Wei Lee <mai@mai0313.com>" \
    org.label-schema.name="go-template" \
    org.label-schema.vendor="Wei Lee" \
    org.label-schema.schema-version="1.0" \
    com.centurylinklabs.watchtower.stop-signal="SIGINT"

RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/build/go-template-linux-amd64 /usr/local/bin/go-template
RUN chmod +x /usr/local/bin/go-template
CMD ["go-template"]
