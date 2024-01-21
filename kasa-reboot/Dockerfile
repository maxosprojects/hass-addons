ARG BUILD_FROM

FROM golang:1.21.4-bullseye as golang

COPY go-kasa-reboot /kasa-reboot
WORKDIR /kasa-reboot
RUN CGO_ENABLED=0 go build -o bin/go-kasa-reboot

# FROM $BUILD_FROM
FROM homeassistant/amd64-base:latest

COPY --from=golang /kasa-reboot/bin/go-kasa-reboot /go-kasa-reboot

# Default entrypoint from the base image does pass through the TZ env var for some reason
ENTRYPOINT []

CMD ["/go-kasa-reboot"]
