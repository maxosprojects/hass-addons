#!/bin/bash
set -e

export GCLOUD_HOSTED_LOGS_ID=$(/get-token gcloud_hosted_logs_id)
export TOKEN=$(/get-token grafana_cloud_token)

# echo "TOKEN=$TOKEN"
# sleep infinity

/bin/grafana-agent --config.file=/etc/agent/agent.yaml --metrics.wal-directory=/data/wal -config.expand-env
