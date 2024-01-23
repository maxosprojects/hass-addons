ARG BUILD_FROM
# FROM $BUILD_FROM

FROM golang:1.21.4-bullseye as golang

COPY go-config /go-config
WORKDIR /go-config
RUN go build -o bin/go-config

FROM grafana/agent:v0.39.1

COPY --from=golang /go-config/bin/go-config /go-config

COPY agent/grafana-agent-config.tmpl /grafana-agent-config.tmpl
COPY agent/run.sh /
RUN chmod a+x /run.sh

# Default entrypoint from the base image does not pass through the TZ env var for some reason
ENTRYPOINT []

CMD ["/run.sh"]

LABEL \
  io.hass.version="2024.1.3" \
  io.hass.type="addon" \
  io.hass.arch="amd64"
