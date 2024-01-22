### Forked from https://github.com/thomasfr/hass-addons

Altered to tighten security:
- prevent it from messing around with backups
- only upload backups to S3
- more secure recommended bucket settings

# Home Assistant Add-on: AWS S3 Backup

## Installation

Follow these steps to get the add-on installed on your system:

1. Enable **Advanced Mode** in your Home Assistant user profile.
2. Navigate in your Home Assistant frontend to **Supervisor** -> **Add-on Store**.
3. Search for "AWS S3 Backup" add-on and click on it.
4. Click on the "INSTALL" button.

## How to use

1. Set the `aws_access_key`, `aws_secret_access_key`, and `bucket_name`. 
2. Optionally / if necessary, change `bucket_region`, `storage_class`, and `delete_local_backups` and `local_backups_to_keep` configuration options.
3. Start the add-on to sync the `/backup/` directory to the configured `bucket_name` on AWS S3. You can also automate this of course, see example below:

## Automation

```yaml
alias: Upload backups to S3
description: ""
trigger:
  - platform: time_pattern
    minutes: /10
action:
  - service: hassio.addon_start
    data:
      addon: c33755f6_aws_s3_backup
```

The automation above first creates a full backup at 4am, and then at 4:15am syncs to AWS S3 and if configured deletes local backups according to your configuration.

## Configuration

Example add-on configuration:

```
aws_access_key: AKXXXXXXXXXXXXXXXX
aws_secret_access_key: XXXXXXXXXXXXXXXX
bucket_name: my-bucket
bucket_region: us-west-2
storage_class: STANDARD
```

### Option: `aws_access_key` (required)
AWS IAM access key used to access the S3 bucket.

### Option: `aws_secret_access_key` (required)
AWS IAM secret access key used to access the S3 bucket.

### Option: `bucket_name` (required)
AWS S3 bucket used to store backups.

### Option: `bucket_region` (optional, Default: us-west-2)
AWS region where the S3 bucket was created. See https://aws.amazon.com/about-aws/global-infrastructure/ for all available regions.

### Option: `storage_class` (optional, Default: STANDARD)
AWS S3 storage class to use for the synced objects, when uploading files to S3. One of STANDARD, REDUCED_REDUNDANCY, STANDARD_IA, ONEZONE_IA, INTELLIGENT_TIERING, GLACIER, DEEP_ARCHIVE. For more information see https://aws.amazon.com/s3/storage-classes/.

If you also want to automatically delete backups to keep your AWS S3 Bucket clean, or change the storage class for backups to save some money, you should take a look at S3 Lifecycle Rules (https://docs.aws.amazon.com/AmazonS3/latest/userguide/how-to-set-lifecycle-configuration-intro.html).

## Security

### AWS IAM

Create a new AWS IAM user, which:
- can not login to the AWS Console
- can only access AWS programmatically
- is used by this add-on only
- uses the lowest possible IAM Policy, which is this:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AllowAWSS3Sync",
            "Effect": "Allow",
            "Action": [
                "s3:PutObject",
                "s3:ListBucket"
            ],
            "Resource": [
                "arn:aws:s3:::YOUR-S3-BUCKET-NAME/*",
                "arn:aws:s3:::YOUR-S3-BUCKET-NAME"
            ]
        }
    ]
}
```

### Bucket

1. Enable versioning
2. Enable Object Lock with default retention period of 10 years. Objects can still be deleted from the AWS Console (JS takes care of the lock)
3. Create Lifecycle Rule `Delete expired object delete markers` with `Delete expired object delete markers` set to 7 days

## Development

See [README](https://github.com/maxosprojects/hass-addons)

### Docker

#### Build image

Image tag is formed after the version in [config.yaml](./config.yaml).

```shell
docker build --tag=maxosprojects/hass-aws-s3-backup:1.2.2 .
docker push maxosprojects/hass-aws-s3-backup:1.2.2
```

## Support

Usage of the addon requires knowledge of AWS S3 and AWS IAM.
Under the hood it uses the aws cli version 1, specifically the `aws s3 sync` command.

## Thanks
This addon is highly inspired by https://github.com/gdrapp/hass-addons and https://github.com/rrostt/hassio-backup-s3
