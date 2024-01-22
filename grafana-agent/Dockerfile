ARG BUILD_FROM
# FROM $BUILD_FROM

FROM golang:1.21.4-bullseye as golang

COPY go-get-token /go-get-token
WORKDIR /go-get-token
RUN go build -o bin/get-token

FROM grafana/agent:v0.39.0

COPY --from=golang /go-get-token/bin/get-token /get-token

COPY agent/grafana-agent-config.yaml /etc/agent/agent.yaml
COPY agent/run.sh /
RUN chmod a+x /run.sh

# Default entrypoint from the base image does not pass through the TZ env var for some reason
ENTRYPOINT []

CMD ["/run.sh"]
