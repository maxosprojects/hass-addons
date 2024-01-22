#!/usr/bin/with-contenv bashio
# ==============================================================================
# Home Assistant Community Add-on: AWS S3 Backup
# ==============================================================================
#bashio::log.level "debug"

bashio::log.info "Starting AWS S3 Backup..."

bucket_name="$(bashio::config 'bucket_name')"
storage_class="$(bashio::config 'storage_class' 'STANDARD')"
bucket_region="$(bashio::config 'bucket_region' 'eu-central-1')"
delete_local_backups="$(bashio::config 'delete_local_backups' 'true')"
local_backups_to_keep="$(bashio::config 'local_backups_to_keep' '3')"
monitor_path="/backup"
jq_filter=".backups|=sort_by(.date)|.backups|reverse|.[$local_backups_to_keep:]|.[].slug"

export AWS_ACCESS_KEY_ID="$(bashio::config 'aws_access_key')"
export AWS_SECRET_ACCESS_KEY="$(bashio::config 'aws_secret_access_key')"
export AWS_REGION="$bucket_region"

bashio::log.debug "Using AWS CLI version: '$(aws --version)'"
bashio::log.debug "Command: 'aws s3 sync $monitor_path s3://$bucket_name/ --no-progress --region $bucket_region --storage-class $storage_class'"
aws s3 sync $monitor_path s3://"$bucket_name"/ --no-progress --region "$bucket_region" --storage-class "$storage_class"

bashio::log.info "Finished AWS S3 Backup."
