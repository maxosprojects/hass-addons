# Home Assistant Addon: Backup to S3

See [README](https://github.com/maxosprojects/hass-addons)

## Use in HA

### Install in HA

1. See [README](https://github.com/maxosprojects/hass-addons)
2. Enable `Watchdog`

### Configure in HA

## Development

See [README](https://github.com/maxosprojects/hass-addons)

### Docker

#### Build image

Image tag is formed after the version in [config.yaml](./config.yaml).

```shell
docker build --tag=maxosprojects/hass-backup-to-s3:1.0.2 .
docker push maxosprojects/hass-backup-to-s3:1.0.2
```

#### Deploy locally

1. See [README](https://github.com/maxosprojects/hass-addons)
2. Click on the addon in the `Local` repo > Build > Start
3. Enable `Watchdog`
