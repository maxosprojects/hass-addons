#!/bin/bash
set -e

/go-config
# sleep infinity

/bin/grafana-agent --config.file=/grafana-agent-config.yaml --metrics.wal-directory=/data/wal
