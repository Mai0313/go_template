FROM golang:1.24-alpine AS builder
LABEL maintainer="Wei Lee <mai@mai0313.com>"

RUN apk add --no-cache make git
WORKDIR /app
COPY . .
RUN make

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
CMD ["go_template"]
