ARG BUILD_FROM

FROM golang:1.21.4-bullseye as golang

COPY go-backup /go-backup
WORKDIR /go-backup
RUN CGO_ENABLED=0 go build -o bin/go-backup

# FROM $BUILD_FROM
FROM homeassistant/amd64-base:latest

COPY --from=golang /go-backup/bin/go-backup /go-backup

# Default entrypoint from the base image does not pass through the TZ env var for some reason
ENTRYPOINT []

CMD ["/go-backup"]
